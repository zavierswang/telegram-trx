package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
	"runtime"
)

type CallerEncoder struct {
	zapcore.Encoder
}

func NewEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	return CallerEncoder{
		Encoder: zapcore.NewJSONEncoder(cfg),
	}
}

func (enc CallerEncoder) Clone() zapcore.Encoder {
	return CallerEncoder{
		enc.Encoder.Clone(),
	}
}

func (enc CallerEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	pc, _, line, _ := runtime.Caller(5)
	filename := runtime.FuncForPC(pc).Name()
	fields = append(fields, zap.String("func", filename))
	fields = append(fields, zap.Int("lineno", line))
	//fields = append(fields, zap.String("file", file))
	//fields = append(fields, zap.String("tag", "TelegramBot"))
	return enc.Encoder.EncodeEntry(entry, fields)
}

type entryCaller struct {
	*zapcore.EntryCaller
}

func (c entryCaller) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	pc, _, _, _ := runtime.Caller(12)
	runtime.FuncForPC(pc).Name()
	//fmt.Println(filename)
	//enc.AddString("filename", filename)
	return nil
}
