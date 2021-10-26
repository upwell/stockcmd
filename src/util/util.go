package util

import (
	"fmt"
	"time"
)

func DateToStr(time time.Time) string {
	return time.Format("2006-01-02")
}

func GetLastWorkDay() time.Time {
	today := time.Now()

	lastWorkDay := today.AddDate(0, 0, -1)
	switch lastWorkDay.Weekday() {
	case time.Sunday:
		lastWorkDay = lastWorkDay.AddDate(0, 0, -2)
	case time.Saturday:
		lastWorkDay = lastWorkDay.AddDate(0, 0, -1)
	}
	return lastWorkDay
}

func printTimeMeasurement(funcName string, start time.Time) func() {
	return func() {
		fmt.Printf("Time taken by %s function is %v \n", funcName, time.Since(start))
	}
}

func MeasureTime(funcName string) func() {
	start := time.Now()
	return printTimeMeasurement(funcName, start)
}

// DateEqual check if two time objects are in same day
func DateEqual(date1, date2 time.Time) bool {
	y1, m1, d1 := date1.Date()
	y2, m2, d2 := date2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
