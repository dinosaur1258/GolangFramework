package postgres

import (
	"context"
	"database/sql"

	"github.com/dinosaur1258/GolangFramework/db/sqlc"
	"github.com/dinosaur1258/GolangFramework/internal/domain/contract"
	"github.com/dinosaur1258/GolangFramework/internal/domain/entity"
)

type userRepository struct {
	queries *sqlc.Queries
}

// 確保實作了 Interface
var _ contract.UserRepository = (*userRepository)(nil)

func NewUserRepository(db *sql.DB) contract.UserRepository {
	return &userRepository{
		queries: sqlc.New(db),
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	params := sqlc.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	createdUser, err := r.queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	// 更新 user 的 ID 和時間
	user.ID = createdUser.ID
	user.CreatedAt = createdUser.CreatedAt
	user.UpdatedAt = createdUser.UpdatedAt

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (*entity.User, error) {
	sqlcUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           sqlcUser.ID,
		Username:     sqlcUser.Username,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		CreatedAt:    sqlcUser.CreatedAt,
		UpdatedAt:    sqlcUser.UpdatedAt,
	}, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	sqlcUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           sqlcUser.ID,
		Username:     sqlcUser.Username,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		CreatedAt:    sqlcUser.CreatedAt,
		UpdatedAt:    sqlcUser.UpdatedAt,
	}, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	sqlcUser, err := r.queries.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:           sqlcUser.ID,
		Username:     sqlcUser.Username,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		CreatedAt:    sqlcUser.CreatedAt,
		UpdatedAt:    sqlcUser.UpdatedAt,
	}, nil
}

func (r *userRepository) List(ctx context.Context, limit, offset int32) ([]*entity.User, error) {
	params := sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	}

	sqlcUsers, err := r.queries.ListUsers(ctx, params)
	if err != nil {
		return nil, err
	}

	users := make([]*entity.User, len(sqlcUsers))
	for i, sqlcUser := range sqlcUsers {
		users[i] = &entity.User{
			ID:           sqlcUser.ID,
			Username:     sqlcUser.Username,
			Email:        sqlcUser.Email,
			PasswordHash: sqlcUser.PasswordHash,
			CreatedAt:    sqlcUser.CreatedAt,
			UpdatedAt:    sqlcUser.UpdatedAt,
		}
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) error {
	params := sqlc.UpdateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	updatedUser, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return err
	}

	user.UpdatedAt = updatedUser.UpdatedAt
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	return r.queries.DeleteUser(ctx, id)
}
