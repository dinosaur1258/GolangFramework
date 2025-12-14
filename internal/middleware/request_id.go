package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

// RequestID 為每個請求產生唯一 ID
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 檢查請求是否已經帶有 Request ID
		requestID := c.GetHeader(RequestIDHeader)

		// 如果沒有，產生新的 UUID
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 存入 Context（讓其他地方可以取用）
		c.Set("request_id", requestID)

		// 加入回應 Header（讓前端可以追蹤）
		c.Header(RequestIDHeader, requestID)

		c.Next()
	}
}

// GetRequestID 從 Context 取得 Request ID
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get("request_id"); exists {
		return requestID.(string)
	}
	return ""
}
