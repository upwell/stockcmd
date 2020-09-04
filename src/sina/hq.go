package sina

import (
	"math"
	"strconv"
	"strings"
	"sync"

	"hehan.net/my/stockcmd/hq"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

type SinaHQApi struct {
}

func GetLivePrices(codes []string) map[string]hq.HQ {
	var wg sync.WaitGroup
	ret := make(map[string]hq.HQ)

	api := SinaHQApi{}
	for _, code := range codes {
		wg.Add(1)
		go func(code string) {
			defer wg.Done()
			hq, err := api.GetHQ(code)
			if err == nil {
				ret[code] = hq
			}
		}(code)
	}
	wg.Wait()
	return ret
}

// GetLivePrice code format should be: sh.000001
func (api SinaHQApi) GetHQ(code string) (hq.HQ, error) {
	code = hq.ConvertCode(code)
	ret := hq.HQ{}
	resp, err := grequests.Get(HqURL+code, nil)
	if err != nil {
		return ret, errors.Wrapf(err, "failed to get live price for [%s]", code)
	}
	if resp.StatusCode != 200 {
		return ret, errors.Errorf("failed to get live price for [%s], status code [%d]",
			code, resp.StatusCode)
	}

	rawResult, err := hq.ConvertGB2UTF8(resp.String())
	if err != nil {
		return ret, errors.Wrapf(err, "failed to get live price")
	}
	rawResult = strings.Split(rawResult, "=")[1]
	rawResult = strings.ReplaceAll(rawResult, "\"", "")
	parts := strings.Split(rawResult, ",")
	if len(parts) < 4 {
		return ret, errors.Errorf("wrong result from live price [%s]", rawResult)
	}

	now, err := strconv.ParseFloat(parts[3], 64)
	close, err := strconv.ParseFloat(parts[2], 64)
	isSuspendFloat, err := strconv.ParseFloat(parts[1], 64)
	ret.Now = now
	ret.ChgToday = math.Round(((now-close)/close)*100*100) / 100
	ret.Last = close
	ret.IsSuspend = isSuspendFloat == 0.0
	return ret, nil
}
