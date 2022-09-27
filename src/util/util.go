package util

import (
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/fatih/color"
)

func DateToStr(time time.Time) string {
	return time.Format("2006-01-02")
}

func DateToStr2(time time.Time) string {
	return time.Format("20060102")
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

func RoundChgRate(rate float64) float64 {
	return math.Round(rate*100*100) / 100
}

func Round2(val float64) float64 {
	return math.Round(val*100) / 100
}

func Round3(val float64) float64 {
	return math.Round(val*1000) / 1000
}

func Float64String(f float64) string {
	return fmt.Sprintf("%.2f", f)
}

func JsonString(j map[string]interface{}) string {
	jBytes, err := json.Marshal(j)
	if err != nil {
		return ""
	}
	return string(jBytes)
}

func ChgString(chg float64, fallRate float64, riseRate float64) string {
	chgStr := Float64String(chg)
	post := ""
	switch {
	case chg > riseRate:
		post = "✨"
	case chg > 0:
		post = "↑"
	case chg == 0:
		post = "⁃"
	case chg < fallRate:
		post = "⚡"
	case chg < 0:
		post = "↓"
	}
	chgStr = fmt.Sprintf("%s %s", chgStr, post)
	if chg >= riseRate {
		chgStr = color.RedString(chgStr)
	} else if chg <= fallRate {
		chgStr = color.GreenString(chgStr)
	}
	return chgStr
}

func IsTradeTime(t time.Time) bool {
	weekday := t.Weekday()
	if weekday == time.Sunday || weekday == time.Saturday {
		return false
	}

	if t.Hour() < 9 || t.Hour() >= 15 {
		return false
	}

	return true
}
