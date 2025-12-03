package router

import (
	"github.com/dinosaur1258/GolangFramework/internal/handler"
	"github.com/gin-gonic/gin"
)

// SetupAuthRoutes 設定認證相關路由
func SetupAuthRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, authHandler *handler.AuthHandler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/register", authHandler.Register) // 註冊
		auth.POST("/login", authHandler.Login)       // 登入
	}
}
