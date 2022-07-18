package akshare

import (
	mapset "github.com/deckarep/golang-set"
)

var shCodeSet mapset.Set
var szCodeSet mapset.Set
var bjCodeSet mapset.Set

func init() {
	shCodeSet = mapset.NewSet("600", "601", "603", "605", "688", "689")
	szCodeSet = mapset.NewSet("000", "001", "002", "003", "300", "301")
	bjCodeSet = mapset.NewSet("43", "83", "87", "88")
}

func DetermineExchangeByCode(code string) string {
	prefix := code[0:3]
	if shCodeSet.Contains(prefix) {
		return "sh"
	} else if szCodeSet.Contains(prefix) {
		return "sz"
	} else {
		prefix = code[0:2]
		if bjCodeSet.Contains(prefix) {
			return "bj"
		}
	}
	return ""
}
