package db

import (
	"github.com/qq754174349/ht-frame/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"time"
)

var Mysql *gorm.DB

type AutoConfig struct{}

var gormConfig = &gorm.Config{
	NamingStrategy: schema.NamingStrategy{
		SingularTable: true,
	},
}

func (AutoConfig) Init(cfg config.AppConfig) error {
	mysqlConfig := cfg.Datasource.Mysql
	var err error
	Mysql, err = gorm.Open(mysql.Open(GetDataSourceName(mysqlConfig)), gormConfig)
	if err != nil {
		panic("failed to connect database")
	}
	dbPoll, err := Mysql.DB()
	dbPoll.SetMaxIdleConns(10)
	dbPoll.SetMaxOpenConns(100)
	dbPoll.SetConnMaxLifetime(time.Hour)
	return nil
}

func GetDataSourceName(mysqlConfig config.MysqlConfig) string {
	if mysqlConfig.Host == "" || mysqlConfig.Port == "" || mysqlConfig.User == "" || mysqlConfig.Password == "" || mysqlConfig.Database == "" {
		panic("mysql配置错误，请检查mysql配置")
	}
	return mysqlConfig.User + ":" + mysqlConfig.Password + "@tcp(" + mysqlConfig.Host + ":" + mysqlConfig.Port + ")/" + mysqlConfig.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
}
