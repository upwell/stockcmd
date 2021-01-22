package store

import (
	"fmt"
	"sync"
	"testing"

	"hehan.net/my/stockcmd/tencent"

	"gotest.tools/assert"
	"hehan.net/my/stockcmd/sina"
)

func TestWriteHQ(t *testing.T) {
	code := "sh.000961"
	api := sina.HQApi{}
	hq, _ := api.GetHQ(code)

	WriteHQ(&StockHQ{
		Code:  code,
		Price: fmt.Sprintf("%f", hq.Now),
	})

	price := GetHQ(code)

	assert.Equal(t, price, hq.Now)
}

func TestBulkWriteHQ(t *testing.T) {
	codes := GetCodes()

	var wg sync.WaitGroup
	hqs := make([]*StockHQ, 0, 512)
	for _, code := range codes {
		wg.Add(1)
		api := tencent.HQApi{}
		go func(code string) {
			defer wg.Done()
			v, err := api.GetHQ(code)
			if err != nil {
				fmt.Printf("failed to get price for [%s]\n", code)
				return
			}
			if v.IsSuspend {
				fmt.Printf("%s is suspend\n", code)
			}
			if v.Now == 0.00 && v.Last == 0.00 {
				fmt.Printf("%s now and last is zero\n", code)
			}
			hq := &StockHQ{
				Code:  code,
				Price: fmt.Sprintf("%f", v.Now),
			}
			hqs = append(hqs, hq)
		}(code)
	}

	wg.Wait()
	BulkWriteHQ(hqs)
}

func TestGetHQ(t *testing.T) {
	code := "sh.000939"
	price := GetHQ(code)
	fmt.Println(price)
}
