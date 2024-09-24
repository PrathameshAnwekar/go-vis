package log

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

func Init() {
	once.Do(func() {

		config := zapcore.EncoderConfig{
			MessageKey: "message",       // Log message key
			LevelKey:   "level",         // Log level key
			TimeKey:    zapcore.OmitKey, // Remove timestamp
			EncodeLevel: zapcore.LowercaseColorLevelEncoder,
		}

		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(config),        // Console output
			zapcore.AddSync(zapcore.Lock(os.Stdout)), // Output to stdout
			zapcore.InfoLevel,                        // Log level (can be changed)
		)

		instance = zap.New(core)
	})
}

func I(s ...any) {
	str := fmt.Sprint(s...)
	instance.Info(str)
	syn()
}

func D(s ...any) {
	str := fmt.Sprint(s...)
	instance.Debug(str)
	syn()
}

func W(s ...any) {
	str := fmt.Sprint(s...)
	instance.Warn(str)
	syn()
}

// compiler inlines this function
func syn() {
	instance.Sync()
}
