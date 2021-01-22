package sina

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestGetLivePrice(t *testing.T) {
	//code := "sz.300284"
	code := "sh.000911"
	//code := "sh.000009"
	api := HQApi{}
	ret, err := api.GetHQ(code)
	if err != nil {
		t.Errorf("failed with [%v]", err)
		return
	}
	fmt.Println(ret)
}

func TestGetLivePrices(t *testing.T) {
	ret := GetLivePrices([]string{"sz002475", "sz300433"})
	assert.Equal(t, len(ret), 2)
	fmt.Println(ret)
}
