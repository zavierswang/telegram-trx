package logger

import (
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"telegram-trx/pkg/core/cst"
	"telegram-trx/pkg/utils"
)

const LogDir = "./logs"

func New(writer io.Writer, level Level, extMap ...map[string]string) *Logger {
	if ok, _ := utils.PathExists(LogDir); !ok {
		_ = os.Mkdir(LogDir, os.ModePerm)
	}
	if writer == nil {
		panic("the writer is nil")
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:     "msg",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "name",
		FunctionKey:    "func",
		StacktraceKey:  "strace",
		SkipLineEnding: false,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
	core := zapcore.NewCore(
		NewEncoder(cfg.EncoderConfig),
		getLogWriter(),
		level,
	)
	logger := &Logger{
		zap:   zap.New(core),
		level: level,
	}

	// 初始化默认字段
	fs := make([]zap.Field, 0)
	for _, ext := range extMap {
		for key, value := range ext {
			fs = append(fs, zap.String(key, value))
		}
	}
	logger = logger.With(fs...)
	return logger
}

// 使用 lumberjack 作为日志写入器
func getLogWriter() zapcore.WriteSyncer {
	pf, err := os.Open(fmt.Sprintf("%s.yaml", cst.AppName))
	if err != nil {
		panic(err)
	}
	viper.SetConfigType("yaml")
	err = viper.ReadConfig(pf)
	if err != nil {
		panic(err)
	}
	env := viper.Get("app.env").(string)
	var file io.Writer
	if env == "release" || env == "prod" || env == "ttkj" {
		file = &lumberjack.Logger{
			Filename:   fmt.Sprintf("./logs/%s.log", cst.AppName),
			MaxSize:    200,
			MaxBackups: 7,
			MaxAge:     5,
			Compress:   true,
		}
	} else {
		file = os.Stdout
	}
	return zapcore.AddSync(file)
}
