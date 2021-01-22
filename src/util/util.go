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
