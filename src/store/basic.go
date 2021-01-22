package store

import (
	"encoding/json"
	"fmt"

	"hehan.net/my/stockcmd/logger"

	"hehan.net/my/stockcmd/hq"

	"hehan.net/my/stockcmd/sina"

	"go.etcd.io/bbolt"
)

type StockBasic struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func GetName(code string, force bool) string {
	var basic *StockBasic
	if !force {
		basic = GetBasic(code)
	}
	if basic == nil {
		logger.SugarLog.Infof("no basic info for [%s], fetch now", code)
		sinaCode := hq.ConvertCode(code)
		ret := sina.Suggest(sinaCode)
		if len(ret) == 0 {
			return ""
		}
		name := ret[0]["name"]
		WriteBasic(code, &StockBasic{
			Code: code,
			Name: name,
		})
		return name
	} else {
		return basic.Name
	}
}

func GetBasic(code string) *StockBasic {
	var ret *StockBasic
	DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(BasicBucketName))
		v := b.Get([]byte(code))
		if v != nil {
			var basic StockBasic
			err := json.Unmarshal(v, &basic)
			if err != nil {
				return err
			}
			ret = &basic
		}
		return nil
	})
	return ret
}

func WriteBasic(code string, basic *StockBasic) {
	DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(BasicBucketName))
		bytes, _ := json.Marshal(basic)
		b.Put([]byte(code), bytes)
		return nil
	})
}

func WriteBasics(arrs [][]string) {
	DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(BasicBucketName))
		for _, arr := range arrs {
			basic := &StockBasic{
				Code: arr[0],
				Name: arr[2],
			}
			if arr[1] == "0" {
				fmt.Printf("%s %s\n", arr[0], arr[2])
			}
			bytes, _ := json.Marshal(basic)
			b.Put([]byte(basic.Code), bytes)
		}
		return nil
	})
}

func GetCodes() []string {
	ret := make([]string, 0, 512)
	DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(BasicBucketName))
		c := b.Cursor()
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			ret = append(ret, string(k))
		}
		return nil
	})
	return ret
}

func GetBasics() []*StockBasic {
	ret := make([]*StockBasic, 0, 512)
	DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(BasicBucketName))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			basic := &StockBasic{}
			if err := json.Unmarshal(v, basic); err != nil {
				return err
			}
			ret = append(ret, basic)
		}
		return nil
	})
	return ret
}

func RecreateBasicBucket() {
	DB.Update(func(tx *bbolt.Tx) error {
		tx.DeleteBucket([]byte(BasicBucketName))
		tx.CreateBucket([]byte(BasicBucketName))
		return nil
	})
}
