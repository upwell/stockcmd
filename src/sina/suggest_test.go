package sina

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuggest(t *testing.T) {
	ret := Suggest("sz002475")
	if len(ret) == 0 {
		t.Error("failed to get suggest result, empty")
		return
	}
	assert.Equal(t, ret[0]["code"], "sz.002475")
	assert.Equal(t, ret[0]["name"], "立讯精密")
}
