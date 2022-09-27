package eastmoney

import (
	"testing"

	"hehan.net/my/stockcmd/logger"
)

func TestHQApi_GetAllHQ(t *testing.T) {
	logger.InitLogger()
	api := HQApi{}
	result, err := api.GetAllHQ()
	if err != nil {
		t.Errorf("failed with [%v]", err)
	}
	if len(result) == 0 {
		t.Errorf("empty result")
	}

	println(result[0].Code)
	println(result[0].Price)
}

func TestHQApi_GetHQ(t *testing.T) {
	logger.InitLogger()
	api := HQApi{}
	result, err := api.GetHQ("sh.603530")
	if err != nil {
		t.Errorf("%v", err)
	}

	println(result.Now)
}
