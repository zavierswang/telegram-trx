package logger

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestCustomField(t *testing.T) {
	var ops = []TreeOption{
		{
			FileName: "access.log",
			Rpt: RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   true,
			},
			Lef: func(level zapcore.Level) bool {
				return level <= zap.InfoLevel
			},
		},
		{
			FileName: "error.log",
			Rpt: RotateOptions{
				MaxSize:    1,
				MaxAge:     1,
				MaxBackups: 3,
				Compress:   true,
			},
			Lef: func(level zapcore.Level) bool {
				return level > zap.InfoLevel
			},
		},
	}

	logger := NewRotate(ops)
	ResetDefault(logger)
	for i := 0; i < 2000000; i++ {
		field := &CustomField{
			UID:      fmt.Sprintf("%d", i),
			UserName: fmt.Sprintf("username_%d", i),
		}
		Warn("testing warn", zap.Inline(field))
	}

	assert.FileExists(t, "access.log")
	assert.FileExists(t, "error.log")
}

type CustomField struct {
	UID      string
	UserName string
}

func (f CustomField) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("uid", f.UID)
	enc.AddString("username", f.UserName)
	return nil
}
