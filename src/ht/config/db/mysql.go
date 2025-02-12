package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ht-crm/autoconfigure"
	"time"
)

var Mysql *gorm.DB

func init() {
	var err error
	Mysql, err = gorm.Open(mysql.Open(GetDataSourceName()), gormConfig)
	if err != nil {
		panic("failed to connect database")
	}
	dbPoll, err := Mysql.DB()
	dbPoll.SetMaxIdleConns(10)
	dbPoll.SetMaxOpenConns(100)
	dbPoll.SetConnMaxLifetime(time.Hour)
}

func GetDataSourceName() string {
	mysqlConfig := autoconfigure.GetAppCig().Datasource.Mysql
	if mysqlConfig.Host == "" || mysqlConfig.Port == "" || mysqlConfig.User == "" || mysqlConfig.Password == "" || mysqlConfig.Database == "" {
		panic("mysql配置错误，请检查mysql配置")
	}
	return mysqlConfig.User + ":" + mysqlConfig.Password + "@tcp(" + mysqlConfig.Host + ":" + mysqlConfig.Port + ")/" + mysqlConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
}
