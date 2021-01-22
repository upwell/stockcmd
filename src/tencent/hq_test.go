package tencent

import (
	"fmt"
	"testing"
)

func TestTencentHQApi_GetHQ(t *testing.T) {
	api := HQApi{}
	ret, err := api.GetHQ("sh.000939")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ret.Now)
}
