package logger

import (
	"go.uber.org/zap"
	"hehan.net/my/stockcmd/config"
)

var Log *zap.Logger
var SugarLog *zap.SugaredLogger

func InitLogger() {
	if config.Verbose {
		Log, _ = zap.NewDevelopment()
	} else {
		Log, _ = zap.NewProduction()
	}
	SugarLog = Log.Sugar()
}
