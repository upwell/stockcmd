package store

import (
	"bytes"
	"encoding/gob"

	"go.etcd.io/bbolt"
)

/*
以 key / value 的形式记录配置或者过程数据
*/

type ConfigStore struct {
	bucketName string
}

var RunningConfig ConfigStore

// Marshal encodes a Go value to gob.
func gobMarshal(v interface{}) ([]byte, error) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(v)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

// Unmarshal decodes a gob value into a Go value.
func gobUnmarshal(data []byte, v interface{}) error {
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)
	return decoder.Decode(v)
}

func (store ConfigStore) GetString(key string) (value string, existed bool, err error) {
	val := new(string)

	err, existed = store.Get(key, val)
	if err != nil {
		return
	}
	return *val, existed, nil
}

func (store ConfigStore) Get(key string, val interface{}) (err error, exist bool) {
	var data []byte
	err = DB.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(store.bucketName))
		v := b.Get([]byte(key))
		if v != nil {
			data = make([]byte, len(v))
			copy(data, v)
		}
		return nil
	})

	if err != nil {
		return nil, false
	}
	if data == nil {
		return nil, false
	}

	return gobUnmarshal(data, val), true
}

func (store ConfigStore) GetInt64(key string) (value int64, existed bool, err error) {
	val := new(int64)

	err, existed = store.Get(key, val)
	if err != nil {
		return
	}
	return *val, existed, nil
}

func (store ConfigStore) GetInt64OrDefault(key string, defaultVal int64) (value int64, err error) {
	value, existed, err := store.GetInt64(key)
	if err != nil && existed {
		value = defaultVal
	}
	return
}

func (store ConfigStore) Set(key string, val interface{}) error {
	data, err := gobMarshal(val)
	if err != nil {
		return err
	}

	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(store.bucketName))
		return b.Put([]byte(key), data)
	})
}

func (store ConfigStore) Delete(key string) error {
	return DB.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(store.bucketName))
		return b.Delete([]byte(key))
	})
}

func init() {
	RunningConfig = ConfigStore{ConfigBucketName}
}
