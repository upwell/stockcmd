package hq

type HQ struct {
	Now       float64
	Last      float64
	ChgToday  float64
	MarketCap float64
	IsSuspend bool
}

type HQApi interface {
	GetHQ(code string) (HQ, error)
}
