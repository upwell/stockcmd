package baostock

import (
	"fmt"
	"hash/crc32"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"hehan.net/my/stockcmd/config"
	"hehan.net/my/stockcmd/logger"
	"hehan.net/my/stockcmd/util"

	"github.com/pkg/errors"
)

var BS *BaoStock

type BaoStock struct {
	Conn     net.Conn
	IsLogin  bool
	UserID   string
	mux      sync.Mutex
	loginMux sync.Mutex
}

type ResponseMessage struct {
	MsgType    string
	BodyLength int
	ErrCode    string
	ErrMsg     string
	BodyAttrs  []string
}

func init() {
	BS = NewBaoStockInstance()
}

func NewBaoStockInstance() *BaoStock {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", Server, ServerPort))
	if err != nil {
		logger.SugarLog.Fatalf("failed to connect", Server)
		return nil
	}
	return &BaoStock{Conn: conn}
}

func parseResp(respStr string) *ResponseMessage {
	respHeader := respStr[0:MessageHeaderLength]
	respBody := respStr[MessageHeaderLength:]
	headerAttrs := strings.Split(respHeader, MessageSplit)
	bodyAttrs := strings.Split(respBody, MessageSplit)

	bodyLength, _ := strconv.Atoi(headerAttrs[2])
	return &ResponseMessage{
		MsgType:    headerAttrs[1],
		BodyLength: bodyLength,
		ErrCode:    bodyAttrs[0],
		ErrMsg:     bodyAttrs[1],
		BodyAttrs:  bodyAttrs,
	}
}

func (respMsg ResponseMessage) IsSucceed() bool {
	return respMsg.ErrCode == ErrSuccess
}

func requestElapse(start time.Time, what string) {
	logger.SugarLog.Debugf("%s took %v", what, time.Since(start))
}

func composeRequestString(msgType string, msgBody string) []byte {
	msgHeader := toMessageHeader(msgType, len(msgBody))
	msg := msgHeader + msgBody
	crc32Str := strconv.FormatUint(uint64(crc32.Checksum([]byte(msg), crc32Table)), 10)
	return []byte(msg + MessageSplit + crc32Str + "\n")
}

func (bs *BaoStock) request(msgType string, msgBody string) (*ResponseMessage, error) {
	if config.Verbose {
		defer requestElapse(time.Now(), msgType)
	}

	bs.mux.Lock()
	defer bs.mux.Unlock()

	if _, err := bs.Conn.Write(composeRequestString(msgType, msgBody)); err != nil {
		return nil, errors.Wrap(err, "write request to baostock failed")
	}

	respStr := ""
	for {
		const bufSize = 1024
		// in case there is condition that it reads exact 1024 bytes and no data any more
		if err := bs.Conn.SetReadDeadline(time.Now().Add(time.Second * 2)); err != nil {
			return nil, errors.Wrap(err, "set read dead line failed")
		}
		buf := make([]byte, bufSize)
		n, err := bs.Conn.Read(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				break
			}
			return nil, errors.Wrap(err, "read from baostock failed")
		}
		respStr += string(buf[:n])
		if len(respStr) > 13 && respStr[len(respStr)-13:] == "<![CDATA[]]>\n" {
			break
		}
	}
	if len(respStr) == 0 {
		return nil, errors.New("empty response from baostock")
	}

	respMsg := parseResp(respStr)
	return respMsg, nil
}

func (bs *BaoStock) Login() error {
	// only login once
	bs.loginMux.Lock()
	defer bs.loginMux.Unlock()

	if bs.IsLogin {
		return nil
	}

	userName := AnonymousUserID
	password := "123456"
	bs.UserID = AnonymousUserID

	parts := []string{"login", userName, password, "0"}
	msgBody := strings.Join(parts, MessageSplit)
	respMsg, err := bs.request(MessageTypeLoginRequest, msgBody)
	if err != nil {
		return errors.Wrap(err, "login failed")
	}
	if respMsg.IsSucceed() {
		bs.IsLogin = true
		return nil
	} else {
		return errors.Errorf("login failed with error [%s]:[%s]", respMsg.ErrCode, respMsg.ErrMsg)
	}
}

func (bs *BaoStock) Logout() error {
	if !bs.IsLogin {
		return errors.New("not login yet")
	}
	nowStr := time.Now().Format("%Y%m%d%H%M%S")
	msgBody := strings.Join([]string{
		"logout", bs.UserID, nowStr,
	}, MessageSplit)

	respMsg, err := bs.request(MessageTypeLogoutRequest, msgBody)
	if err != nil {
		return errors.Wrap(err, "logout failed")
	}

	if respMsg.IsSucceed() {
		bs.IsLogin = false
		return nil
	} else {
		return errors.Errorf("logout failed with error [%s]:[%s]", respMsg.ErrCode, respMsg.ErrMsg)
	}
}

func (bs *BaoStock) QueryAllStock(day time.Time) (*ResultSet, error) {
	if day.IsZero() {
		day = time.Now()
	}

	parts := []string{
		"query_all_stock", bs.UserID, "1", "10000", util.DateToStr(day),
	}
	msgBody := strings.Join(parts, MessageSplit)
	respMsg, err := bs.request(MessageTypeQueryAllStockRequest, msgBody)
	if err != nil {
		return nil, errors.Wrap(err, "query all stock failed")
	}
	if !respMsg.IsSucceed() {
		return nil, errors.Errorf("error code [%s], error message [%s]", respMsg.ErrCode, respMsg.ErrMsg)
	}

	rs := &ResultSet{
		MsgType:      MessageTypeQueryAllStockRequest,
		ReqBodyParts: parts,
		Fields:       []string{},
		BS:           bs,
	}

	if len(respMsg.BodyAttrs) < 7 {
		return nil, errors.Errorf("invalid body attrs [%s]", respMsg.BodyAttrs)
	}
	rs.CurPageNum, _ = strconv.Atoi(respMsg.BodyAttrs[4])
	rs.RespMsg = respMsg
	rs.setData(respMsg.BodyAttrs[6])
	return rs, nil
}

//QueryHistoryKDataPage
// code format: sh.600000
func (bs *BaoStock) QueryHistoryKDataPage(curPageNum int, perPageCount int, code string, fields string,
	startDate time.Time, endDate time.Time, frequency string, adjustFlag string) (*ResultSet, error) {
	if len(code) != StockCodeLength {
		return nil, errors.New("invalid code, the format should be: sh.6000000")
	}
	if len(fields) == 0 {
		return nil, errors.New("fields cannot be empty")
	}
	if startDate.After(endDate) {
		return nil, errors.New("startDate large than endDate")
	}
	if len(frequency) == 0 {
		return nil, errors.New("frequency cannot be empty")
	}
	if len(adjustFlag) == 0 {
		return nil, errors.New("adjustFlag cannot be empty")
	}

	code = strings.ToLower(code)

	parts := []string{
		"query_history_k_data", bs.UserID, strconv.Itoa(curPageNum), strconv.Itoa(perPageCount), code,
		fields, util.DateToStr(startDate), util.DateToStr(endDate), frequency, adjustFlag,
	}
	msgBody := strings.Join(parts, MessageSplit)
	rs := &ResultSet{
		MsgType:      MessageTypeGetKDataRequest,
		ReqBodyParts: parts,
		Fields:       strings.Split(fields, ","),
		BS:           bs,
	}

	respMsg, err := bs.request(MessageTypeGetKDataRequest, msgBody)
	if err != nil {
		return nil, errors.Wrap(err, "get kdata failed")
	}
	if !respMsg.IsSucceed() {
		return nil, errors.Errorf("error code [%s], error message [%s]", respMsg.ErrCode, respMsg.ErrMsg)
	}

	if len(respMsg.BodyAttrs) < 7 {
		return nil, errors.Errorf("invalid body attrs [%s]", respMsg.BodyAttrs)
	}
	rs.CurPageNum, _ = strconv.Atoi(respMsg.BodyAttrs[4])
	rs.RespMsg = respMsg
	rs.setData(respMsg.BodyAttrs[6])
	return rs, nil
}

// GetDailyKData wrap QueryHistoryKDataPage to get daily k data
func (bs *BaoStock) GetDailyKData(code string, startDay time.Time, endDay time.Time) (*ResultSet, error) {
	logger.SugarLog.Debugf("startDay = %v, endDay = %v", startDay, endDay)
	return bs.QueryHistoryKDataPage(1, 200, code,
		"date,open,high,low,close,preclose,volume,amount,pctChg,peTTM,pbMRQ", startDay, endDay, "d",
		"2")
}
