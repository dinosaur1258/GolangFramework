package router

import (
	"time"

	"github.com/dinosaur1258/GolangFramework/internal/handler"
	"github.com/dinosaur1258/GolangFramework/internal/middleware"
	"github.com/dinosaur1258/GolangFramework/internal/service"
	"github.com/dinosaur1258/GolangFramework/pkg/logger"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 設定主路由
func SetupRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, jwtService *service.JWTService) *gin.Engine {
	r := gin.New()

	// 全域中間件（按順序執行）
	r.Use(middleware.Recovery(logger.Log))      // 1. Panic 恢復（整合日誌）
	r.Use(middleware.RequestID())               // 2. Request ID
	r.Use(middleware.RequestLogger(logger.Log)) // 3. 請求日誌（取代 gin.Logger()）
	r.Use(middleware.CORS())                    // 4. CORS
	r.Use(middleware.Timeout(30 * time.Second)) // 5. 超時控制
	r.Use(middleware.ErrorHandler(logger.Log))  // 6. 錯誤處理（整合日誌）

	// Swagger 文檔路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 群組
	v1 := r.Group("/api/v1")
	v1.Use(middleware.RateLimit()) // API 群組使用一般限流
	{
		// 健康檢查
		v1.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"message": "Server is running",
			})
		})

		// 註冊各模組路由
		SetupAuthRoutes(v1, userHandler, authHandler)
		SetupUserRoutes(v1, userHandler, jwtService)
	}

	return r
}
