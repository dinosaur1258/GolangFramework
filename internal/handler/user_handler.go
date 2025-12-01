package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dinosaur1258/GolangFramework/internal/domain/dto/request"
	"github.com/dinosaur1258/GolangFramework/internal/usecase"
	customerrors "github.com/dinosaur1258/GolangFramework/pkg/errors"
	"github.com/dinosaur1258/GolangFramework/pkg/utils"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
}

func NewUserHandler(userUseCase *usecase.UserUseCase) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
	}
}

// Register 註冊新用戶
func (h *UserHandler) Register(c *gin.Context) {
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
	user, err := h.userUseCase.Register(c.Request.Context(), req)
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

// GetUser 取得用戶資料
func (h *UserHandler) GetUser(c *gin.Context) {
	// 從 URL 參數取得 ID
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 32)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest,
			customerrors.CodeInvalidInput,
			"Invalid user ID")
		return
	}

	// 呼叫 UseCase
	user, err := h.userUseCase.GetUserByID(c.Request.Context(), int32(id))
	if err != nil {
		if err == customerrors.ErrUserNotFound || err == sql.ErrNoRows {
			utils.ErrorResponse(c, http.StatusNotFound,
				customerrors.CodeUserNotFound,
				customerrors.MsgUserNotFound)
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError,
			customerrors.CodeInternalServer,
			customerrors.MsgInternalServer)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User retrieved successfully", user)
}

// GetProfile 取得當前登入用戶的資料（需要認證）
func (h *UserHandler) GetProfile(c *gin.Context) {
	// 從 Context 取得用戶 ID（由 middleware 設定）
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized,
			customerrors.CodeUnauthorized,
			customerrors.MsgUnauthorized)
		return
	}

	// 取得用戶資料
	user, err := h.userUseCase.GetUserByID(c.Request.Context(), userID.(int32))
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			customerrors.CodeInternalServer,
			customerrors.MsgInternalServer)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Profile retrieved successfully", user)
}
