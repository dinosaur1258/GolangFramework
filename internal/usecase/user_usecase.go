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
