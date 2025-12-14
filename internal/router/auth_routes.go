package router

import (
	"github.com/dinosaur1258/GolangFramework/internal/handler"
	"github.com/dinosaur1258/GolangFramework/internal/middleware"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 設定認證相關路由
func SetupAuthRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		// 註冊和登入使用嚴格限流（每分鐘 10 次）
		auth.POST("/register", middleware.RateLimitStrict(), authHandler.Register)
		auth.POST("/login", middleware.RateLimitStrict(), authHandler.Login)
	}
}
