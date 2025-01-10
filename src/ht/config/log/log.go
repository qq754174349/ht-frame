package log

import (
	"ht-crm/src/ht/config"
	"io"
)

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	TraceDebug(traceId string, args ...interface{})
	TraceInfo(traceId string, args ...interface{})
	TraceWarn(traceId string, args ...interface{})
	TraceError(traceId string, args ...interface{})
	TraceDebugf(traceId string, template string, args ...interface{})
	TraceInfof(traceId string, template string, args ...interface{})
	TraceWarnf(traceId string, template string, args ...interface{})
	TraceErrorf(traceId string, template string, args ...interface{})
	RedirectStdLog()
	Writer() io.Writer
}

const (
	defaultLevel       = "info"
	defaultOutputPaths = "tmp/logs/"
)

var log Logger

func init() {
	appLogCfg := config.GetEnvCfg().Log
	if appLogCfg.Level == "" {
		appLogCfg.Level = defaultLevel
	}
	if appLogCfg.FileUrl == "" {
		appLogCfg.FileUrl = defaultOutputPaths
	}

	log = newZapLog(&appLogCfg)
	// 替换标准日志
	log.RedirectStdLog()
}

func Debug(args ...interface{})                   { log.Debug(args...) }
func Info(args ...interface{})                    { log.Info(args...) }
func Warn(args ...interface{})                    { log.Warn(args...) }
func Error(args ...interface{})                   { log.Error(args...) }
func Debugf(template string, args ...interface{}) { log.Debugf(template, args...) }
func Infof(template string, args ...interface{})  { log.Infof(template, args...) }
func Warnf(template string, args ...interface{})  { log.Warnf(template, args...) }
func Errorf(template string, args ...interface{}) { log.Errorf(template, args...) }

func TraceDebug(traceId string, args ...interface{}) { log.TraceDebug(traceId, args...) }
func TraceInfo(traceId string, args ...interface{})  { log.TraceInfo(traceId, args...) }
func TraceWarn(traceId string, args ...interface{})  { log.TraceWarn(traceId, args...) }
func TraceError(traceId string, args ...interface{}) { log.TraceError(traceId, args...) }
func TraceDebugf(traceId string, template string, args ...interface{}) {
	log.TraceDebugf(traceId, template, args...)
}
func TraceInfof(traceId string, template string, args ...interface{}) {
	log.TraceInfof(traceId, template, args...)
}
func TraceWarnf(traceId string, template string, args ...interface{}) {
	log.TraceWarnf(traceId, template, args...)
}
func TraceErrorf(traceId string, template string, args ...interface{}) {
	log.TraceErrorf(traceId, template, args...)
}

func Writer() io.Writer {
	return log.Writer()
}
