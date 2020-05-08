package store

import (
	"encoding/json"

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
		sinaCode := sina.ConvertCode(code)
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
