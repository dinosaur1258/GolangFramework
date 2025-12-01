package handler

import (
	"net/http"

	"github.com/dinosaur1258/GolangFramework/internal/domain/dto/request"
	"github.com/dinosaur1258/GolangFramework/internal/service"
	"github.com/dinosaur1258/GolangFramework/internal/usecase"
	customerrors "github.com/dinosaur1258/GolangFramework/pkg/errors"
	"github.com/dinosaur1258/GolangFramework/pkg/utils"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userUseCase *usecase.UserUseCase
	jwtService  *service.JWTService
}

func NewAuthHandler(userUseCase *usecase.UserUseCase, jwtService *service.JWTService) *AuthHandler {
	return &AuthHandler{
		userUseCase: userUseCase,
		jwtService:  jwtService,
	}
}

// Login 用戶登入
func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest

	// 解析並驗證請求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest,
			customerrors.CodeValidationFailed,
			customerrors.MsgValidationFailed,
			err.Error())
		return
	}

	// 呼叫 UseCase 驗證用戶
	loginResp, err := h.userUseCase.Login(c.Request.Context(), req)
	if err != nil {
		if err == customerrors.ErrInvalidCredentials {
			utils.ErrorResponse(c, http.StatusUnauthorized,
				customerrors.CodeInvalidCredentials,
				customerrors.MsgInvalidCredentials)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError,
			customerrors.CodeInternalServer,
			customerrors.MsgInternalServer)
		return
	}

	// 生成 JWT Token
	token, err := h.jwtService.GenerateToken(
		loginResp.User.ID,
		loginResp.User.Username,
		loginResp.User.Email,
	)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			customerrors.CodeInternalServer,
			"Failed to generate token")
		return
	}

	loginResp.Token = token

	utils.SuccessResponse(c, http.StatusOK, "Login successful", loginResp)
}
