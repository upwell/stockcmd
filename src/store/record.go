package store

import (
	"bytes"
	"context"
	"fmt"
	"strings"
	"time"

	"hehan.net/my/stockcmd/logger"

	"github.com/pkg/errors"

	"github.com/rocketlaunchr/dataframe-go/imports"

	"github.com/rocketlaunchr/dataframe-go"

	"go.etcd.io/bbolt"
)

type Record struct {
	Code string
	T    time.Time
	Val  string
}

var ErrDBColNotMatch = errors.New("columns not match the expected fields, db schema might change")

// the key is designed to be like:
//     sh000001#2000-01-01T00:00:00Z
// format: code#RFC3389 encoded time keys
func genKey(code string, t time.Time) string {
	return fmt.Sprintf("%s#%s", code, t.Format(time.RFC3339))
}

func RecreateDailyBucket() {
	DB.Update(func(tx *bbolt.Tx) error {
		tx.DeleteBucket([]byte(DailyBucketName))
		tx.CreateBucket([]byte(DailyBucketName))
		return nil
	})
}

func WriteRecord(code string, t time.Time, val string) {
	DB.Update(func(tx *bbolt.Tx) error {
		err := tx.Bucket([]byte(DailyBucketName)).Put([]byte(genKey(code, t)), []byte(val))
		return err
	})
}

func DeleteCodeRecords(code string) {
	DB.Batch(func(tx *bbolt.Tx) error {
		c := tx.Bucket([]byte(DailyBucketName)).Cursor()
		prefix := []byte(code)
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			c.Delete()
		}
		return nil
	})
}

func WriteRecords(records []*Record) {
	DB.Update(func(tx *bbolt.Tx) error {
		for _, record := range records {
			tx.Bucket([]byte(DailyBucketName)).Put([]byte(genKey(record.Code, record.T)), []byte(record.Val))
		}
		return nil
	})
}

func GetLastTime(code string) time.Time {
	var t time.Time
	DB.View(func(tx *bbolt.Tx) error {
		c := tx.Bucket([]byte(DailyBucketName)).Cursor()
		prefix := []byte(code)
		var key string
		for k, _ := c.Seek(prefix); k != nil && bytes.HasPrefix(k, prefix); k, _ = c.Next() {
			key = string(k)
		}
		if len(key) > 0 {
			t, _ = time.Parse(time.RFC3339, strings.Split(key, "#")[1])
		}
		return nil
	})
	return t
}

func GetRecords(code string, start time.Time, end time.Time) (*dataframe.DataFrame, error) {
	df := &dataframe.DataFrame{}
	if start.After(end) {
		return df, nil
	}

	var dbErr error

	DB.View(func(tx *bbolt.Tx) error {
		FieldsStr := "date,open,high,low,close,preclose,volume,amount,pctChg,peTTM,pbMRQ"
		c := tx.Bucket([]byte(DailyBucketName)).Cursor()
		startKey := []byte(genKey(code, start))
		endKey := []byte(genKey(code, end))
		csvRows := make([]string, 0, 1024)
		csvRows = append(csvRows, FieldsStr)

		ctx := context.Background()
		dataTypes := make(map[string]interface{})
		for _, field := range strings.Split(FieldsStr, ",") {
			dataTypes[field] = float64(0)
		}
		for k, v := c.Seek(startKey); k != nil && bytes.Compare(k, endKey) <= 0; k, v = c.Next() {
			csvRow := string(v)
			if len(strings.Split(csvRow, ",")) != len(strings.Split(FieldsStr, ",")) {
				dbErr = ErrDBColNotMatch
			}

			csvRows = append(csvRows, string(v))
		}
		dataTypes["date"] = imports.Converter{
			ConcreteType: time.Time{},
			ConverterFunc: func(in interface{}) (i interface{}, err error) {
				return time.Parse("2006-01-02", in.(string))
			},
		}
		opts := imports.CSVLoadOptions{
			TrimLeadingSpace: true,
			LargeDataSet:     false,
			DictateDataType:  dataTypes,
			InferDataTypes:   false,
		}
		csvStr := strings.Join(csvRows, "\n")
		retDf, err := imports.LoadFromCSV(ctx, strings.NewReader(csvStr), opts)
		if err != nil {
			logger.SugarLog.Errorf("load csv error [%v] for [%s]", err, code)
			df = &dataframe.DataFrame{}
			return nil
		}
		retDf.Sort(ctx, []dataframe.SortKey{
			{Key: "date", Desc: true},
		})
		df = retDf
		return nil
	})

	if dbErr != nil {
		return nil, dbErr
	}

	return df, nil
}
