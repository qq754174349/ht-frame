package log

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"ht-crm/src/ht/config"
	"io"
	"os"
	"time"
)

type zapLog struct {
	logger *zap.SugaredLogger
}

func newZapLog(appLogCfg *config.LogConfig) Logger {
	level, err := zap.ParseAtomicLevel(appLogCfg.Level)
	if err != nil {
		panic(err)
	}
	encoder := zapcore.NewConsoleEncoder(*newEncoderConfig())

	// 控制台输出日志
	stdWriter := zapcore.AddSync(os.Stdout)
	stdoutEncoder := zapcore.NewCore(encoder, stdWriter, level)

	// 文件输出日志
	fileWriter := zapcore.AddSync(newFileWriteSyncer(appLogCfg.FileUrl))
	fileEncoder := zapcore.NewCore(encoder, fileWriter, level)

	logger := zap.New(zapcore.NewTee(stdoutEncoder, fileEncoder), zap.WithCaller(true))
	// 全局使用
	zap.ReplaceGlobals(logger)
	sugar := logger.Sugar()
	sugar.Info("init logger complete")
	sugar.Info("use zap logger")
	sugar.Infof("current log level %s", appLogCfg.Level)

	zapLog := zapLog{logger: sugar.WithOptions(zap.AddCallerSkip(3))}

	return zapLog
}

func newEncoderConfig() *zapcore.EncoderConfig {
	return &zapcore.EncoderConfig{
		MessageKey:   "msg",
		LevelKey:     "level",
		TimeKey:      "time",
		CallerKey:    "caller",
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeTime:   zapcore.TimeEncoderOfLayout(time.DateTime),
		EncodeCaller: zapcore.FullCallerEncoder,
	}
}

func newFileWriteSyncer(outPutPath string) *lumberjack.Logger {
	logFile := fmt.Sprintf("%s%s-%s.log", outPutPath, config.GetAppName(), time.Now().Format(time.DateOnly))
	return &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10,   // 每个日志文件的最大大小（MB）
		MaxBackups: 3,    // 最多保留的备份日志文件数
		MaxAge:     28,   // 保留旧日志的最大天数
		Compress:   true, // 是否压缩旧日志
	}
}

func (log zapLog) Debug(args ...interface{})                   { log.logger.Debug(args...) }
func (log zapLog) Info(args ...interface{})                    { log.logger.Info(args...) }
func (log zapLog) Warn(args ...interface{})                    { log.logger.Warn(args...) }
func (log zapLog) Error(args ...interface{})                   { log.logger.Error(args...) }
func (log zapLog) Debugf(template string, args ...interface{}) { log.logger.Debugf(template, args...) }
func (log zapLog) Infof(template string, args ...interface{})  { log.logger.Infof(template, args...) }
func (log zapLog) Warnf(template string, args ...interface{})  { log.logger.Warnf(template, args...) }
func (log zapLog) Errorf(template string, args ...interface{}) { log.logger.Errorf(template, args...) }

func (log zapLog) TraceDebug(traceId string, args ...interface{}) {
	args = append([]interface{}{"\t[" + traceId + "]\t"}, args...)
	log.Debug(args...)
}
func (log zapLog) TraceInfo(traceId string, args ...interface{}) {
	args = append([]interface{}{"\t[" + traceId + "]\t"}, args...)
	log.Info(args...)
}
func (log zapLog) TraceWarn(traceId string, args ...interface{}) {
	args = append([]interface{}{"\t[" + traceId + "]\t"}, args...)
	log.Warn(args...)
}
func (log zapLog) TraceError(traceId string, args ...interface{}) {
	args = append([]interface{}{"\t[" + traceId + "]\t"}, args...)
	log.Error(args...)
}
func (log zapLog) TraceDebugf(traceId string, template string, args ...interface{}) {
	log.Debugf("\t["+traceId+"]\t"+template, args...)
}
func (log zapLog) TraceInfof(traceId string, template string, args ...interface{}) {
	log.Infof("\t["+traceId+"]\t"+template, args...)
}
func (log zapLog) TraceWarnf(traceId string, template string, args ...interface{}) {
	log.Warnf("\t["+traceId+"]\t"+template, args...)
}
func (log zapLog) TraceErrorf(traceId string, template string, args ...interface{}) {
	log.Errorf("\t["+traceId+"]\t"+template, args...)
}

func (log zapLog) RedirectStdLog() {
	zap.RedirectStdLog(log.logger.Desugar().WithOptions(zap.AddCallerSkip(-3)))
}

func (log zapLog) Writer() io.Writer {
	ginLog := log.logger.Desugar().WithOptions(zap.AddCallerSkip(-3))
	return zap.NewStdLog(ginLog).Writer()
}
