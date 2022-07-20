package redisstore

import (
	"fmt"
	"testing"

	"hehan.net/my/stockcmd/store"

	"hehan.net/my/stockcmd/util"

	"github.com/jinzhu/now"
	"hehan.net/my/stockcmd/eastmoney"
	"hehan.net/my/stockcmd/logger"
)

const testCode = "sz.300821"

func TestWriteRecords(t *testing.T) {
	logger.InitLogger()

	start, _ := now.Parse("2022-07-01")
	end, _ := now.Parse("2022-07-19")
	kdataList, _ := eastmoney.EM.GetDailyKData(testCode, start, end)

	records := make([]*store.Record, 0, 32)
	for _, kdata := range kdataList {
		t, _ := now.Parse(kdata.Date)
		val := fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s,%s,0,0",
			kdata.Date, kdata.Open, kdata.High, kdata.Low, kdata.Close, kdata.PreClose,
			kdata.Volume, kdata.Amount, kdata.ChgRate)
		record := &store.Record{
			Code: testCode,
			T:    t,
			Val:  val,
		}
		records = append(records, record)
	}
	WriteRecords(records)
}

func TestGetRecords(t *testing.T) {
	logger.InitLogger()

	start, _ := now.Parse("2022-07-01")
	end, _ := now.Parse("2022-07-15")
	df, _ := GetRecords(testCode, start, end)
	print(df.Table())
}

func TestGetLastTime(t *testing.T) {
	tt := GetLastTime(testCode)
	println(util.DateToStr2(tt))
}
