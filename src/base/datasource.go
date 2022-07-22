package base

import (
	"time"
)

type DataSource interface {
	GetDailyKData(code string, startDay time.Time, endDay time.Time) ([]KlineDaily, error)
}
