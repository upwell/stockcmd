package baostock

import (
	"strconv"
	"strings"
	"time"
)

func paddingZeroForString(content string, length int, direction string) string {
	paddingStr := ""
	strLen := len(content)
	for i := 0; i < length-strLen; i++ {
		paddingStr += "0"
	}

	var result string
	switch direction {
	case "left":
		result = paddingStr + content
	case "right":
		result = content + paddingStr
	default:
		result = content
	}
	return result
}

func paddingLeftZero(content string, length int) string {
	return paddingZeroForString(content, length, "left")
}

func paddingRightZero(content string, length int) string {
	return paddingZeroForString(content, length, "right")
}

func toMessageHeader(msgType string, msgLen int) string {
	lenStr := paddingLeftZero(strconv.Itoa(msgLen), 10)
	return strings.Join([]string{ClientVersion, msgType, lenStr}, MessageSplit)
}

func dateToStr(time time.Time) string {
	return time.Format("2006-01-02")
}
