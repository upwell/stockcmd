package store_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"gonum.org/v1/gonum/floats"

	"github.com/rocketlaunchr/dataframe-go"

	"hehan.net/my/stockcmd/store"

	"github.com/jinzhu/now"
	"gonum.org/v1/gonum/stat"
	"hehan.net/my/stockcmd/baostock"
)

const testCode = "sz.300821"

func TestWriteRecord(t *testing.T) {
	baostock.BS.Login()
	fromDate, _ := now.Parse("2020-01-01")
	toDate, _ := now.Parse("2020-03-01")
	rs, err := baostock.BS.QueryHistoryKDataPage(1, 10, testCode,
		"date,open,high,low,close,preclose,volume,amount,pctChg", fromDate, toDate, "d",
		"3")
	if err != nil {
		t.Errorf("get daily data failed [%v]", err)
	} else {
		for {
			hasNext, err := rs.Next()
			if !hasNext || err != nil {
				break
			}
			seps := rs.GetRowData()
			t, _ := time.Parse("2006-01-02", seps[0])
			store.WriteRecord(testCode, t, strings.Join(seps, ","))
		}
	}
}

func TestGetLastTime(t *testing.T) {
	d := store.GetLastTime(testCode)
	if d.IsZero() {
		t.Errorf("failed to get last time")
	} else {
		fmt.Println(d)
	}
}

func TestGetRecords(t *testing.T) {
	fromDate, _ := now.Parse("2020-01-01")
	toDate, _ := now.Parse("2020-05-01")
	df, err := store.GetRecords(testCode, fromDate, toDate)
	if err != nil {
		t.Errorf("failed to get record [%v]", err)
	} else {
		if len(df.Series) > 0 {
			fmt.Println(df.Table())
		} else {
			fmt.Println("empty result")
		}
	}
}

func TestCalAvg(t *testing.T) {
	fromDate, _ := now.Parse("2020-01-01")
	toDate, _ := now.Parse("2020-03-01")
	df, err := store.GetRecords(testCode, fromDate, toDate)
	if err != nil {
		t.Errorf("failed to get record [%v]", err)
	} else {
		if len(df.Series) > 0 {
			idx, err := df.NameToColumn("close")
			if err != nil {
				t.Errorf("failed to get column index [%v]", err)
				return
			}
			opens := df.Series[idx].(*dataframe.SeriesFloat64)

			avg := stat.Mean(opens.Values, nil)
			max := floats.Max(opens.Values)
			min := floats.Min(opens.Values)
			fmt.Printf("avg: %f, max: %f, min: %f\n", avg, max, min)
		}
	}
}
