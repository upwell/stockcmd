package store

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestWriteBasic(t *testing.T) {
	code := "sz002475"
	basic := &StockBasic{
		Code: code,
		Name: "立讯精密",
	}
	WriteBasic(code, basic)
}

func TestGetBasic(t *testing.T) {
	code := "sz002475"
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
	code := "sz002049"
	name := GetName(code, true)
	assert.Equal(t, name, "紫光国微")
}
