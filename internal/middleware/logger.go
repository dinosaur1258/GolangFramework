package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// RequestLogger 請求日誌中間件
func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 處理請求
		c.Next()

		// 計算請求耗時
		cost := time.Since(start)

		// 記錄請求信息
		fields := []zap.Field{
			zap.Int("status", c.Writer.Status()),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.Duration("latency", cost),
			zap.String("request_id", c.GetString("request_id")),
		}

		// 如果有錯誤，添加錯誤信息
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		// 根據狀態碼選擇日誌級別
		switch {
		case c.Writer.Status() >= 500:
			logger.Error("Server error", fields...)
		case c.Writer.Status() >= 400:
			logger.Warn("Client error", fields...)
		case c.Writer.Status() >= 300:
			logger.Info("Redirection", fields...)
		default:
			logger.Info("Success", fields...)
		}
	}
}
