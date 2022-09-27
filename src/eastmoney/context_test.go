package eastmoney

import (
	"testing"
	"time"

	"hehan.net/my/stockcmd/store"

	"github.com/jinzhu/now"

	"hehan.net/my/stockcmd/logger"
)

func TestEastMoney_GetCodeIDMap(t *testing.T) {
	logger.InitLogger()
	EM.GetCodeIDMap()
}

func TestEastMoney_GetDailyKData(t *testing.T) {
	logger.InitLogger()
	start, _ := now.Parse("2022-07-10")
	end, _ := now.Parse("2022-07-20")
	dataList, err := EM.GetDailyKData("sz.300982", start, end)
	if err != nil {
		println(err)
	}
	println(time.Now().Hour())
	print(dataList[len(dataList)-1].Date)
}

func TestEastMoney_GetLastDividendDay(t *testing.T) {
	logger.InitLogger()

	lastT, err := EM.GetLastDividendDay("sh.688008")
	if err != nil {
		println(err)
		t.FailNow()
	}
	println(lastT.String())
}

func TestEastMoney_AllGetLastDividendDay(t *testing.T) {
	start := time.Now()

	codeSet := store.GetAllStockCodes()
	for code := range codeSet.Iter() {
		d, _ := EM.GetLastDividendDay(code.(string))
		println(d.String())
	}

	println(time.Since(start).String())
}
