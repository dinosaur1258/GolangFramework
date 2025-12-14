package usecase

import (
	"context"
	"database/sql"

	"github.com/dinosaur1258/GolangFramework/internal/domain/contract"
	"github.com/dinosaur1258/GolangFramework/internal/domain/dto/request"
	"github.com/dinosaur1258/GolangFramework/internal/domain/dto/response"
	"github.com/dinosaur1258/GolangFramework/internal/domain/entity"
	"github.com/dinosaur1258/GolangFramework/pkg/database"
	customerrors "github.com/dinosaur1258/GolangFramework/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase struct {
	userRepo contract.UserRepository
	db       *sql.DB // ⭐ 新增:需要 DB 來執行事務
}

// ⭐ 修改:建構子需要傳入 db
func NewAuthUseCase(userRepo contract.UserRepository, db *sql.DB) *AuthUseCase {
	return &AuthUseCase{
		userRepo: userRepo,
		db:       db,
	}
}

// Register 註冊新用戶(原本的版本,保持不變)
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

// ⭐ 新增:使用事務的註冊方法
// RegisterWithTransaction 在事務中註冊用戶
// 使用場景:需要確保多個操作的原子性
func (a *AuthUseCase) RegisterWithTransaction(ctx context.Context, req request.RegisterRequest) (*response.UserResponse, error) {
	var result *response.UserResponse

	// 使用事務執行
	err := database.WithTransaction(a.db, func(txCtx context.Context) error {
		// 1. 檢查 email 是否已存在(在事務中)
		existingUser, err := a.userRepo.GetByEmail(txCtx, req.Email)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if existingUser != nil {
			return customerrors.ErrUserAlreadyExists
		}

		// 2. 檢查 username 是否已存在(在事務中)
		existingUser, err = a.userRepo.GetByUsername(txCtx, req.Username)
		if err != nil && err != sql.ErrNoRows {
			return err
		}
		if existingUser != nil {
			return customerrors.ErrUserAlreadyExists
		}

		// 3. 密碼加密
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		// 4. 建立用戶(在事務中)
		user := &entity.User{
			Username:     req.Username,
			Email:        req.Email,
			PasswordHash: string(hashedPassword),
		}

		if err := a.userRepo.Create(txCtx, user); err != nil {
			return err // 失敗會自動 rollback
		}

		// 5. 如果將來需要做其他操作(例如:寫入 audit log)
		// 都會在同一個事務中,要麼全成功,要麼全失敗

		// 6. 準備返回結果
		result = &response.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}

		return nil // 成功,會自動 commit
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Login 用戶登入(保持不變)
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

	// 生成 Token (這裡先返回空字串,等等會在 handler 生成)
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
