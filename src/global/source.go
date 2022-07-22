package global

import (
	"hehan.net/my/stockcmd/baostock"
	"hehan.net/my/stockcmd/base"
	"hehan.net/my/stockcmd/eastmoney"
	"hehan.net/my/stockcmd/sina"
	"hehan.net/my/stockcmd/store"
	"hehan.net/my/stockcmd/tencent"
)

func GetDataSource() base.DataSource {
	dailySource, _ := store.RunningConfig.GetStringOrDefault("dailySource", "baostock")
	switch dailySource {
	case "baostock":
		return baostock.BSP
	case "eastmoney":
		return eastmoney.EM
	default:
		return baostock.BSP
	}
}

func GetHQSource() base.HQApi {
	hqSource, _ := store.RunningConfig.GetStringOrDefault("hqSource", "tencent")
	switch hqSource {
	case "tencent":
		return tencent.HQApi{}
	case "sina":
		return sina.HQApi{}
	default:
		return tencent.HQApi{}
	}
}
