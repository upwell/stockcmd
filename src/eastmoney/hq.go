package eastmoney

import (
	"encoding/json"
	"strings"

	"hehan.net/my/stockcmd/base"

	"hehan.net/my/stockcmd/logger"

	"hehan.net/my/stockcmd/util"

	"hehan.net/my/stockcmd/store"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

type HQApi struct {
}

func (api HQApi) GetHQ(code string) (base.HQ, error) {
	secId := strings.ReplaceAll(code, "sh", "1")
	secId = strings.ReplaceAll(secId, "sz", "0")
	ret := base.HQ{}

	url := "https://push2.eastmoney.com/api/qt/ulist.np/get"
	params := map[string]string{
		"OSVersion":     "14.3",
		"appVersion":    "6.3.8",
		"fields":        "f12,f14,f3,f2,f15,f16,f17,f4,f8,f10,f9,f5,f6,f18,f20,f21,f13,f124,f297",
		"fltt":          "2",
		"plat":          "Iphone",
		"product":       "EFund",
		"secids":        secId,
		"serverVersion": "6.3.6",
		"version":       "6.3.8",
	}
	headers := map[string]string{
		"User-Agent":      "Mozilla/5.0 (Windows NT 6.3; WOW64; Trident/7.0; Touch; rv:11.0) like Gecko",
		"Accept":          "*/*",
		"Accept-Language": "zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2",
	}

	resp, err := grequests.Get(url, &grequests.RequestOptions{Params: params, Headers: headers})
	if err != nil {
		return ret, errors.Wrapf(err, "failed to request [%s] with [%v]", url, err)
	}

	if resp.StatusCode != 200 {
		return ret, errors.Errorf("failed to request [%s] with status code [%d]", url, resp.StatusCode)
	}

	var respJson map[string]interface{}
	err = json.Unmarshal(resp.Bytes(), &respJson)
	if err != nil {
		return ret, errors.Wrapf(err, "failed to parse response [%s]", resp.String())
	}

	dataJson := respJson["data"].(map[string]interface{})
	diffJson, exists := dataJson["diff"]
	if !exists {
		return ret, errors.Errorf("empty result [%s]", resp.String())
	}
	hqsInterfs := diffJson.([]interface{})
	if len(hqsInterfs) == 0 {
		return ret, errors.Errorf("empty result [%s]", resp.String())
	}
	hqJson := hqsInterfs[0].(map[string]interface{})

	priceFloat, ok := hqJson["f2"].(float64)
	if !ok {
		logger.SugarLog.Warnf("invalid price [%s]", util.JsonString(hqJson))
		ret.IsSuspend = true
	}
	ret = base.HQ{
		Now:       priceFloat,
		ChgToday:  hqJson["f3"].(float64),
		Last:      hqJson["f18"].(float64),
		MarketCap: hqJson["f20"].(float64),
	}

	return ret, nil
}

func (api HQApi) GetAllHQ() ([]*store.StockHQ, error) {
	url := "http://82.push2.eastmoney.com/api/qt/clist/get"
	params := map[string]string{
		"pn":     "1",
		"pz":     "50000",
		"po":     "1",
		"np":     "1",
		"ut":     "bd1d9ddb04089700cf9c27f6f7426281",
		"fltt":   "2",
		"invt":   "2",
		"fid":    "f3",
		"fs":     "m:0 t:6,m:0 t:80,m:1 t:2,m:1 t:23,m:0 t:81 s:2048",
		"fields": "f1,f2,f3,f4,f5,f6,f7,f8,f9,f10,f12,f13,f14,f15,f16,f17,f18,f20,f21,f23,f24,f25,f22,f11,f62,f128,f136,f115,f152",
		"_":      "1623833739532",
	}

	resp, err := grequests.Get(url, &grequests.RequestOptions{Params: params})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to request [%s] with [%v]", url, err)
	}

	if resp.StatusCode != 200 {
		return nil, errors.Errorf("failed to request [%s] with status code [%d]", url, resp.StatusCode)
	}

	var respJson map[string]interface{}
	err = json.Unmarshal(resp.Bytes(), &respJson)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse response [%s]", resp.String())
	}

	dataJson := respJson["data"].(map[string]interface{})
	diffJson, exists := dataJson["diff"]
	if !exists {
		return nil, errors.Errorf("empty results return from [%s]", url)
	}
	result := make([]*store.StockHQ, 0, 64)

	hqsInterfs := diffJson.([]interface{})
	hqsArray := make([]map[string]interface{}, len(hqsInterfs))
	for i, _ := range hqsArray {
		hqsArray[i] = hqsInterfs[i].(map[string]interface{})
	}
	for _, hqJson := range hqsArray {
		//close, _ := strconv.ParseFloat(parts[17], 64)
		//chgRate, _ := strconv.ParseFloat(parts[2], 64)

		//hqStr, _ := json.Marshal(hqJson)
		//println(string(hqStr))

		code := hqJson["f12"].(string)
		marketF := hqJson["f13"].(float64)
		var marketCode string
		switch marketF {
		case 0:
			marketCode = "sz"
			break
		case 1:
			marketCode = "sh"
			break
		case 2:
			marketCode = "bj"
			break
		default:
			logger.SugarLog.Warnf("unknown market [%f] for [%s]", marketF, code)
			continue
		}
		priceFloat, ok := hqJson["f2"].(float64)
		if !ok {
			logger.SugarLog.Warnf("failed to get price for [%s]", code)
			continue
		}

		hq := &store.StockHQ{
			Code:  marketCode + "." + code,
			Price: util.Float64String(priceFloat),
		}
		result = append(result, hq)
	}
	return result, nil
}
