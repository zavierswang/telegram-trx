package logger

import (
	"os"
	"testing"
)

func TestDefault(t *testing.T) {
	file, err := os.OpenFile(".log/access.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		panic(err)
	}

	logger := New(file, InfoLevel)
	ResetDefault(logger)
	defer Sync()

	Info("默认测试")
}
