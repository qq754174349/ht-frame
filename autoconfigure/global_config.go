// Package config 全局配置文件
package autoconfigure

import (
	"github.com/qq754174349/ht-frame/config"
	"github.com/qq754174349/ht-frame/logger"
	"github.com/qq754174349/ht-frame/web"
	"github.com/spf13/viper"
	"log"
)

const (
	defaultConfigFileName = "config"
	defaultConfigFileType = "yaml"
)

var (
	appCfg       *AppConfig
	initializers = make(map[string]config.Configuration)
)

type AppConfig struct {
	Active     string
	AppName    string `yaml:"app_name" json:"app_name" mapstructure:"app_name"`
	Web        web.Web
	Log        logger.LogConfig
	Datasource datasource
}

type datasource struct {
	Mysql MysqlConfig
	Redis RedisConfig
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

func Register(name string, conf config.Configuration) {
	initializers[name] = conf
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

	autoConfigure()
	logger.InitLogger(appCfg.Log)
}

func autoConfigure() {

}

func GetAppCig() *AppConfig {
	return appCfg
}
