package baostock_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"hehan.net/my/stockcmd/util"

	"hehan.net/my/stockcmd/store"

	"github.com/jinzhu/now"
	"github.com/rocketlaunchr/dataframe-go/imports"

	"hehan.net/my/stockcmd/baostock"
	"hehan.net/my/stockcmd/logger"
)

func setup() {
	logger.InitLogger()
}

func TestMain(m *testing.M) {
	fmt.Println("inside TestMain")
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestBaoStock_LoginLogout(t *testing.T) {
	fmt.Println("start testing ..")
	err := baostock.BS.Login()
	if err != nil {
		t.Errorf("login failed with [%v]", err)
	}

	err = baostock.BS.Logout()
	if err != nil {
		t.Errorf("logout failed with [%v]", err)
	}
}

func TestBaoStock_QueryAll(t *testing.T) {
	baostock.BS.Login()
	defer baostock.BS.Logout()

	rs, err := baostock.BS.QueryAllStock(util.GetLastWorkDay())
	if err != nil {
		t.Errorf("query failed with [%v]", err)
	}

	store.WriteBasics(rs.Data)
}

func TestBaoStock_GetDailyKData(t *testing.T) {
	baostock.BS.Login()
	defer baostock.BS.Logout()

	fromDate, _ := now.Parse("2020-08-24")
	toDate, _ := now.Parse("2020-08-25")

	rs, _ := baostock.BS.GetDailyKData("sz.002475", fromDate, toDate)
	fmt.Println(rs.RespMsg.BodyAttrs)
}

func TestBaoStock_QueryDividendData(t *testing.T) {
	baostock.BS.Login()
	defer baostock.BS.Logout()

	rs, _ := baostock.BS.QueryDividendData("sz.002475", "2021")
	rows := make([]string, 0, 20)
	for {
		hasNext, err := rs.Next()
		if !hasNext || err != nil {
			break
		}
		rows = append(rows, strings.Join(rs.GetRowData(), ","))
	}

	for _, row := range rows {
		println(row)
	}
}

func TestBaoStock_GetLastDividendDay(t *testing.T) {
	baostock.BS.Login()
	defer baostock.BS.Logout()

	start := time.Now()
	day, err := baostock.BS.GetLastDividendDay("sz.002475")
	if err != nil {
		println(err)
		return
	}

	println(day.String())
	println(time.Since(start).String())
}

func TestBaoStock_AllGetLastDividendDay(t *testing.T) {
	baostock.BS.Login()
	defer baostock.BS.Logout()

	start := time.Now()

	codeSet := store.GetAllStockCodes()
	for code := range codeSet.Iter() {
		d, _ := baostock.BS.GetLastDividendDay(code.(string))
		println(d.String())
	}

	println(time.Since(start).String())
}

func TestBaoStock_AllGetLastDividendDayWithPool(t *testing.T) {
	start := time.Now()
	codeSet := store.GetAllStockCodes()
	var wg sync.WaitGroup
	for code := range codeSet.Iter() {
		wg.Add(1)
		go func(c string) {
			defer wg.Done()
			v, err := baostock.BSPool.Get()
			if err != nil {
				fmt.Printf("failed to obtain baostock instance [%v]\n", err)
				return
			}
			bs := v.(*baostock.BaoStock)
			d, err := bs.GetLastDividendDay(c)
			if err != nil {
				fmt.Printf("something is wrong [%v]\n", err)
				baostock.BSPool.Close(v)
				return
			}
			baostock.BSPool.Put(v)
			println(d.String())
		}(code.(string))
	}

	wg.Wait()
	println(time.Since(start).String())
}

func TestBaoStock_QueryHistoryKDataPage(t *testing.T) {
	baostock.BS.Login()
	fromDate, _ := now.Parse("2020-06-13")
	toDate, _ := now.Parse("2020-07-01")
	rs, err := baostock.BS.QueryHistoryKDataPage(1, 200, "sz.002475",
		baostock.DailyDataFields, fromDate, toDate, "d",
		"2")
	if err != nil {
		t.Errorf("get daily data failed [%v]", err)
	}

	//dfRecords := make([][]string, 0, 1024)
	//dfRecords = append(dfRecords, rs.Fields)
	//for rs.Next() {
	//	dfRecords = append(dfRecords, rs.GetRowData())
	//}

	//df := dataframe.LoadRecords(dfRecords, dataframe.DetectTypes(false), dataframe.DefaultType(series.Float),
	//	dataframe.WithTypes(map[string]series.Type{
	//		"date": series.String,
	//	}))
	//fmt.Println(df)
	csvRows := make([]string, 0, 1024)
	csvRows = append(csvRows, strings.Join(rs.Fields, ","))
	for {
		hasNext, err := rs.Next()
		if !hasNext || err != nil {
			break
		}
		csvRows = append(csvRows, strings.Join(rs.GetRowData(), ","))
	}
	csvStr := strings.Join(csvRows, "\n")

	ctx := context.Background()
	dataTypes := make(map[string]interface{})
	for _, field := range rs.Fields {
		dataTypes[field] = float64(0)
	}
	dataTypes["date"] = imports.Converter{
		ConcreteType: time.Time{},
		ConverterFunc: func(in interface{}) (i interface{}, err error) {
			return time.Parse("2006-01-02", in.(string))
		},
	}
	opts := imports.CSVLoadOptions{
		TrimLeadingSpace: true,
		LargeDataSet:     false,
		DictateDataType:  dataTypes,
		InferDataTypes:   false,
	}
	df, err := imports.LoadFromCSV(ctx, strings.NewReader(csvStr), opts)
	if err != nil {
		t.Errorf("load data to dataframe failed [%v]", err)
	}
	fmt.Println(df)
}

func TestBaoStockConnNumber(t *testing.T) {
	var ops uint64
	var wg sync.WaitGroup

	for i := 0; i < 10000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			instance := baostock.NewBaoStockInstance()
			err := instance.Login()
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(-1)
			}
			atomic.AddUint64(&ops, 1)
			fmt.Println(ops)
			time.Sleep(time.Minute * 5)
		}()
	}

	wg.Wait()
}

func TestTimeCompare(t *testing.T) {
	var t1 time.Time
	t1 = now.With(t1).BeginningOfDay()
	t2 := now.BeginningOfDay()
	fmt.Println(t1.After(t2))
	fmt.Println(t1.Equal(t2))
}
