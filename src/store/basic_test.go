package store

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestWriteBasic(t *testing.T) {
	code := "sz.002475"
	basic := &StockBasic{
		Code: code,
		Name: "立讯精密",
	}
	WriteBasic(code, basic)
}

func TestGetBasic(t *testing.T) {
	code := "sz.002475"
	basic := &StockBasic{
		Code: code,
		Name: "立讯精密",
	}
	WriteBasic(code, basic)
	basic = GetBasic(code)
	if basic == nil {
		t.Error("failed to get basic")
		return
	}
	fmt.Println(basic)
}

func TestGetName(t *testing.T) {
	code := "sz.002049"
	name := GetName(code, true)
	assert.Equal(t, name, "紫光国微")
}

func TestGetBasics(t *testing.T) {
	ret := GetCodes()
	assert.Equal(t, ret[0], "sh.000001")
}
