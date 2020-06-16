package baostock_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

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

func TestBaoStock_GetDailyKData(t *testing.T) {
	baostock.BS.Login()

	fromDate, _ := now.Parse("2020-04-01")
	toDate, _ := now.Parse("2020-05-01")

	rs, _ := baostock.BS.GetDailyKData("sz.002475", fromDate, toDate)
	fmt.Println(rs.RespMsg.BodyAttrs)

	baostock.BS.Logout()
}

func TestBaoStock_QueryHistoryKDataPage(t *testing.T) {
	baostock.BS.Login()
	fromDate, _ := now.Parse("2019-05-01")
	toDate, _ := now.Parse("2020-05-01")
	rs, err := baostock.BS.QueryHistoryKDataPage(1, 200, "sz.002475",
		baostock.DailyDataFields, fromDate, toDate, "d",
		"3")
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
	for rs.Next() {
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
