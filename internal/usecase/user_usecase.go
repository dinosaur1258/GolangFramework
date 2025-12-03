package usecase

import (
	"context"
	"database/sql"

	"github.com/dinosaur1258/GolangFramework/internal/domain/dto/request"
	"github.com/dinosaur1258/GolangFramework/internal/domain/dto/response"
	"github.com/dinosaur1258/GolangFramework/internal/domain/entity"
	"github.com/dinosaur1258/GolangFramework/internal/domain/repository"
	customerrors "github.com/dinosaur1258/GolangFramework/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (u *UserUseCase) Register(ctx context.Context, req request.RegisterRequest) (*response.UserResponse, error) {
	// 檢查 email 是否已存在
	existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, customerrors.ErrUserAlreadyExists
	}

	// 檢查 username 是否已存在
	existingUser, err = u.userRepo.GetByUsername(ctx, req.Username)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, customerrors.ErrUserAlreadyExists
	}

	// 密碼加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 建立用戶
	user := &entity.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *UserUseCase) GetUserByID(ctx context.Context, id int32) (*response.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrUserNotFound
		}
		return nil, err
	}

	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *UserUseCase) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	// 根據 email 取得用戶
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrInvalidCredentials
		}
		return nil, err
	}

	// 驗證密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, customerrors.ErrInvalidCredentials
	}

	// 生成 Token (這裡先返回空字串，等等會在 handler 生成)
	return &response.LoginResponse{
		Token: "", // 將在 handler 層生成
		User: &response.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		},
	}, nil
}

// UpdateUser 更新用戶資料
func (u *UserUseCase) UpdateUser(ctx context.Context, userID int32, req request.UpdateUserRequest) (*response.UserResponse, error) {
	// 取得當前用戶
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customerrors.ErrUserNotFound
		}
		return nil, err
	}

	// 如果要更新 email，檢查是否已被使用
	if req.Email != "" && req.Email != user.Email {
		existingUser, err := u.userRepo.GetByEmail(ctx, req.Email)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if existingUser != nil {
			return nil, customerrors.ErrUserAlreadyExists
		}
		user.Email = req.Email
	}

	// 如果要更新 username，檢查是否已被使用
	if req.Username != "" && req.Username != user.Username {
		existingUser, err := u.userRepo.GetByUsername(ctx, req.Username)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		if existingUser != nil {
			return nil, customerrors.ErrUserAlreadyExists
		}
		user.Username = req.Username
	}

	// 更新用戶
	if err := u.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// DeleteUser 刪除用戶
func (u *UserUseCase) DeleteUser(ctx context.Context, userID int32) error {
	// 檢查用戶是否存在
	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return customerrors.ErrUserNotFound
		}
		return err
	}

	// 刪除用戶
	return u.userRepo.Delete(ctx, userID)
}

// ListUsers 列出所有用戶（分頁）
func (u *UserUseCase) ListUsers(ctx context.Context, page, limit int) ([]*response.UserResponse, error) {
	// 預設值
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	// 計算 offset
	offset := (page - 1) * limit

	// 取得用戶列表
	users, err := u.userRepo.List(ctx, int32(limit), int32(offset))
	if err != nil {
		return nil, err
	}

	// 轉換成 Response
	userResponses := make([]*response.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = &response.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}
	}

	return userResponses, nil
}

// ChangePassword 修改密碼
func (u *UserUseCase) ChangePassword(ctx context.Context, userID int32, req request.ChangePasswordRequest) error {
	// 取得用戶
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return customerrors.ErrUserNotFound
		}
		return err
	}

	// 驗證舊密碼
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return customerrors.ErrInvalidCredentials
	}

	// 加密新密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)

	// 更新用戶
	return u.userRepo.Update(ctx, user)
}
