package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Timeout 設定請求超時
func Timeout(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 建立帶有超時的 context
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		// 替換 request 的 context
		c.Request = c.Request.WithContext(ctx)

		// 用 channel 來處理請求完成
		finished := make(chan struct{})

		go func() {
			c.Next()
			finished <- struct{}{}
		}()

		select {
		case <-finished:
			// 請求正常完成
			return
		case <-ctx.Done():
			// 請求超時
			c.JSON(http.StatusRequestTimeout, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "REQUEST_TIMEOUT",
					"message": "Request timeout",
				},
			})
			c.Abort()
		}
	}
}
