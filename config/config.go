package config

const (
	Logger = "log"
	WEB    = "web"
	MYSQL  = "db"
	REDIS  = "redis"
	CONSUL = "consul"
)

var appCfg *AppConfig

type AppConfig struct {
	Active     string
	AppName    string `yaml:"app_name" json:"app_name" mapstructure:"app_name"`
	Web        WebConfig
	Log        LogConfig
	Datasource datasource
	Consul     ConsulConfig
}

type LogConfig struct {
	Level       string
	OutputPaths string `json:"output_paths" yaml:"output_paths" mapstructure:"output_paths"`
}

type WebConfig struct {
	Port string
}

type datasource struct {
	Mysql map[string]MysqlConfig `yaml:"mysql" mapstructure:"mysql"`
	Redis map[string]RedisConfig `yaml:"redis" mapstructure:"redis"`
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

type ConsulConfig struct {
	Addr string
}

type Configuration interface {
	Init(config *AppConfig) error
}

func SetAppCfg(appConfig *AppConfig) {
	appCfg = appConfig
}

func GetAppCfg() *AppConfig {
	return appCfg
}
