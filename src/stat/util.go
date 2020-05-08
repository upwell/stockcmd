package stat

import "math"

func RoundChgRate(rate float64) float64 {
	return math.Round(rate*100*100) / 100
}

func Round2(val float64) float64 {
	return math.Round(val*100) / 100
}
