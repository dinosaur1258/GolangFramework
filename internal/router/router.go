package router

import (
	"github.com/dinosaur1258/GolangFramework/internal/handler"
	"github.com/dinosaur1258/GolangFramework/internal/middleware"
	"github.com/dinosaur1258/GolangFramework/internal/service"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 設定主路由
func SetupRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, jwtService *service.JWTService) *gin.Engine {
	r := gin.New()

	// 全域中間件
	r.Use(middleware.Recovery())     // 自定義 panic 恢復
	r.Use(gin.Logger())              // 日誌
	r.Use(middleware.ErrorHandler()) // 錯誤處理

	// Swagger 文檔路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 群組
	v1 := r.Group("/api/v1")
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
