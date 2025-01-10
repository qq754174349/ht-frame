// Package config 全局配置文件
package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

const (
	defaultConfigPath = "conf/config.json"
	defaultActive     = "dev"
)

var (
	appCfg     *appConfig
	currEvnCfg *EnvConfig
)

type appConfig struct {
	Active  string    `json:"active"`
	AppName string    `json:"app_name"`
	Dev     EnvConfig `json:"dev"`
	Test    EnvConfig `json:"test"`
	Prod    EnvConfig `json:"prod"`
}

type EnvConfig struct {
	Log   LogConfig   `json:"log"`
	Port  string      `json:"port"`
	Mysql mysqlConfig `json:"mysql"`
	Redis redisConfig `json:"redis"`
}

type LogConfig struct {
	Level   string `json:"level"`
	FileUrl string `json:"fileUrl"`
}

type mysqlConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type redisConfig struct {
	Addr     string `json:"addr"`
	User     string `json:"user"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func init() {
	appCfg = newAppConfig()

	if appCfg.AppName == "" {
		panic("请配置应用名称application_name")
	}

	active := appCfg.Active
	if active == "" {
		active = defaultActive
	}

	switch active {
	case "dev":
		currEvnCfg = &appCfg.Dev
	case "test":
		currEvnCfg = &appCfg.Test
	case "prod":
		currEvnCfg = &appCfg.Prod
	default:
		log.Fatalf("无效的 active 配置: %s", active)
	}

}

func (mysql mysqlConfig) GetDataSourceName() string {
	if mysql.Host == "" || mysql.Port == "" || mysql.User == "" || mysql.Password == "" || mysql.Database == "" {
		panic("mysql配置错误，请检查mysql配置")
	}
	return mysql.User + ":" + mysql.Password + "@tcp(" + mysql.Host + ":" + mysql.Port + ")/" + mysql.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
}

func newAppConfig() *appConfig {
	file, err := os.Open(defaultConfigPath)
	if err != nil {
		log.Fatalf("配置文件读取错误,%s", err.Error())
	}
	defer file.Close()

	all, err := io.ReadAll(file)
	if err != nil {
		log.Fatalf("配置文件读取错误,%s", err.Error())
	}

	var appConfig appConfig
	err = json.Unmarshal(all, &appConfig)
	if err != nil {
		log.Fatalf("配置文件格式错误,%s", err.Error())
	}
	return &appConfig
}

func GetEnvCfg() EnvConfig {
	return *currEvnCfg
}

func GetAppName() string {
	return appCfg.AppName
}

func GetActive() string {
	return appCfg.Active
}
