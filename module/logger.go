// @Project -> File    : bare-disk-perform -> logger
// @IDE    : GoLand
// @Author    : wuji
// @Date   : 2023/8/23 11:11

package module

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"runtime"
)

type MyLogger struct {
	logger *zap.Logger
}

var logger *MyLogger

func InitMyLogger() {
	zapLogger, _ := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
	}.Build()

	logger = &MyLogger{
		logger: zapLogger,
	}
}

func GetLogger() *MyLogger {
	return logger
}

func (l *MyLogger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	pc, file, line, _ := runtime.Caller(1)
	l.logger.Info(message, zap.String("caller", runtime.FuncForPC(pc).Name()), zap.String("file", file), zap.Int("line", line))
}

func (l *MyLogger) Debugf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	pc, file, line, _ := runtime.Caller(1)
	l.logger.Debug(message, zap.String("caller", runtime.FuncForPC(pc).Name()), zap.String("file", file), zap.Int("line", line))
}

func (l *MyLogger) Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	pc, file, line, _ := runtime.Caller(1)
	l.logger.Warn(message, zap.String("caller", runtime.FuncForPC(pc).Name()), zap.String("file", file), zap.Int("line", line))
}

func (l *MyLogger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	pc, file, line, _ := runtime.Caller(1)
	l.logger.Error(message, zap.String("caller", runtime.FuncForPC(pc).Name()), zap.String("file", file), zap.Int("line", line))
}

func (l *MyLogger) Fatalf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	pc, file, line, _ := runtime.Caller(1)
	l.logger.Fatal(message, zap.String("caller", runtime.FuncForPC(pc).Name()), zap.String("file", file), zap.Int("line", line))
}
