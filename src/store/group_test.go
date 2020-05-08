package store

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestAddGroup(t *testing.T) {
	name := "hold"
	AddGroup(name)
	isExist := CheckGroupExist(name)
	assert.Equal(t, isExist, true)
}

func TestGetGroup(t *testing.T) {
	name := "hold"
	AddGroup(name)
	g := GetGroup(name)
	assert.Equal(t, g.Name, name)
}

func TestGroup_AddStock(t *testing.T) {
	name := "hold"
	AddGroup(name)
	g := GetGroup(name)
	g.AddStock("sz002475", "立讯精密")
	fmt.Println(g)
}
