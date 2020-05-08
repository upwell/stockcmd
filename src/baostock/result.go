package baostock

import (
	"encoding/json"
	"strconv"
	"strings"

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

func (rs *ResultSet) Next() bool {
	if len(rs.Data) == 0 {
		return false
	}
	if rs.CurRowNum < len(rs.Data) {
		return true
	}

	// request for next page
	nextPage := rs.CurPageNum + 1
	rs.ReqBodyParts[2] = strconv.Itoa(nextPage)
	msgBody := strings.Join(rs.ReqBodyParts, MessageSplit)
	rspMsg, err := rs.BS.request(rs.MsgType, msgBody)
	if err != nil {
		logger.SugarLog.Warnf("failed to get data with error [%v]", err)
		return false
	}

	rs.CurPageNum, _ = strconv.Atoi(rspMsg.BodyAttrs[4])
	rs.CurRowNum = 0
	rs.setData(rspMsg.BodyAttrs[6])

	if len(rs.Data) == 0 {
		return false
	}
	return true
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
		// {"record":[["2020-01-02","20.3800","20.4800","20.1400","20.3200","20.3200","9934211","202012745.4100","0.000000"]
		var f map[string][][]string
		err := json.Unmarshal([]byte(jsonData), &f)
		if err != nil {
			logger.SugarLog.Warnf("parse json err: [%v]", err)
			rs.Data = nil
		} else {
			rs.Data = f["record"]
		}
	}
}
