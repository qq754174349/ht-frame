package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ht-crm/src/ht/config"
	"time"
)

var Mysql *gorm.DB

func init() {
	envCfg := config.GetEnvCfg()
	var err error
	Mysql, err = gorm.Open(mysql.Open(envCfg.Mysql.GetDataSourceName()), gormConfig)
	if err != nil {
		panic("failed to connect database")
	}
	dbPoll, err := Mysql.DB()
	dbPoll.SetMaxIdleConns(10)
	dbPoll.SetMaxOpenConns(100)
	dbPoll.SetConnMaxLifetime(time.Hour)
}
