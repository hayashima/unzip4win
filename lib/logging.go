package unzip4win

import (
	"go.uber.org/zap"
)

var LibLogger *zap.Logger

func debugLog(msg string, fields ...zap.Field) {
	if LibLogger == nil {
		return
	}
	LibLogger.Debug(msg, fields...)
}

func AppLogger(isDebug bool) (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stderr"}
	config.Level.SetLevel(zap.ErrorLevel)
	if isDebug {
		config.Level.SetLevel(zap.DebugLevel)
	}
	return config.Build()
}
