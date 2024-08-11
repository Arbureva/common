package czap

import (
	"fmt"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

type CustomZapEncoder struct {
	zapcore.Encoder
	appName string
}

func (ce *CustomZapEncoder) Clone() zapcore.Encoder {
	return &CustomZapEncoder{
		Encoder: ce.Encoder.Clone(),
		appName: ce.appName,
	}
}

func (ce *CustomZapEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	entry.Message = fmt.Sprintf("[%s] %s", ce.appName, entry.Message)
	return ce.Encoder.EncodeEntry(entry, fields)
}

func NewCustomZapEncoder(encoderConfig zapcore.EncoderConfig, appName string) zapcore.Encoder {
	return &CustomZapEncoder{
		Encoder: zapcore.NewConsoleEncoder(encoderConfig),
		appName: appName,
	}
}
