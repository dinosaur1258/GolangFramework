package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// CORS 跨域資源共享設定
func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		// 允許的來源（開發階段允許所有，生產環境應該指定特定域名）
		AllowOrigins: []string{"*"}, // 或 []string{"http://localhost:3000", "https://yourdomain.com"}

		// 允許的 HTTP 方法
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},

		// 允許的 Header
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
		},

		// 暴露的 Header（讓前端能讀取）
		ExposeHeaders: []string{
			"Content-Length",
			"X-Request-ID",
		},

		// 允許攜帶認證資訊（Cookies、Authorization header）
		AllowCredentials: true,

		// 預檢請求的快取時間
		MaxAge: 12 * time.Hour,
	})
}
