package store

import (
	"os"
	"path/filepath"

	"hehan.net/my/stockcmd/logger"

	bolt "go.etcd.io/bbolt"
)

var DB *bolt.DB

const DailyBucketName = "Daily"
const BasicBucketName = "Basic"
const ConfigBucketName = "ConfigRunning"
const GroupBucketName = "Group"
const HQBucketName = "HQ"

var appPath string

const dbName = "data.db"

func init() {
	homeDir, _ := os.UserHomeDir()
	joinedPath := filepath.Join(homeDir, ".config", "stockcmd")
	appPath = joinedPath

	err := os.MkdirAll(appPath, os.ModePerm)

	db, err := bolt.Open(filepath.Join(appPath, dbName), 0600, nil)
	if err != nil {
		logger.SugarLog.Fatalf("failed to open bolt db error [%v]", err)
	}

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(DailyBucketName))
		if err != nil {
			logger.SugarLog.Fatalf("failed to get/create db bucket [%v]", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(BasicBucketName))
		if err != nil {
			logger.SugarLog.Fatalf("failed to get/create db bucket [%v]", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(ConfigBucketName))
		if err != nil {
			logger.SugarLog.Fatalf("failed to get/create db bucket [%v]", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(GroupBucketName))
		if err != nil {
			logger.SugarLog.Fatalf("failed to get/create db bucket [%v]", err)
		}

		_, err = tx.CreateBucketIfNotExists([]byte(HQBucketName))
		if err != nil {
			logger.SugarLog.Fatalf("failed to get/create db bucket [%v]", err)
		}
		return nil
	})

	DB = db
}
