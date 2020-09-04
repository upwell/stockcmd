package tencent

import (
	"strconv"
	"strings"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
	"hehan.net/my/stockcmd/hq"
)

const hqURL = "http://qt.gtimg.cn/q="

type TencentHQApi struct {
}

func (api TencentHQApi) GetHQ(code string) (hq.HQ, error) {
	code = hq.ConvertCode(code)
	ret := hq.HQ{}
	resp, err := grequests.Get(hqURL+code, nil)
	if err != nil {
		return ret, errors.Wrapf(err, "failed to get live price for [%s]", code)
	}
	if resp.StatusCode != 200 {
		return ret, errors.Errorf("failed to get live price for [%s], status code [%d]",
			code, resp.StatusCode)
	}

	rawResult, err := hq.ConvertGBK2UTF8(resp.String())
	if err != nil {
		return ret, errors.Wrapf(err, "convert GBK2UTF8 failed")
	}
	rawResult = strings.Split(rawResult, "=")[1]
	rawResult = strings.ReplaceAll(rawResult, "\"", "")
	parts := strings.Split(rawResult, "~")
	if len(parts) <= 49 {
		return ret, errors.Errorf("wrong number of parts")
	}
	now, err := strconv.ParseFloat(parts[3], 64)
	last, err := strconv.ParseFloat(parts[4], 64)
	chgToday, err := strconv.ParseFloat(parts[32], 64)

	ret.Now = now
	ret.ChgToday = chgToday
	ret.Last = last
	ret.IsSuspend = false
	return ret, nil
}
