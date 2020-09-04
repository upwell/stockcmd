package util

import "time"

func DateToStr(time time.Time) string {
	return time.Format("2006-01-02")
}
