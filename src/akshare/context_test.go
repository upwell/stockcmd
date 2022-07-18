package akshare

import (
	"testing"
	"time"

	"hehan.net/my/stockcmd/util"

	"hehan.net/my/stockcmd/logger"
)

func TestAKShare_QueryAllStock(t *testing.T) {
	entries := AK.QueryAllStock()
	if entries == nil {
		println("nil result")
	}
	for idx, entry := range entries {
		if idx > 5 {
			break
		}
		println(entry.Name, entry.Code)
	}
}

func TestAKShare_WriteBasics(t *testing.T) {
	logger.InitLogger()
	entries := AK.QueryAllStock()
	if entries == nil {
		println("nil result")
	}
	AK.WriteBasics(entries)
}

func TestAKShare_GetDailyKData(t *testing.T) {
	logger.InitLogger()
	startTime, _ := time.Parse("2006-01-02", "2022-05-01")
	endTime, _ := time.Parse("2006-01-02", "2022-05-20")
	entries := AK.GetDailyKData("sh.603777", startTime, endTime)

	for _, entry := range entries {
		println(util.DateToStr(time.Time(entry.Date)))
	}
}
