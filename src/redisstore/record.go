package redisstore

import (
	"context"
	"fmt"
	"strings"
	"time"

	"hehan.net/my/stockcmd/store"

	"github.com/go-redis/redis"

	"hehan.net/my/stockcmd/util"

	"github.com/pkg/errors"

	"github.com/rocketlaunchr/dataframe-go"
	"github.com/rocketlaunchr/dataframe-go/imports"
	"hehan.net/my/stockcmd/logger"
)

const prefixDaily = "sc:daily:"
const zPrefixDaily = "sc:daily:z:"

var ErrDBColNotMatch = errors.New("columns not match the expected fields, db schema might change")

func genKey(code string) string {
	return fmt.Sprintf("%s%s", prefixDaily, code)
}

func genZKey(code string) string {
	return fmt.Sprintf("%s%s", zPrefixDaily, code)
}

func RecreateDailyBucket() {
	iter := Redis.Scan(0, prefixDaily+"*", 0).Iterator()
	for iter.Next() {
		Redis.Del(iter.Val())
	}
	iter = Redis.Scan(0, zPrefixDaily+"*", 0).Iterator()
	for iter.Next() {
		Redis.Del(iter.Val())
	}
}

func WriteRecord(code string, t time.Time, val string) {
	Redis.HSet(genKey(code), util.DateToStr2(t), val)
	Redis.ZAdd(genZKey(code), redis.Z{Member: val, Score: float64(t.Unix())})
}

func DeleteCodeRecords(code string) {
	Redis.Del(genKey(code))
	Redis.Del(genZKey(code))
}

func WriteRecords(records []*store.Record) {
	for _, record := range records {
		WriteRecord(record.Code, record.T, record.Val)
	}
}

func GetLastTime(code string) time.Time {
	var t time.Time
	max, _ := Redis.ZRevRangeWithScores(genZKey(code), 0, 0).Result()
	if len(max) == 0 {
		return t
	}

	tUnix := int64(max[0].Score)
	return time.Unix(tUnix, 0)
}

func GetRecords(code string, start time.Time, end time.Time) (*dataframe.DataFrame, error) {
	//defer util.MeasureTime(code + "GetRecords")()

	df := &dataframe.DataFrame{}
	if start.After(end) {
		return df, nil
	}
	var dbErr error

	records, _ := Redis.ZRangeByScoreWithScores(genZKey(code), redis.ZRangeBy{
		Min: fmt.Sprint(start.Unix()),
		Max: fmt.Sprint(end.Unix()),
	}).Result()

	FieldsStr := "date,open,high,low,close,preclose,volume,amount,pctChg,peTTM,pbMRQ"
	csvRows := make([]string, 0, 1024)
	csvRows = append(csvRows, FieldsStr)
	dataTypes := make(map[string]interface{})
	for _, field := range strings.Split(FieldsStr, ",") {
		dataTypes[field] = float64(0)
	}
	ctx := context.Background()
	for _, record := range records {
		csvRow := record.Member.(string)
		if len(strings.Split(csvRow, ",")) != len(strings.Split(FieldsStr, ",")) {
			dbErr = ErrDBColNotMatch
		}
		csvRows = append(csvRows, csvRow)
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
		return df, errors.Wrapf(err, "load csv error for [%s]", code)
	}
	retDf.Sort(ctx, []dataframe.SortKey{
		{Key: "date", Desc: true},
	})
	df = retDf

	if dbErr != nil {
		return nil, dbErr
	}

	return df, nil
}
