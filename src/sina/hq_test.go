package sina

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestGetLivePrice(t *testing.T) {
	ret, err := GetLivePrice("sz002475")
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
