package stat

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"hehan.net/my/stockcmd/eastmoney"

	"hehan.net/my/stockcmd/akshare"

	"hehan.net/my/stockcmd/tencent"

	"hehan.net/my/stockcmd/util"

	"gonum.org/v1/gonum/floats"

	"hehan.net/my/stockcmd/logger"

	"gonum.org/v1/gonum/stat"

	"github.com/rocketlaunchr/dataframe-go"

	"github.com/fatih/color"
	"github.com/iancoleman/strcase"
	"github.com/jinzhu/now"
	"github.com/pkg/errors"
	"hehan.net/my/stockcmd/baostock"
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
	row = append(row, Float64String(ds.Now))
	row = append(row, ChgString(ds.ChgToday, -3, 3))
	row = append(row, Float64String(ds.Last))
	row = append(row, ChgString(ds.ChgLast, -3, 3))
	row = append(row, Float64String(ds.PE))
	row = append(row, Float64String(ds.ChgMonth))
	row = append(row, Float64String(ds.ChgLastMonth))
	row = append(row, Float64String(ds.ChgYear))
	row = append(row, ChgString(ds.ChgMax, -6, 5))
	row = append(row, Float64String(ds.ChgMin))
	row = append(row, Float64String(ds.Chg5))
	row = append(row, Float64String(ds.Chg10))
	row = append(row, Float64String(ds.Chg90))
	row = append(row, Float64String(ds.Avg20))
	row = append(row, Float64String(ds.Avg60))
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
	return RoundChgRate((lastClose - firstPreClose) / firstPreClose)
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
	return RoundChgRate((lastClose - firstPreClose) / firstPreClose)
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
	return Round2(stat.Mean(values, nil))
}

func GetDataFrame(code string) (*dataframe.DataFrame, error) {
	t := store.GetLastTime(code)
	endDay := now.BeginningOfDay()
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

	if endDay.After(startDay) {
		logger.SugarLog.Infof("get history data for [%s] between [%s] and [%s]", code,
			util.DateToStr(startDay), util.DateToStr(endDay))
		v, err := baostock.BSPool.Get()
		if err != nil {
			return nil, errors.Wrap(err, "failed to get baostock instance from pool")
		}
		bs := v.(*baostock.BaoStock)
		t1 := time.Now()
		rs, err := bs.GetDailyKData(code, startDay, endDay)
		if err != nil {
			baostock.BSPool.Close(v)
			return nil, errors.Wrap(err, "get daily state failed")
		}
		records := make([]*store.Record, 0, 1024)
		for {
			hasNext, err := rs.Next()
			if err != nil {
				baostock.BSPool.Close(v)
				return nil, errors.Wrap(err, "get daily state, error in loop")
			}
			if !hasNext {
				break
			}

			seps := rs.GetRowData()
			skipThisRow := false
			for _, sep := range seps {
				if len(sep) == 0 {
					skipThisRow = true
				}
			}
			if skipThisRow {
				continue
			}
			date, _ := now.Parse(seps[0])
			records = append(records, &store.Record{
				Code: code,
				T:    date,
				Val:  strings.Join(seps, ","),
			})
		}
		baostock.BSPool.Put(v)
		store.WriteRecords(records)
		logger.SugarLog.Debugf("[%s] get remote data takes [%v]", code, time.Since(t1))
	}

	df, err := store.GetRecords(code, endDay.AddDate(-1, 0, 0), endDay)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get records from db for [%s]", code)
	}
	return df, nil
}

func GetDataFrameAKShare(code string) (*dataframe.DataFrame, error) {
	t := store.GetLastTime(code)
	endDay := now.BeginningOfDay()
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

	if endDay.After(startDay) {
		logger.SugarLog.Infof("get history data for [%s] between [%s] and [%s]", code,
			util.DateToStr(startDay), util.DateToStr(endDay))

		t1 := time.Now()
		dailyDataArray := akshare.AK.GetDailyKData(code, startDay, endDay)
		records := make([]*store.Record, 0, 1024)

		for _, kdata := range dailyDataArray {
			//date,open,high,low,close,preclose,volume,amount,pctChg,peTTM,pbMRQ
			val := fmt.Sprintf("%s,%f,%f,%f,%f,%f,%f,%f,%f,0,0",
				util.DateToStr(time.Time(kdata.Date)), kdata.Open, kdata.High, kdata.Low, kdata.Close, 0.0,
				kdata.Volume, kdata.Amount, kdata.ChgRate)

			record := &store.Record{
				Code: code,
				T:    time.Time(kdata.Date),
				Val:  val,
			}
			records = append(records, record)
		}

		store.WriteRecords(records)
		logger.SugarLog.Debugf("[%s] get remote data takes [%v]", code, time.Since(t1))
	}

	df, err := store.GetRecords(code, endDay.AddDate(-1, 0, 0), endDay)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get records from db for [%s]", code)
	}
	return df, nil
}

func GetDataFrameEastMoney(code string) (*dataframe.DataFrame, error) {
	t := store.GetLastTime(code)
	endDay := now.BeginningOfDay()
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

	if endDay.After(startDay) {
		logger.SugarLog.Infof("get history data for [%s] between [%s] and [%s]", code,
			util.DateToStr(startDay), util.DateToStr(endDay))

		t1 := time.Now()
		dailyDataArray, err := eastmoney.EM.GetDailyKData(code, startDay, endDay)
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
		store.WriteRecords(records)
		logger.SugarLog.Debugf("[%s] write records takes [%v]", code, time.Since(t2))
	}

	df, err := store.GetRecords(code, endDay.AddDate(-1, 0, 0), endDay)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get records from db for [%s]", code)
	}
	return df, nil
}

func GetDailyState(code string, period int) (*DailyStat, error) {
	df, err := GetDataFrame(code)
	if err != nil {
		return nil, err
	}

	name := store.GetName(code, false)
	api := tencent.HQApi{}
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
		ChgMax:       RoundChgRate((hq.Now - max) / max),
		ChgMin:       RoundChgRate((hq.Now - min) / min),
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
