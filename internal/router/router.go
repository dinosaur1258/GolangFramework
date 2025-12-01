package router

import (
	"github.com/dinosaur1258/GolangFramework/internal/handler"
	"github.com/dinosaur1258/GolangFramework/internal/middleware"
	"github.com/dinosaur1258/GolangFramework/internal/service"
	"github.com/gin-gonic/gin"
)

func SetupRouter(userHandler *handler.UserHandler, authHandler *handler.AuthHandler, jwtService *service.JWTService) *gin.Engine {
	r := gin.New()

	// 全域中間件
	r.Use(middleware.Recovery())     // 自定義 panic 恢復
	r.Use(gin.Logger())              // 日誌
	r.Use(middleware.ErrorHandler()) // 錯誤處理

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

		// Auth 相關路由（不需要認證）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", userHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// User 相關路由
		users := v1.Group("/users")
		{
			// 公開路由
			users.GET("/:id", userHandler.GetUser)

			// 需要認證的路由
			protected := users.Group("")
			protected.Use(middleware.AuthMiddleware(jwtService))
			{
				protected.GET("/profile", userHandler.GetProfile)
			}
		}
	}

	return r
}
