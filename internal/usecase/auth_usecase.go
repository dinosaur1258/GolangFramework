package usecase

import (
	"context"
	"database/sql"

	"github.com/dinosaur1258/GolangFramework/internal/domain/contract"
	"github.com/dinosaur1258/GolangFramework/internal/domain/dto/request"
	"github.com/dinosaur1258/GolangFramework/internal/domain/dto/response"
	"github.com/dinosaur1258/GolangFramework/internal/domain/entity"
	customerrors "github.com/dinosaur1258/GolangFramework/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo contract.UserRepository
}

func NewAuthUseCase(userRepo contract.UserRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
	}
}

// Register 註冊新用戶
func (a *AuthUseCase) Register(ctx context.Context, req request.RegisterRequest) (*response.UserResponse, error) {
	// 檢查 email 是否已存在
	existingUser, err := a.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, customerrors.ErrUserAlreadyExists
	}

	// 檢查 username 是否已存在
	existingUser, err = a.userRepo.GetByUsername(ctx, req.Username)
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

	if err := a.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return &response.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}, nil
}

// Login 用戶登入
func (a *AuthUseCase) Login(ctx context.Context, req request.LoginRequest) (*response.LoginResponse, error) {
	// 根據 email 取得用戶
	user, err := a.userRepo.GetByEmail(ctx, req.Email)
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
