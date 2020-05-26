package util

import (
	"io/ioutil"

	"github.com/jinzhu/gorm"

	"github.com/binance-chain/bep3-deputy/store"
)

func GetTestConfig() *Config {
	config := ParseConfigFromFile("../config/test_config_eth.json")
	return config
}

func PrepareDB(config *Config) (*gorm.DB, error) {
	config.DBConfig.DBPath = "tmp.db"
	tmpDBFile, err := ioutil.TempFile("", config.DBConfig.DBPath)
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(config.DBConfig.Dialect, tmpDBFile.Name())
	if err != nil {
		return nil, err
	}
	store.InitTables(db)
	return db, nil
}
