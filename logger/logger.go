package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// L ...
var L *zap.Logger

// S ...
var S *zap.SugaredLogger

var (
	ColorBlack        = ColorString("\033[1;30m%s\033[0m")
	ColorRed          = ColorString("\033[1;31m%s\033[0m")
	ColorGreen        = ColorString("\033[1;32m%s\033[0m")
	ColorYellow       = ColorString("\033[1;33m%s\033[0m")
	ColorBlue         = ColorString("\033[1;34m%s\033[0m")
	ColorMagenta      = ColorString("\033[1;35m%s\033[0m")
	ColorCyan         = ColorString("\033[1;36m%s\033[0m")
	ColorLightGray    = ColorString("\033[1;37m%s\033[0m")
	ColorDefaultColor = ColorString("\033[1;39m%s\033[0m")
	ColorDarkGray     = ColorString("\033[1;90m%s\033[0m")
	ColorLightRed     = ColorString("\033[1;91m%s\033[0m")
	ColorLightGreen   = ColorString("\033[1;92m%s\033[0m")
	ColorLightYellow  = ColorString("\033[1;93m%s\033[0m")
	ColorLightBlue    = ColorString("\033[1;94m%s\033[0m")
	ColorLightMagenta = ColorString("\033[1;95m%s\033[0m")
	ColorLightCyan    = ColorString("\033[1;96m%s\033[0m")
	ColorWhite        = ColorString("\033[1;97m%s\033[0m")
)

func ColorString(colorString string) func(...interface{}) string {
	return func(args ...interface{}) string {
		return fmt.Sprintf(colorString, fmt.Sprint(args...))
	}
}

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
