package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Level = zapcore.Level

const (
	InfoLevel   = zap.InfoLevel
	WarnLevel   = zap.WarnLevel
	ErrorLevel  = zap.ErrorLevel
	DPanicLevel = zap.DPanicLevel
	PanicLevel  = zap.PanicLevel
	FatalLevel  = zap.FatalLevel
	DebugLevel  = zap.DebugLevel
)
