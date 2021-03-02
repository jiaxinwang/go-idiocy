package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// L ...
var L *zap.Logger

// S ...
var S *zap.SugaredLogger

func init() {
	L, _ = zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{"stdout"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "lv",
			EncodeLevel:  zapcore.CapitalColorLevelEncoder,
			TimeKey:      "@",
			EncodeTime:   zapcore.TimeEncoderOfLayout("01/02 15:04:05"),
			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}.Build()
	S = L.Sugar()
}
