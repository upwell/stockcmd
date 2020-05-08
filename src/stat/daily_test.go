package stat

import (
	"encoding/json"
	"fmt"
	"testing"
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
