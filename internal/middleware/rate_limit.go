package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimit 限制請求頻率
func RateLimit() gin.HandlerFunc {
	// 設定限制規則：每分鐘 100 次請求
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	}

	// 使用記憶體儲存（生產環境建議用 Redis）
	store := memory.NewStore()

	// 建立限制器
	instance := limiter.New(store, rate)

	// 建立 Gin 中間件
	middleware := mgin.NewMiddleware(instance)

	return middleware
}

// RateLimitStrict 嚴格限制（例如登入 API）
func RateLimitStrict() gin.HandlerFunc {
	// 每分鐘只能 10 次請求
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  10,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)
	middleware := mgin.NewMiddleware(instance)

	return middleware
}

// RateLimitWithCustomError 自定義錯誤訊息的限流
func RateLimitWithCustomError(requestsPerMinute int64) gin.HandlerFunc {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  requestsPerMinute,
	}

	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// 取得客戶端 IP
		key := c.ClientIP()

		// 檢查是否超過限制
		context, err := instance.Get(c, key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limiter error",
			})
			c.Abort()
			return
		}

		// 設定 Header
		c.Header("X-RateLimit-Limit", strconv.FormatInt(context.Limit, 10))
		c.Header("X-RateLimit-Remaining", strconv.FormatInt(context.Remaining, 10))
		c.Header("X-RateLimit-Reset", strconv.FormatInt(context.Reset, 10))

		// 如果超過限制
		if context.Reached {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"success": false,
				"error": gin.H{
					"code":    "RATE_LIMIT_EXCEEDED",
					"message": "Too many requests, please try again later",
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
