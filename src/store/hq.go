package store

import (
	"strconv"

	"go.etcd.io/bbolt"
)

type StockHQ struct {
	Code  string `json:"code"`
	Price string `json:"price"`
}

func WriteHQ(hq *StockHQ) error {
	DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(HQBucketName))
		b.Put([]byte(hq.Code), []byte(hq.Price))
		return nil
	})
	return nil
}

func BulkWriteHQ(hqs []*StockHQ) error {
	DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(HQBucketName))
		for _, hq := range hqs {
			b.Put([]byte(hq.Code), []byte(hq.Price))
		}
		return nil
	})
	return nil
}

func GetHQ(code string) float64 {
	var ret float64
	DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(HQBucketName))
		v := b.Get([]byte(code))
		floatV, err := strconv.ParseFloat(string(v), 64)
		ret = floatV
		return err
	})
	return ret
}
