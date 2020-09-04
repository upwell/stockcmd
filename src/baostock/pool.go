package baostock

import (
	"errors"
	"time"

	"hehan.net/my/stockcmd/logger"

	"github.com/silenceper/pool"
)

var BSPool pool.Pool

func factory() (interface{}, error) {
	bs := NewBaoStockInstance()
	if bs != nil {
		err := bs.Login()
		if err != nil {
			return nil, err
		}
		return bs, nil
	}
	return bs, errors.New("failed to init baostock instance")
}

func close(v interface{}) error {
	bs := v.(*BaoStock)
	bs.Logout()

	return bs.Conn.Close()
}

func init() {
	poolConfig := &pool.Config{
		InitialCap:  0,
		MaxIdle:     3,
		MaxCap:      4,
		Factory:     factory,
		Close:       close,
		IdleTimeout: 60 * time.Second,
	}
	p, err := pool.NewChannelPool(poolConfig)
	if err != nil {
		logger.SugarLog.Fatalf("failed to create baostock pool with [%v]", err)
	}
	BSPool = p
}
