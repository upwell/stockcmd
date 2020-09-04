package hq

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func ConvertGB2UTF8(raw string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(raw)), simplifiedchinese.GB18030.NewDecoder())
	rawBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", errors.Wrapf(err, "failed to convert GB18030 to UTF8")
	}
	return string(rawBytes), nil
}

func ConvertGBK2UTF8(raw string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(raw)), simplifiedchinese.GBK.NewDecoder())
	rawBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", errors.Wrapf(err, "failed to convert GB18030 to UTF8")
	}
	return string(rawBytes), nil
}

// ConvertCode convert "sh.000001" to sina accept code "sh000001"
func ConvertCode(code string) string {
	return strings.ReplaceAll(code, ".", "")
}

// ConvertCodeBack convert "sh000001" to sina accept code "sh.000001"
func ConvertCodeBack(code string) string {
	if len(code) < 2 {
		return code
	}
	return fmt.Sprintf("%s.%s", code[:2], code[2:])
}
