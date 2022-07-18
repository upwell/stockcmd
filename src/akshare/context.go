package akshare

import (
	"encoding/json"
	"time"

	"go.etcd.io/bbolt"
	"hehan.net/my/stockcmd/store"

	"hehan.net/my/stockcmd/logger"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

var AK *AKShare

type AKShare struct {
}

func init() {
	AK = &AKShare{}
}

func dateParamStr(t time.Time) string {
	return t.Format("20060102")
}

func (ak AKShare) request(msgType string, params map[string]string) ([]byte, error) {
	var ro *grequests.RequestOptions
	if params == nil {
		ro = nil
	} else {
		ro = &grequests.RequestOptions{
			Params: params,
		}
	}
	resp, err := grequests.Get(SERVER+msgType, ro)
	if err != nil {
		return nil, errors.Wrapf(err, "failed on request to [%s]", msgType)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Wrapf(err, "failed on request to [%s], status code [%d]",
			msgType, resp.StatusCode)
	}

	return resp.Bytes(), nil
}

func (ak AKShare) QueryAllStock() []StockInfoResponse {
	respBytes, err := ak.request("stock_info_a_code_name", nil)
	if err != nil {
		return nil
	}

	var entries []StockInfoResponse
	err = json.Unmarshal(respBytes, &entries)
	if err != nil {
		logger.SugarLog.Warnf("get all stock info response parse error [%v]", err)
		return nil
	}

	return entries
}

func (ak AKShare) GetDailyKData(code string, startDay time.Time, endDay time.Time) []KDataDaily {
	params := make(map[string]string)
	params["symbol"] = code[3:]
	//params["symbol"] = strings.ReplaceAll(code, ".", "")
	params["period"] = "daily"
	params["start_date"] = dateParamStr(startDay)
	params["end_date"] = dateParamStr(endDay)
	//respBytes, err := ak.request("stock_zh_a_hist_163", params)
	respBytes, err := ak.request("stock_zh_a_hist", params)
	if err != nil {
		logger.SugarLog.Warnf("get daily kdata for [%s] failed with error [%v]", code, err)
		return nil
	}

	var entries []KDataDaily
	err = json.Unmarshal(respBytes, &entries)
	if err != nil {
		logger.SugarLog.Warnf("get daily kdata parse error [%v] with [%s]",
			err, string(respBytes[:]))
		return nil
	}
	return entries
}

func (ak AKShare) WriteBasics(stockInfos []StockInfoResponse) {
	store.DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(store.BasicBucketName))
		var basic *store.StockBasic
		for _, stockInfo := range stockInfos {
			codePrefix := DetermineExchangeByCode(stockInfo.Code)
			if len(codePrefix) == 0 {
				logger.SugarLog.Errorf("unknown code format [%s]", stockInfo.Code)
				continue
			}

			// TODO ignore bj
			if codePrefix == "bj" {
				continue
			}

			code := codePrefix + "." + stockInfo.Code
			basic = &store.StockBasic{
				Code: code,
				Name: stockInfo.Name,
			}

			bytes, _ := json.Marshal(basic)
			b.Put([]byte(basic.Code), bytes)

		}
		return nil
	})
}
