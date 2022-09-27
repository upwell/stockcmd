package eastmoney

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jinzhu/now"

	"hehan.net/my/stockcmd/base"
	"hehan.net/my/stockcmd/util"

	"github.com/pkg/errors"

	"github.com/levigross/grequests"
)

var EM *EastMoney

// EastMoney east money DataSource
type EastMoney struct {
}

func init() {
	EM = &EastMoney{}
}

func (em EastMoney) GetCodeIDMap() (map[string]string, error) {
	url := "http://80.push2.eastmoney.com/api/qt/clist/get"
	params := map[string]string{
		"pn":     "1",
		"pz":     "50000",
		"po":     "1",
		"np":     "1",
		"ut":     "bd1d9ddb04089700cf9c27f6f7426281",
		"fltt":   "2",
		"invt":   "2",
		"fid":    "f3",
		"fs":     "m:1 t:2,m:1 t:23",
		"fields": "f12",
		"_":      "1623833739532",
	}
	resp, err := grequests.Get(url, &grequests.RequestOptions{Params: params})
	if err != nil {
		return nil, errors.Wrapf(err, "failed on request to [%s]", url)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Wrapf(err, "failed on request to [%s], status code [%d]", url, resp.StatusCode)
	}

	params = map[string]string{
		"pn":     "1",
		"pz":     "50000",
		"po":     "1",
		"np":     "1",
		"ut":     "bd1d9ddb04089700cf9c27f6f7426281",
		"fltt":   "2",
		"invt":   "2",
		"fid":    "f3",
		"fs":     "m:0 t:6,m:0 t:80",
		"fields": "f12",
		"_":      "1623833739532",
	}
	resp, err = grequests.Get(url, &grequests.RequestOptions{Params: params})
	println(resp.String())

	return nil, nil
}

func (em EastMoney) GetDailyKData(code string, startDay time.Time, endDay time.Time) ([]base.KlineDaily, error) {
	secId := strings.ReplaceAll(code, "sh", "1")
	secId = strings.ReplaceAll(secId, "sz", "0")
	url := "http://push2his.eastmoney.com/api/qt/stock/kline/get"
	params := map[string]string{
		"fields1": "f1,f2,f3,f4,f5,f6",
		"fields2": "f51,f52,f53,f54,f55,f56,f57,f58,f59,f60,f61,f116",
		"ut":      "7eea3edcaed734bea9cbfc24409ed989",
		"klt":     "101",
		"fqt":     "0",
		"secid":   secId,
		"beg":     util.DateToStr2(startDay),
		"end":     util.DateToStr2(endDay),
		"_":       "1623766962675",
	}
	resp, err := grequests.Get(url, &grequests.RequestOptions{Params: params})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to request [%s] with [%v]", url, err)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Errorf("failed to request [%s] with status code [%d]",
			url, resp.StatusCode)
	}

	var respJson map[string]interface{}
	err = json.Unmarshal(resp.Bytes(), &respJson)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse response [%s]", resp.String())
	}

	result := make([]base.KlineDaily, 0, 64)
	dataJson := respJson["data"].(map[string]interface{})
	preCloseFloat := dataJson["preKPrice"].(float64)
	preClose := fmt.Sprintf("%.2f", preCloseFloat)
	if klinesJson, ok := dataJson["klines"]; ok {
		klinesAInterfs := klinesJson.([]interface{})
		klinesArray := make([]string, len(klinesAInterfs))
		for i, _ := range klinesAInterfs {
			klinesArray[i] = klinesAInterfs[i].(string)
		}
		for _, klineStr := range klinesArray {
			parts := strings.Split(klineStr, ",")
			KlineDaily := base.KlineDaily{
				Date:     parts[0],
				Open:     parts[1],
				Close:    parts[2],
				High:     parts[3],
				Low:      parts[4],
				Volume:   parts[5],
				Amount:   parts[6],
				ChgRate:  parts[8],
				PreClose: preClose,
			}
			preClose = KlineDaily.Close
			result = append(result, KlineDaily)
		}
	} else {
		return result, nil
	}

	return result, nil
}

func (em EastMoney) GetLastDividendDay(code string) (time.Time, error) {
	codeParts := strings.Split(code, ".")
	url := "https://datacenter-web.eastmoney.com/api/data/v1/get"
	params := map[string]string{
		"callback":     "jQuery112306373219885568009_1664246756161",
		"sortColumns":  "REPORT_DATE",
		"sortTypes":    "-1",
		"pageSize":     "50",
		"pageNumber":   "1",
		"reportName":   "RPT_SHAREBONUS_DET",
		"columns":      "ALL",
		"quoteColumns": "",
		"js":           `{"data":(x),"pages":(tp)}`,
		"source":       "WEB",
		"client":       "WEB",
		"filter":       `(SECURITY_CODE="` + codeParts[1] + `")`,
	}

	resp, err := grequests.Get(url, &grequests.RequestOptions{Params: params})
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "failed to request [%s] with [%v]", url, err)
	}
	if resp.StatusCode != 200 {
		return time.Time{}, errors.Errorf("failed to request [%s] with status code [%d]",
			url, resp.StatusCode)
	}

	// response string format: jQuery1122_11122({"version": ""})
	// get the json string
	rawRespStr := resp.String()
	startIdx := strings.Index(rawRespStr, "(")
	endIdx := strings.LastIndex(rawRespStr, ")")
	jsonStr := rawRespStr[startIdx+1 : endIdx]

	var respJson map[string]interface{}
	err = json.Unmarshal([]byte(jsonStr), &respJson)
	if err != nil {
		return time.Time{}, errors.Wrapf(err, "failed to parse response [%s]", jsonStr)
	}

	resultItr := respJson["result"]
	if resultItr == nil {
		return time.Time{}, nil
	}
	resultJson := resultItr.(map[string]interface{})

	dataJsonArray := resultJson["data"].([]interface{})
	if len(dataJsonArray) == 0 {
		return time.Time{}, nil
	}

	latestRecord := dataJsonArray[0].(map[string]interface{})
	exDivDateItr := latestRecord["EX_DIVIDEND_DATE"]
	if exDivDateItr == nil {
		return time.Time{}, nil
	}

	ret, err := now.Parse(exDivDateItr.(string))
	if err != nil {
		return time.Time{}, errors.Errorf("wrong date string format: [%s]", latestRecord["EX_DIVIDEND_DATE"].(string))
	}
	return ret, nil
}
