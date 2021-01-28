package application

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"go-api-samp/util/config"
)

var (
	sqlDB *sql.DB
)

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
	//db.SetConnMaxIdleTime(conf.ConnMaxIdleTime)
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)

	//time.Sleep(time.Second * 20) // for docker-compose
	if err = db.Ping(); err != nil {
		fmt.Println("ping error")
		return err
	}

	sqlDB = db

	return nil
}
