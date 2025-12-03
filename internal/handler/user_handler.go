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

// UpdateProfile 更新當前用戶資料（需要認證）
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 從 Context 取得用戶 ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized,
			customerrors.CodeUnauthorized,
			customerrors.MsgUnauthorized)
		return
	}

	var req request.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest,
			customerrors.CodeValidationFailed,
			customerrors.MsgValidationFailed,
			err.Error())
		return
	}

	// 呼叫 UseCase
	user, err := h.userUseCase.UpdateUser(c.Request.Context(), userID.(int32), req)
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

	utils.SuccessResponse(c, http.StatusOK, "User updated successfully", user)
}

// DeleteUser 刪除用戶（需要認證，只能刪除自己）
func (h *UserHandler) DeleteUser(c *gin.Context) {
	// 從 Context 取得用戶 ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized,
			customerrors.CodeUnauthorized,
			customerrors.MsgUnauthorized)
		return
	}

	// 呼叫 UseCase
	if err := h.userUseCase.DeleteUser(c.Request.Context(), userID.(int32)); err != nil {
		if err == customerrors.ErrUserNotFound {
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

	utils.SuccessResponse(c, http.StatusOK, "User deleted successfully", nil)
}

// ListUsers 列出所有用戶（需要認證）
func (h *UserHandler) ListUsers(c *gin.Context) {
	var req request.ListUsersRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest,
			customerrors.CodeValidationFailed,
			customerrors.MsgValidationFailed,
			err.Error())
		return
	}

	// 呼叫 UseCase
	users, err := h.userUseCase.ListUsers(c.Request.Context(), req.Page, req.Limit)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError,
			customerrors.CodeInternalServer,
			customerrors.MsgInternalServer)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Users retrieved successfully", gin.H{
		"users": users,
		"page":  req.Page,
		"limit": req.Limit,
	})
}

// ChangePassword 修改密碼（需要認證）
func (h *UserHandler) ChangePassword(c *gin.Context) {
	// 從 Context 取得用戶 ID
	userID, exists := c.Get("user_id")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized,
			customerrors.CodeUnauthorized,
			customerrors.MsgUnauthorized)
		return
	}

	var req request.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest,
			customerrors.CodeValidationFailed,
			customerrors.MsgValidationFailed,
			err.Error())
		return
	}

	// 呼叫 UseCase
	if err := h.userUseCase.ChangePassword(c.Request.Context(), userID.(int32), req); err != nil {
		if err == customerrors.ErrInvalidCredentials {
			utils.ErrorResponse(c, http.StatusUnauthorized,
				customerrors.CodeInvalidCredentials,
				"Old password is incorrect")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError,
			customerrors.CodeInternalServer,
			customerrors.MsgInternalServer)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Password changed successfully", nil)
}
