// Package autoconfigure Package config 全局配置文件
package autoconfigure

import (
	"github.com/qq754174349/ht-frame/config"
	"github.com/spf13/viper"
	"log"
)

const (
	defaultConfigFileName = "config"
	defaultConfigFileType = "yaml"
)

var (
	appCfg       *config.AppConfig
	initializers = make(map[string]config.Configuration)
)

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

	appCfg = &config.AppConfig{}
	err := viper.Unmarshal(appCfg)
	if err != nil {
		log.Fatal("配置文件格式错误")
	}
	config.SetAppCfg(appCfg)

	autoConfigure()
}

func autoConfigure() {
	for _, v := range initializers {
		v.Init(config.GetAppCfg())
	}
}
