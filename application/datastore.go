package application

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-api-samp/util/config"
	"go-api-samp/util/log"
)

/*
 * アプリケーションで扱うdatastore
 * DBインスタンス作成取得、locationMapの保持
 */

var (
	sqlDB        *sql.DB
	locationsMap map[int]string
)

func GetLocationsMap() map[int]string {
	return locationsMap
}

func NewLocationsMap() {
	locationsMap = make(map[int]string)
}

func GetDB() *sql.DB {
	return sqlDB
}

func NewDBOpen(conf *config.DBConfig) error {
	dataSource := fmt.Sprintf("%s:%s@%s/%s", conf.User, conf.Password, conf.Addr, conf.DBName)
	db, err := sql.Open(conf.DriverName, dataSource)
	if err != nil {
		fmt.Println("open error")
		return err
	}

	db.SetConnMaxLifetime(conf.ConnMaxLifeTime)
	//db.SetConnMaxIdleTime(conf.ConnMaxIdleTime) //MaxIdleConnsがあれば不要
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)

	if err = db.Ping(); err != nil {
		log.GetLogger().Error(nil, "ping error", err)
		return err
	}

	sqlDB = db

	return nil
}
