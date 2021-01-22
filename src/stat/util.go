package stat

import (
	"fmt"
	"math"

	"github.com/fatih/color"
)

func RoundChgRate(rate float64) float64 {
	return math.Round(rate*100*100) / 100
}

func Round2(val float64) float64 {
	return math.Round(val*100) / 100
}

func Float64String(f float64) string {
	return fmt.Sprintf("%.2f", f)
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
