package base

type DailyStat struct {
	Name         string
	Now          float64
	ChgToday     float64
	Last         float64
	ChgLast      float64
	PE           float64 `sc:"PE"`
	ChgMonth     float64 `sc:"chg_m"`
	ChgLastMonth float64 `sc:"chg_lm"`
	ChgYear      float64 `sc:"chg_y"`
	ChgMax       float64
	ChgMin       float64
	Chg5         float64
	Chg10        float64
	Chg90        float64
	Avg20        float64
	Avg60        float64
	Avg200       float64
	Code         string
	PB           float64 `sc:"PB"`
}

// KlineDaily 日线数据记录
type KlineDaily struct {
	Date     string
	Open     string
	Close    string
	High     string
	Low      string
	Volume   string
	Amount   string
	ChgRate  string
	PreClose string
}
