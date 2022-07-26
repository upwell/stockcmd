package stat

import (
	"testing"
	"time"

	"hehan.net/my/stockcmd/logger"
	"hehan.net/my/stockcmd/store"
)

func TestGetRPS(t *testing.T) {
	logger.InitLogger()
	basics := store.GetBasics()
	GetRPS(basics, 10, time.Time{})
}
