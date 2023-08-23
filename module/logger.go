// @Project -> File    : bare-disk-perform -> logger
// @IDE    : GoLand
// @Author    : wuji
// @Date   : 2023/8/23 11:11

package module

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type MyLogger struct {
	logger *zap.Logger
}

func NewMyLogger() *MyLogger {
	logger, _ := zap.Config{
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
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}.Build()

	return &MyLogger{
		logger: logger,
	}
}

func (l *MyLogger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logger.Info(message)
}

func (l *MyLogger) Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logger.Warn(message)
}

func (l *MyLogger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logger.Error(message)
}

func (l *MyLogger) Fatalf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	l.logger.Fatal(message)
}
