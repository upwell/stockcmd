package baostock

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/pkg/errors"

	"hehan.net/my/stockcmd/logger"
)

type ResultSet struct {
	Fields       []string
	MsgType      string
	ReqBodyParts []string
	CurRowNum    int
	CurPageNum   int
	Data         [][]string
	RespMsg      *ResponseMessage
	BS           *BaoStock
}

func (rs *ResultSet) Next() (bool, error) {
	if len(rs.Data) == 0 {
		return false, nil
	}
	if rs.CurRowNum < len(rs.Data) {
		return true, nil
	}

	// request for next page
	nextPage := rs.CurPageNum + 1
	rs.ReqBodyParts[2] = strconv.Itoa(nextPage)
	msgBody := strings.Join(rs.ReqBodyParts, MessageSplit)
	rspMsg, err := rs.BS.request(rs.MsgType, msgBody)
	if err != nil {
		return false, errors.Wrap(err, "failed to get data with error")
	}
	if len(rspMsg.BodyAttrs) < 7 {
		return false, errors.Errorf("wrong number of body attrs of response message: [%s]", rspMsg.BodyAttrs)
	}

	rs.CurPageNum, _ = strconv.Atoi(rspMsg.BodyAttrs[4])
	rs.CurRowNum = 0
	rs.setData(rspMsg.BodyAttrs[6])

	if len(rs.Data) == 0 {
		return false, nil
	}
	return true, nil
}

func (rs *ResultSet) GetRowData() []string {
	if rs.CurRowNum < len(rs.Data) {
		ret := rs.Data[rs.CurRowNum]
		rs.CurRowNum += 1
		return ret
	} else {
		return []string{}
	}
}

func (rs *ResultSet) setData(rawData string) {
	rawData = strings.TrimSpace(rawData)
	if len(rawData) == 0 {
		rs.Data = nil
	} else {
		parts := strings.Split(rawData, " ")
		jsonData := strings.Join(parts, "")
		// jsonData example:
		// {"record":[["2020-01-02","20.3800","20.4800","20.1400","20.3200","20.3200","9934211","202012745.4100","0.000000"]]}
		var f map[string][][]string
		err := json.Unmarshal([]byte(jsonData), &f)
		if err != nil {
			logger.SugarLog.Warnf("parse json err: [%v]", err)
			rs.Data = make([][]string, 0)
		} else {
			record := f["record"]
			if len(record) == 0 || (len(rs.Fields) != 0) && len(record[0]) != len(rs.Fields) {
				logger.SugarLog.Warnf("empty record or invalid record missing fields [%s]", jsonData)
				rs.Data = make([][]string, 0)
			}
			rs.Data = f["record"]
		}
	}
}
