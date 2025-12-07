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
	authUseCase *usecase.AuthUseCase
	jwtService  *service.JWTService
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase, jwtService *service.JWTService) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		jwtService:  jwtService,
	}
}

// Register godoc
// @Summary      註冊新用戶
// @Description  註冊一個新的用戶帳號
// @Tags         認證
// @Accept       json
// @Produce      json
// @Param        request body request.RegisterRequest true "註冊資料"
// @Success      201  {object}  utils.Response{data=response.UserResponse}
// @Failure      400  {object}  utils.Response
// @Failure      409  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest

	// 解析並驗證請求
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest,
			customerrors.CodeValidationFailed,
			customerrors.MsgValidationFailed,
			err.Error())
		return
	}

	// 呼叫 UseCase
	user, err := h.authUseCase.Register(c.Request.Context(), req)
	if err != nil {
		if err == customerrors.ErrUserAlreadyExists {
			utils.ErrorResponse(c, http.StatusConflict,
				customerrors.CodeUserAlreadyExists,
				customerrors.MsgUserAlreadyExists)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError,
			customerrors.CodeInternalServer,
			customerrors.MsgInternalServer)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "User registered successfully", user)
}

// Login godoc
// @Summary      用戶登入
// @Description  使用 Email 和密碼登入
// @Tags         認證
// @Accept       json
// @Produce      json
// @Param        request body request.LoginRequest true "登入資料"
// @Success      200  {object}  utils.Response{data=response.LoginResponse}
// @Failure      400  {object}  utils.Response
// @Failure      401  {object}  utils.Response
// @Failure      500  {object}  utils.Response
// @Router       /auth/login [post]
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
	loginResp, err := h.authUseCase.Login(c.Request.Context(), req)
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
