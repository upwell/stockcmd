package stat

import (
	"sort"
	"sync"
	"time"

	"hehan.net/my/stockcmd/global"

	"hehan.net/my/stockcmd/redisstore"

	"hehan.net/my/stockcmd/util"

	"gonum.org/v1/gonum/floats"

	"github.com/rocketlaunchr/dataframe-go"

	"hehan.net/my/stockcmd/logger"

	"github.com/jinzhu/now"
	"hehan.net/my/stockcmd/store"
)

type RPS struct {
	Code   string
	Name   string
	Change float64
	Value  float64
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetRPS(basics []*store.StockBasic, days int, endDay time.Time) []*RPS {
	rpss := make([]*RPS, 0, 512)
	changes := make([]float64, 0, 512)

	if endDay.IsZero() {
		endDay = now.BeginningOfDay()
	}
	var wg sync.WaitGroup

	if now.New(endDay).BeginningOfDay() == now.BeginningOfDay() {
		// 获取当前行情数据
		FetchAllHQ()
	}

	// 取2倍日期数据
	startDay := endDay.AddDate(0, 0, -(days * 2))
	for _, basic := range basics {
		wg.Add(1)
		go func(basic *store.StockBasic) {
			defer wg.Done()

			code := basic.Code
			rps := &RPS{
				Code: basic.Code,
				Name: basic.Name,
			}
			GetDataFrame(global.GetDataSource(), code)
			df, err := redisstore.GetRecords(code, startDay, endDay)
			if err != nil {
				logger.SugarLog.Errorf("get records for [%s] error [%v]", code, err)
				return
			}
			if df.NRows() == 0 {
				logger.SugarLog.Errorf("get records return zero rows for [%s]", code)
				return
			}

			closeIdx, err := df.NameToColumn("close")
			if err != nil {
				logger.SugarLog.Errorf("failed to get column index [%v]", err)
				return
			}
			closes := df.Series[closeIdx].(*dataframe.SeriesFloat64)

			var price float64
			if now.New(endDay).BeginningOfDay() == now.BeginningOfDay() {
				price = store.GetHQ(code)
				if price == 0.00 {
					logger.SugarLog.Infof("failed to get price for [%s]", code)
					return
				}
			} else {
				price = closes.Values[0]
			}

			//if code == "sz.001236" {
			//	println(df.Table())
			//	println(price)
			//	println(closes.Values[5])
			//}

			rps.Change = price / closes.Values[min(len(closes.Values)-1, days)]

			rpss = append(rpss, rps)
			changes = append(changes, rps.Change)
		}(basic)
	}
	wg.Wait()

	n := len(changes)
	for _, rps := range rpss {
		cmpF1 := func(v float64) bool { return v < rps.Change }
		cmpF2 := func(v float64) bool { return v <= rps.Change }
		left := floats.Count(cmpF1, changes)
		right := floats.Count(cmpF2, changes)

		tmpV := 0
		if right > left {
			tmpV = 1
		}
		rps.Value = util.Round3(float64(right+left+tmpV) * 50.0 / float64(n))
	}

	sort.Slice(rpss, func(i, j int) bool {
		return rpss[i].Value > rpss[j].Value
	})
	return rpss
}
