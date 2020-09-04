package stat

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"

	"hehan.net/my/stockcmd/logger"

	"hehan.net/my/stockcmd/store"
)

func TestFields(t *testing.T) {
	ds := DailyStat{}
	fmt.Println(Fields(ds))
}

func TestGetDailyState(t *testing.T) {
	code := "sz.002475"
	ds, err := GetDailyState(code)
	if err != nil {
		t.Errorf("get daily state error [%v]", err)
		return
	}
	dsBytes, _ := json.MarshalIndent(ds, "", "  ")
	fmt.Println(string(dsBytes))
	fmt.Println(ds.Row())
}

func TestGetDataFrame(t *testing.T) {
	logger.InitLogger()
	code := "sz.300821"
	df, err := GetDataFrame(code)
	if err != nil {
		t.Error(err)
		return
	}
	if df.NRows() == 0 {
		t.Error("zero row")
	}
}

func TestAllGetDataFrame(t *testing.T) {
	logger.InitLogger()
	codes := store.GetCodes()
	var wg sync.WaitGroup
	for _, code := range codes {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			fmt.Println(code)
			_, err := GetDataFrame(code)
			if err != nil {
				t.Errorf("get data frame error [%v]", err)
				return
			}
		}(code)
	}
	wg.Wait()
}
