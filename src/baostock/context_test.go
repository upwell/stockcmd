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
