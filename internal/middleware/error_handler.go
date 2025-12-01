package middleware

import (
	"log"

	"github.com/dinosaur1258/GolangFramework/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ErrorHandler 統一錯誤處理中間件
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// 檢查是否有錯誤
		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			// 記錄錯誤
			log.Printf("Error: %v", err.Err)

			// 如果還沒有響應，返回 500 錯誤
			if !c.Writer.Written() {
				utils.ErrorResponse(c, 500, "INTERNAL_ERROR", "An unexpected error occurred")
			}
		}
	}
}

// Recovery 自定義 panic 恢復
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				utils.ErrorResponse(c, 500, "PANIC_ERROR", "Server panic occurred")
				c.Abort()
			}
		}()
		c.Next()
	}
}
