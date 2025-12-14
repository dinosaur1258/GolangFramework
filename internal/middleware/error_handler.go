package middleware

import (
	"github.com/dinosaur1258/GolangFramework/pkg/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// ErrorHandler 統一錯誤處理中間件（整合日誌）
func ErrorHandler(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 檢查是否有錯誤
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 使用結構化日誌記錄錯誤
			logger.Error("Request error",
				zap.Error(err.Err),
				zap.String("request_id", c.GetString("request_id")),
				zap.String("method", c.Request.Method),
				zap.String("path", c.Request.URL.Path),
				zap.String("ip", c.ClientIP()),
			)

			// 如果還沒有響應，返回 500 錯誤
			if !c.Writer.Written() {
				utils.ErrorResponse(c, 500, "INTERNAL_ERROR", "An unexpected error occurred")
			}
		}
	}
}

// Recovery 自定義 panic 恢復（整合日誌）
func Recovery(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 使用結構化日誌記錄 panic，包含堆疊追蹤
				logger.Error("Panic recovered",
					zap.Any("error", err),
					zap.String("request_id", c.GetString("request_id")),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("ip", c.ClientIP()),
					zap.String("user_agent", c.Request.UserAgent()),
					zap.Stack("stacktrace"),
				)

				utils.ErrorResponse(c, 500, "PANIC_ERROR", "Server panic occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
