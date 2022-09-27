package stat

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"hehan.net/my/stockcmd/global"

	"hehan.net/my/stockcmd/base"

	"hehan.net/my/stockcmd/redisstore"

	"hehan.net/my/stockcmd/util"

	"gonum.org/v1/gonum/floats"

	"hehan.net/my/stockcmd/logger"

	"gonum.org/v1/gonum/stat"

	"github.com/rocketlaunchr/dataframe-go"

	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
	"hehan.net/my/stockcmd/store"
)

type DailyStat struct {
	Name         string
	Now          float64
	ChgToday     float64
	Last         float64
	ChgLast      float64
	PE           float64 `sc:"PE"`
	ChgMonth     float64 `sc:"chg_m"`
	ChgLastMonth float64 `sc:"chg_lm"`
	ChgYear      float64 `sc:"chg_y"`
	ChgMax       float64
	ChgMin       float64
	Chg5         float64
	Chg10        float64
	Chg90        float64
	Avg20        float64
	Avg60        float64
	Avg200       float64
	Code         string
	PB           float64 `sc:"PB"`
}

func Fields(s interface{}) []string {
	var fields []string
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		display := v.Type().Field(i).Tag.Get("sc")
		if len(display) == 0 {
			display = strcase.ToSnake(v.Type().Field(i).Name)
		}
		fields = append(fields, display)
	}
	return fields
}

func (ds *DailyStat) Row() []string {
	row := make([]string, 0, 32)

	nameStr := ds.Name
	if ds.Avg20 > ds.Now {
		nameStr = color.BlueString(ds.Name)
	}
	if ds.Avg60 > ds.Now {
		nameStr = color.GreenString(ds.Name)
	}
	if ds.Avg200 > ds.Now {
		nameStr = color.HiGreenString(ds.Name)
	}
	row = append(row, nameStr)
	row = append(row, util.Float64String(ds.Now))
	row = append(row, util.ChgString(ds.ChgToday, -3, 3))
	row = append(row, util.Float64String(ds.Last))
	row = append(row, util.ChgString(ds.ChgLast, -3, 3))
	row = append(row, util.Float64String(ds.PE))
	row = append(row, util.Float64String(ds.ChgMonth))
	row = append(row, util.Float64String(ds.ChgLastMonth))
	row = append(row, util.Float64String(ds.ChgYear))
	row = append(row, util.ChgString(ds.ChgMax, -6, 5))
	row = append(row, util.Float64String(ds.ChgMin))
	row = append(row, util.Float64String(ds.Chg5))
	row = append(row, util.Float64String(ds.Chg10))
	row = append(row, util.Float64String(ds.Chg90))
	row = append(row, util.Float64String(ds.Avg20))
	row = append(row, util.Float64String(ds.Avg60))
	//row = append(row, Float64String(ds.Avg200))
	row = append(row, ds.Code)
	//row = append(row, Float64String(ds.PB))

	return row
}

func thisMonthFilterFn(vals map[interface{}]interface{}, row, nRows int) (dataframe.FilterAction, error) {
	now := time.Now()
	month := now.Month()
	year := now.Year()
	date := vals["date"].(time.Time)
	if date.Month() == month && date.Year() == year {
		return dataframe.KEEP, nil
	} else {
		return dataframe.DROP, nil
	}
}

func lastMonthFilterFn(vals map[interface{}]interface{}, row, nRow int) (dataframe.FilterAction, error) {
	now := time.Now()
	lastNow := now.AddDate(0, -1, 0)
	year := lastNow.Year()
	month := lastNow.Month()

	date := vals["date"].(time.Time)
	if date.Month() == month && date.Year() == year {
		return dataframe.KEEP, nil
	} else {
		return dataframe.DROP, nil
	}
}

func thisYearFilterFn(vals map[interface{}]interface{}, row, nRow int) (dataframe.FilterAction, error) {
	year := time.Now().Year()

	date := vals["date"].(time.Time)
	if date.Year() == year {
		return dataframe.KEEP, nil
	} else {
		return dataframe.DROP, nil
	}
}

func chgWithDf(df *dataframe.DataFrame, fn dataframe.FilterDataFrameFn) float64 {
	ctx := context.Background()
	filterRes, _ := dataframe.Filter(ctx, df, fn)
	filterDf := filterRes.(*dataframe.DataFrame)
	n := filterDf.NRows()
	if n == 0 {
		return 0.00
	}
	lastRow := filterDf.Row(0, false, dataframe.SeriesName)
	firstRow := filterDf.Row(n-1, false, dataframe.SeriesName)
	lastClose := lastRow["close"].(float64)
	firstPreClose := firstRow["preclose"].(float64)
	return util.RoundChgRate((lastClose - firstPreClose) / firstPreClose)
}

func GetMaxMin(df *dataframe.DataFrame, days int) (max float64, min float64) {
	n := df.NRows()
	max = 0.00
	min = 0.00
	if n == 0 {
		logger.SugarLog.Debug("rows of df is 0")
		return
	}

	r := days - 1
	if days > n-1 {
		r = n - 1
	}
	if r == 0 {
		return
	}

	idx, err := df.NameToColumn("close")
	if err != nil {
		logger.SugarLog.Error("get column close failed")
		return
	}

	closes := df.Series[idx].(*dataframe.SeriesFloat64).Values[:r]
	max = floats.Max(closes)
	min = floats.Min(closes)
	return
}

func chgDays(df *dataframe.DataFrame, days int) float64 {
	n := df.NRows()
	if n == 0 {
		return 0.00
	}
	lastRow := df.Row(0, false, dataframe.SeriesName)

	r := days - 1
	if days > n-1 {
		r = n - 1
	}
	firstRow := df.Row(r, false, dataframe.SeriesName)
	lastClose := lastRow["close"].(float64)
	firstPreClose := firstRow["preclose"].(float64)
	return util.RoundChgRate((lastClose - firstPreClose) / firstPreClose)
}

func avgDays(df *dataframe.DataFrame, days int) float64 {
	idx, _ := df.NameToColumn("close")
	series := df.Series[idx].(*dataframe.SeriesFloat64)
	values := series.Values
	var n int
	if len(values) < days {
		n = len(values)
	} else {
		n = days
	}
	values = values[0:n]
	return util.Round2(stat.Mean(values, nil))
}

func GetDataFrame(dataSource base.DataSource, code string) (*dataframe.DataFrame, error) {
	t := redisstore.GetLastTime(code)

	currentTime := time.Now()
	var endDay time.Time
	endDay = now.BeginningOfDay()
	if currentTime.Hour() < 15 {
		endDay = endDay.AddDate(0, 0, -1)
	}

	var startDay time.Time
	if t.IsZero() {
		logger.SugarLog.Infof("getting history data for %s, it would take some time ...", code)
		startDay = endDay.AddDate(-1, 0, 0)
	} else {
		startDay = t.AddDate(0, 0, 1)

		// skip weekend
		switch startDay.Weekday() {
		case time.Sunday:
			startDay = startDay.AddDate(0, 0, 1)
		case time.Saturday:
			startDay = startDay.AddDate(0, 0, 2)
		}
	}
	startDay = now.With(startDay).BeginningOfDay()

	if !endDay.Before(startDay) {
		logger.SugarLog.Infof("get history data for [%s] between [%s] and [%s]", code,
			util.DateToStr(startDay), util.DateToStr(endDay))

		t1 := time.Now()
		dailyDataArray, err := dataSource.GetDailyKData(code, startDay, endDay)
		if err != nil {
			logger.SugarLog.Errorf("failed to get daily kdata for [%s]", code)
			return nil, errors.Wrapf(err, "failed to get daily kdata for [%s]", code)
		}
		records := make([]*store.Record, 0, 1024)
		for _, kdata := range dailyDataArray {
			//date,open,high,low,close,preclose,volume,amount,pctChg,peTTM,pbMRQ
			val := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,0,0",
				kdata.Date, kdata.Open, kdata.High, kdata.Low, kdata.Close, kdata.PreClose,
				kdata.Volume, kdata.Amount, kdata.ChgRate)

			t, _ := now.Parse(kdata.Date)
			record := &store.Record{
				Code: code,
				T:    t,
				Val:  val,
			}
			records = append(records, record)
		}

		t2 := time.Now()
		logger.SugarLog.Debugf("[%s] get remote data takes [%v]", code, time.Since(t1))
		redisstore.WriteRecords(records)
		logger.SugarLog.Debugf("[%s] write records takes [%v]", code, time.Since(t2))
	}

	df, err := redisstore.GetRecords(code, endDay.AddDate(-1, 0, 0), endDay)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get records from db for [%s]", code)
	}
	return df, nil
}

func GetDailyState(dataSource base.DataSource, code string, period int) (*DailyStat, error) {
	df, err := GetDataFrame(dataSource, code)
	if err != nil {
		return nil, err
	}

	name := store.GetName(code, false)
	api := global.GetHQSource()
	hq, err := api.GetHQ(code)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get hq from sina")
	}

	hasRecords := true
	var lastRecord map[interface{}]interface{}
	if df.NRows() == 0 {
		logger.SugarLog.Info("empty history records")
		hasRecords = false
	} else {
		lastRecord = df.Row(0, false, dataframe.SeriesName)
	}
	chgLast := 0.00
	pe := 0.00
	pb := 0.00
	if hasRecords {
		chgLast = lastRecord["pctChg"].(float64)
		pe = lastRecord["peTTM"].(float64)
		pb = lastRecord["pbMRQ"].(float64)
	}
	max, min := GetMaxMin(df, period)
	ds := &DailyStat{
		Name:         name,
		ChgToday:     hq.ChgToday,
		Now:          hq.Now,
		Last:         hq.Last,
		ChgLast:      chgLast,
		ChgMonth:     chgWithDf(df, thisMonthFilterFn),
		ChgLastMonth: chgWithDf(df, lastMonthFilterFn),
		ChgYear:      chgWithDf(df, thisYearFilterFn),
		ChgMax:       util.RoundChgRate((hq.Now - max) / max),
		ChgMin:       util.RoundChgRate((hq.Now - min) / min),
		Avg20:        avgDays(df, 20),
		Avg60:        avgDays(df, 60),
		Avg200:       avgDays(df, 200),
		Chg5:         chgDays(df, 5),
		Chg10:        chgDays(df, 10),
		Chg90:        chgDays(df, 90),
		Code:         code,
		PE:           pe,
		PB:           pb,
	}
	return ds, nil
}
