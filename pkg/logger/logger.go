package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

// InitLogger 初始化日誌系統
func InitLogger(env string) error {
	var config zap.Config

	if env == "production" {
		// 生產環境：JSON 格式，Info 級別
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		// 開發環境：彩色輸出，Debug 級別
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	// 設定日誌輸出
	config.OutputPaths = []string{"stdout", "./logs/app.log"}
	config.ErrorOutputPaths = []string{"stderr", "./logs/error.log"}

	// 確保日誌目錄存在（移到前面）
	if err := os.MkdirAll("./logs", 0755); err != nil {
		return err
	}

	// 建立 logger
	var err error
	Log, err = config.Build(
		zap.AddCallerSkip(1),
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
	)
	if err != nil {
		return err
	}

	return nil
}

// Sync 刷新日誌緩衝區
func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}

// 便利方法
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

// WithRequestID 創建帶有 request ID 的 logger
func WithRequestID(requestID string) *zap.Logger {
	return Log.With(zap.String("request_id", requestID))
}
