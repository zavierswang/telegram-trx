package logger

import (
	"fmt"
	"go.uber.org/zap"
)

type Field = zap.Field

type Logger struct {
	zap    *zap.Logger
	level  Level
	fields []zap.Field
}

func (l *Logger) Debug(format string, args ...interface{}) {
	l.zap.Debug(fmt.Sprintf(format, args...), l.fields...)
}

func (l *Logger) Info(format string, args ...interface{}) {
	l.zap.Info(fmt.Sprintf(format, args...), l.fields...)
}

func (l *Logger) Warn(format string, args ...interface{}) {
	l.zap.Warn(fmt.Sprintf(format, args...), l.fields...)
}

func (l *Logger) Error(format string, args ...interface{}) {
	l.zap.Error(fmt.Sprintf(format, args...), l.fields...)
}

func (l *Logger) DPanic(format string, args ...interface{}) {
	l.zap.DPanic(fmt.Sprintf(format, args...), l.fields...)
}

func (l *Logger) Panic(format string, args ...interface{}) {
	l.zap.Panic(fmt.Sprintf(format, args...), l.fields...)
}

func (l *Logger) Fatal(format string, args ...interface{}) {
	l.zap.Fatal(fmt.Sprintf(format, args...), l.fields...)
}

func (l *Logger) With(fields ...Field) *Logger {
	logger := l.zap.With(fields...)
	l.zap = logger
	return l
}

func (l *Logger) Sync() error {
	return l.zap.Sync()
}

func GetLogger() *zap.Logger {
	return zap.L()
}
