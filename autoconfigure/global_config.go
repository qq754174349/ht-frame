// Package config 全局配置文件
package autoconfigure

import (
	"github.com/spf13/viper"
	"log"
)

const (
	defaultConfigFileName = "config"
	defaultConfigFileType = "yaml"
)

var (
	appCfg *AppConfig
)

type AppConfig struct {
	Active     string
	AppName    string `yaml:"app_name" json:"app_name" mapstructure:"app_name"`
	Web        Web
	Log        LogConfig
	Datasource datasource
}

type datasource struct {
	Mysql MysqlConfig
	Redis RedisConfig
}

type Web struct {
	Port string
}

type LogConfig struct {
	Level       string
	OutputPaths string `json:"output_paths" yaml:"output_paths" mapstructure:"output_paths"`
}

type MysqlConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type RedisConfig struct {
	Addr     string
	User     string
	Password string
	DB       int
}

func init() {
	InitConfig("")
}

func InitConfig(active string) {
	viper.AddConfigPath("config/")
	viper.SetConfigType(defaultConfigFileType)
	viper.SetConfigName(defaultConfigFileName + "." + defaultConfigFileType)
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	if active == "" {
		activeVal := viper.Get("active")
		if activeVal == nil {
			log.Fatalf("没有激活的配置")
		}
		active = activeVal.(string)
	}
	log.Printf("Active environment:%s", active)

	viper.SetConfigName(defaultConfigFileName + "-" + active + "." + defaultConfigFileType)
	// 读取配置文件
	if err := viper.MergeInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	appCfg = &AppConfig{}
	err := viper.Unmarshal(appCfg)
	if err != nil {
		log.Fatal("配置文件格式错误")
	}
}

func GetAppCig() *AppConfig {
	return appCfg
}
