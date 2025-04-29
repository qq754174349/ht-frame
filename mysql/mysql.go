package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/qq754174349/ht-frame/autoconfigure"
	"github.com/qq754174349/ht-frame/logger"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
	"time"
)

var (
	mysqlInstances sync.Map // 存储所有MySQL实例
	defaultName    string
	gormConfig     = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		PrepareStmt: true,
	}
	config *Mysql
)

type Mysql struct {
	Mysql map[string]Config `yaml:"mysql" mapstructure:"mysql"`
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

type AutoConfig struct{}

func init() {
	err := autoconfigure.Register(AutoConfig{})
	if err != nil {
		logger.Fatal("mysql 自动配置注册失败")
	}
}

func (AutoConfig) Init() error {
	config = &Mysql{}
	err := viper.Unmarshal(config)
	if err != nil {
		logger.Fatal("配置文件格式错误")
	}
	first := true
	for name, mysqlCfg := range config.Mysql {
		if err := initMySQL(mysqlCfg, name); err != nil {
			return fmt.Errorf("MySQL[%s]初始化失败: %v", name, err)
		}
		if first {
			defaultName = name
			first = false
		}
	}

	if defaultName == "" {
		return fmt.Errorf("至少需要配置一个MySQL数据源")
	}

	return nil
}

func initMySQL(cfg Config, name string) error {
	dsn := buildDSN(cfg)
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return fmt.Errorf("连接MySQL[%s]失败: %v", name, err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取连接池失败: %v", err)
	}

	// 配置连接池
	configurePool(sqlDB, cfg)

	// 健康检查
	if err := ping(sqlDB); err != nil {
		return fmt.Errorf("MySQL[%s]健康检查失败: %v", name, err)
	}

	mysqlInstances.Store(name, db)
	go monitor(name, db, cfg)
	logger.Infof("MySQL[%s]已初始化 @ %s:%s", name, cfg.Host, cfg.Port)
	return nil
}

func Get(name ...string) (*gorm.DB, error) {
	var instancesName string
	if len(name) == 0 {
		instancesName = defaultName
	} else {
		instancesName = name[0]
	}
	val, ok := mysqlInstances.Load(instancesName)
	if !ok {
		return nil, fmt.Errorf("MySQL[%s]未初始化", name)
	}

	db := val.(*gorm.DB)
	if err := verify(db); err != nil {
		return nil, fmt.Errorf("MySQL[%s]连接异常: %v", name, err)
	}
	return db, nil
}

// 内部工具函数
func buildDSN(cfg Config) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=5s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

func configurePool(sqlDB *sql.DB, cfg Config) {
	maxIdle := 10

	sqlDB.SetMaxIdleConns(maxIdle)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}

func ping(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}

func monitor(name string, db *gorm.DB, cfg Config) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if err := verify(db); err != nil {
			logger.Warnf("MySQL[%s]连接异常: %v", name, err)
			_ = reconnect(name, cfg)
		}
	}
}

func reconnect(name string, cfg Config) error {
	if val, ok := mysqlInstances.Load(name); ok {
		if db, ok := val.(*gorm.DB); ok {
			if sqlDB, err := db.DB(); err == nil {
				_ = sqlDB.Close()
			}
		}
		mysqlInstances.Delete(name)
	}
	return initMySQL(cfg, name)
}

// verify 验证数据库连接是否健康
func verify(db *gorm.DB) error {
	// 获取底层sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取连接池失败: %v", err)
	}

	// 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 执行Ping测试
	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("连接不可用: %v", err)
	}

	// 可选：检查数据库版本（增强验证）
	//var version string
	//if err := mysql.Raw("SELECT VERSION()").Scan(&version).Error; err != nil {
	//	return fmt.Errorf("数据库版本检查失败: %v", err)
	//}
	//
	//logger.Debugf("数据库连接正常，版本: %s", version)
	return nil
}

func CloseAllMySQL() {
	mysqlInstances.Range(func(key, value interface{}) bool {
		if db, ok := value.(*gorm.DB); ok {
			if sqlDB, err := db.DB(); err == nil {
				_ = sqlDB.Close()
				logger.Infof("MySQL[%v]已关闭", key)
			}
		}
		return true
	})
}
