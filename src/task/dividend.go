package task

import (
	"sync"
	"time"

	"hehan.net/my/stockcmd/util"

	"github.com/jinzhu/now"

	"hehan.net/my/stockcmd/baostock"

	"hehan.net/my/stockcmd/logger"
	"hehan.net/my/stockcmd/store"
)

const lastGetStockDividendTimeKey string = "lastGetStockDividendTime"

// CheckAllStockDividendDay check the dividend info of all stock
func CheckAllStockDividendDay(force bool) {
	nowTime := time.Now()

	if !force && (nowTime.Weekday() == time.Saturday || nowTime.Weekday() == time.Sunday || nowTime.Hour() < 9) {
		// not trading time, ignore
		return
	}
	lastTimeStamp, err := store.RunningConfig.GetInt64OrDefault(lastGetStockDividendTimeKey, 0)
	if err != nil {
		logger.SugarLog.Errorf("failed to get last getting dividend error [%v]", err)
		return
	}

	lastTime := time.Unix(lastTimeStamp, 0)
	if !force && util.DateEqual(nowTime, lastTime) {
		logger.SugarLog.Debug("already check today")
		return
	}

	codeSet := store.GetAllStockCodes()
	var wg sync.WaitGroup
	for code := range codeSet.Iter() {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			v, err := baostock.BSPool.Get()
			var d time.Time
			if err != nil {
				logger.SugarLog.Warnf("failed to obtain baostock instance [%v]\n", err)
				// delete records if fail
				store.DeleteCodeRecords(c)
				return
			}
			bs := v.(*baostock.BaoStock)
			d, err = bs.GetLastDividendDay(c)
			logger.SugarLog.Debugf("code [%s] last dividend day %s", c, d.String())
			if err != nil {
				logger.SugarLog.Errorf("something is wrong [%v]\n", err)
				// delete records if fail
				store.DeleteCodeRecords(c)
				baostock.BSPool.Close(v)
				return
			}
			baostock.BSPool.Put(v)

			lastDay := now.New(lastTime).BeginningOfDay()
			if !d.IsZero() && (lastDay.Equal(d) || lastDay.Before(d)) {
				logger.SugarLog.Infof("[%s] has dividend event since last query, delete old records", c)
			}
		}(code.(string))
	}
	wg.Wait()

	store.RunningConfig.Set(lastGetStockDividendTimeKey, nowTime.Unix())
}
