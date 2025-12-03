package router

import (
	"github.com/dinosaur1258/GolangFramework/internal/handler"
	"github.com/dinosaur1258/GolangFramework/internal/middleware"
	"github.com/dinosaur1258/GolangFramework/internal/service"
	"github.com/gin-gonic/gin"
)

// SetupUserRoutes 設定用戶相關路由
func SetupUserRoutes(rg *gin.RouterGroup, userHandler *handler.UserHandler, jwtService *service.JWTService) {
	users := rg.Group("/users")
	{
		// 公開路由：查看用戶資料
		users.GET("/:id", userHandler.GetUser)

		// 需要認證的路由
		protected := users.Group("")
		protected.Use(middleware.AuthMiddleware(jwtService))
		{
			// 個人資料管理
			protected.GET("/profile", userHandler.GetProfile)    // 取得個人資料
			protected.PUT("/profile", userHandler.UpdateProfile) // 更新個人資料
			protected.DELETE("/profile", userHandler.DeleteUser) // 刪除帳號

			// 密碼管理
			protected.PUT("/password", userHandler.ChangePassword) // 修改密碼

			// 用戶列表（管理用）
			protected.GET("", userHandler.ListUsers) // 列出所有用戶
		}
	}
}
