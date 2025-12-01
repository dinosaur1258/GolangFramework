package middleware

import (
	"net/http"
	"strings"

	"github.com/dinosaur1258/GolangFramework/internal/service"
	customerrors "github.com/dinosaur1258/GolangFramework/pkg/errors"
	"github.com/dinosaur1258/GolangFramework/pkg/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtService *service.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 從 Header 取得 Token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized,
				customerrors.CodeUnauthorized,
				"Authorization header required")
			c.Abort()
			return
		}

		// 檢查格式：Bearer <token>
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized,
				customerrors.CodeUnauthorized,
				"Invalid authorization header format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 驗證 Token
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized,
				customerrors.CodeUnauthorized,
				"Invalid or expired token")
			c.Abort()
			return
		}

		// 將用戶資訊存入 Context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)

		c.Next()
	}
}
