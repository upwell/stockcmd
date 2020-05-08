package sina

import (
	"strings"

	"hehan.net/my/stockcmd/logger"

	"github.com/levigross/grequests"
)

func isValidCode(code string) bool {
	return strings.Contains(code, "sh") || strings.Contains(code, "sz")
}

func Suggest(input string) []map[string]string {
	ret := make([]map[string]string, 0, 1024)
	if len(input) == 0 {
		return ret
	}

	resp, err := grequests.Get(SuggestURL+input, nil)
	if err != nil {
		return ret
	}
	if resp.StatusCode != 200 {
		return ret
	}

	rawResult, err := ConvertGB2UTF8(resp.String())
	if err != nil {
		return ret
	}
	rawResult = strings.Split(rawResult, "=")[1]
	rawResult = strings.ReplaceAll(rawResult, "\"", "")
	parts := strings.Split(rawResult, ";")
	if len(parts) == 0 {
		return ret
	}
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		seps := strings.Split(part, ",")
		if !isValidCode(seps[3]) {
			logger.SugarLog.Debugf("ignore, not valid sh or sz code [%s/%s]", seps[4], seps[3])
			continue
		}
		ret = append(ret, map[string]string{
			"name": seps[4],
			"code": ConvertCodeBack(seps[3]),
		})
	}
	return ret
}
