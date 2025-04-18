package logger

import (
	"github.com/qq754174349/ht-frame/autoconfigure"
	"github.com/qq754174349/ht-frame/config"
	"io"
)

type AutoConfig struct{}

// Logger 是日志门面接口
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})
	WithTraceID(traceId string) Logger               // 返回一个新的 Logger 实例，附带 TraceID
	WithFields(fields map[string]interface{}) Logger // 返回一个新的 Logger 实例，附带字段
	RedirectStdLog()
	Writer() io.Writer
}

const (
	defaultLevel       = "info"
	defaultOutputPaths = "tmp/logs/"
)

var log Logger

func init() {
	autoconfigure.Register(AutoConfig{})
}

func (AutoConfig) Init(cfg *config.AppConfig) error {
	logConfig := cfg.Log
	InitLogger(logConfig)
	return nil
}

func InitLogger(config config.LogConfig) {
	if config.Level == "" {
		config.Level = defaultLevel
	}
	if config.OutputPaths == "" {
		config.OutputPaths = defaultOutputPaths
	}

	log = newZapLog(&config)
	// 替换标准日志
	log.RedirectStdLog()
}

func Debug(args ...interface{})                   { log.Debug(args...) }
func Info(args ...interface{})                    { log.Info(args...) }
func Warn(args ...interface{})                    { log.Warn(args...) }
func Error(args ...interface{})                   { log.Error(args...) }
func Fatal(args ...interface{})                   { log.Fatal(args...) }
func Debugf(template string, args ...interface{}) { log.Debugf(template, args...) }
func Infof(template string, args ...interface{})  { log.Infof(template, args...) }
func Warnf(template string, args ...interface{})  { log.Warnf(template, args...) }
func Errorf(template string, args ...interface{}) { log.Errorf(template, args...) }
func Fatalf(template string, args ...interface{}) { log.Fatalf(template, args...) }

func WithTraceID(traceId string) Logger {
	return log.WithTraceID(traceId)
}

func WithFields(fields map[string]interface{}) Logger {
	return log.WithFields(fields)
}

func Writer() io.Writer {
	return log.Writer()
}
