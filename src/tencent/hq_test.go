package tencent

import (
	"fmt"
	"testing"
)

func TestTencentHQApi_GetHQ(t *testing.T) {
	api := HQApi{}
	ret, err := api.GetHQ("sh.600036")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ret)
}
