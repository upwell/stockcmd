package stat

import (
	"fmt"
	"sync"
	"time"

	"hehan.net/my/stockcmd/eastmoney"

	"hehan.net/my/stockcmd/global"
	"hehan.net/my/stockcmd/store"

	"hehan.net/my/stockcmd/util"
)

func FetchAllHQ() {
	fmt.Println("Fetch all latest hq data ...")
	//t := time.Now()
	//if !util.IsTradeTime(t) {
	//	fmt.Println("Not trading time, skip fetching")
	//	return
	//}

	api := eastmoney.HQApi{}
	hqs, err := api.GetAllHQ()
	if err != nil {
		fmt.Printf("failed to fetch all hq with [%v]\n", err)
		return
	}
	store.BulkWriteHQ(hqs)
}

func FetchAllHQ2() {
	fmt.Println("Fetch all latest hq data ...")
	t := time.Now()
	if !util.IsTradeTime(t) {
		fmt.Println("Not trading time, skip fetching")
		return
	}

	codes := store.GetCodes()
	start := time.Now()
	var wg sync.WaitGroup
	hqs := make([]*store.StockHQ, 0, 512)
	for _, code := range codes {
		wg.Add(1)
		//api := global.GetHQSource()
		api := global.GetRandomHQSource()
		go func(code string) {
			defer wg.Done()
			t1 := time.Now()
			v, err := api.GetHQ(code)
			if err != nil {
				fmt.Printf("failed to get price for [%s] with error [%v]\n", code, err)
				return
			}
			fmt.Printf("Take [%s] to fetch [%s]\n", time.Since(t1), code)
			if v.IsSuspend {
				fmt.Printf("%s is suspend\n", code)
			}
			if v.Now == 0.00 && v.Last == 0.00 {
				fmt.Printf("%s now and last is zero\n", code)
			}

			hq := &store.StockHQ{
				Code:  code,
				Price: fmt.Sprintf("%f", v.Now),
			}
			hqs = append(hqs, hq)
		}(code)
	}

	wg.Wait()
	store.BulkWriteHQ(hqs)
	fmt.Printf("fetch hq data done, take [%s], start fetch history records ... \n", time.Since(start))
}
