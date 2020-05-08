package sina

import (
	"math"
	"strconv"
	"strings"
	"sync"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

type HQ struct {
	Now       float64
	Last      float64
	ChgToday  float64
	IsSuspend bool
}

func GetLivePrices(codes []string) map[string]*HQ {
	var wg sync.WaitGroup
	ret := make(map[string]*HQ)

	for _, code := range codes {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			hq, err := GetLivePrice(code)
			if err == nil {
				ret[code] = hq
			}
		}(code)
	}
	wg.Wait()
	return ret
}

// GetLivePrice code format should be: sh000001
func GetLivePrice(code string) (*HQ, error) {
	code = ConvertCode(code)
	resp, err := grequests.Get(HqURL+code, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get live price for [%s]", code)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("failed to get live price for [%s], status code [%d]",
			code, resp.StatusCode)
	}

	rawResult, err := ConvertGB2UTF8(resp.String())
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get live price")
	}
	rawResult = strings.Split(rawResult, "=")[1]
	rawResult = strings.ReplaceAll(rawResult, "\"", "")
	parts := strings.Split(rawResult, ",")
	if len(parts) < 4 {
		return nil, errors.Errorf("wrong result from live price [%s]", rawResult)
	}

	now, err := strconv.ParseFloat(parts[3], 64)
	close, err := strconv.ParseFloat(parts[2], 64)
	isSuspendFloat, err := strconv.ParseFloat(parts[1], 64)
	hq := &HQ{
		Now:       now,
		ChgToday:  math.Round(((now-close)/close)*100*100) / 100,
		Last:      close,
		IsSuspend: isSuspendFloat == 0.0,
	}
	return hq, nil
}
