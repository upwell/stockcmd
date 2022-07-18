package akshare

import (
	"encoding/json"
	"strings"
	"time"
)

type JsonDate time.Time

func (j *JsonDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*j = JsonDate(t)
	return nil
}

func (j JsonDate) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(j))
}

func (j JsonDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

type StockInfoResponse struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type KDataDaily struct {
	Date    JsonDate `json:"日期"`
	Open    float64  `json:"开盘"`
	High    float64  `json:"最高"`
	Low     float64  `json:"最低"`
	Close   float64  `json:"收盘"`
	Volume  float64  `json:"成交量"`
	Amount  float64  `json:"成交额"`
	ChgRate float64  `json:"涨跌幅"`
}

type KDataDaily163 struct {
	Date     time.Time `json:"日期"`
	Open     float64   `json:"开盘价"`
	High     float64   `json:"最高价"`
	Low      float64   `json:"最低价"`
	Close    float64   `json:"收盘价"`
	PreClose float64   `json:"前收盘"`
	Volume   float64   `json:"成交量"`
	Amount   float64   `json:"成交额"`
	ChgRate  float64   `json:"涨跌幅"`
}
