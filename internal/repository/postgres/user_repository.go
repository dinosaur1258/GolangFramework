package postgres

import (
	"context"
	"database/sql"

	"github.com/dinosaur1258/GolangFramework/db/sqlc"
	"github.com/dinosaur1258/GolangFramework/internal/domain/contract"
	"github.com/dinosaur1258/GolangFramework/internal/domain/entity"
	"github.com/dinosaur1258/GolangFramework/pkg/database"
)

type userRepository struct {
	db *sql.DB // 保留原本的 DB 連線
}

var _ contract.UserRepository = (*userRepository)(nil)

func NewUserRepository(db *sql.DB) contract.UserRepository {
	return &userRepository{
		db: db,
	}
}

// ⭐ 新增:智能選擇使用 DB 或 TX
func (r *userRepository) getQueries(ctx context.Context) *sqlc.Queries {
	// 如果 context 中有 transaction,就用 transaction
	if tx, ok := database.GetTx(ctx); ok {
		return sqlc.New(tx)
	}
	// 否則使用正常的 DB 連線
	return sqlc.New(r.db)
}

// ⭐ 修改:使用 getQueries 取代原本的 r.queries
func (r *userRepository) Create(ctx context.Context, user *entity.User) error {
	queries := r.getQueries(ctx) // 智能選擇

	params := sqlc.CreateUserParams{
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	createdUser, err := queries.CreateUser(ctx, params)
	if err != nil {
		return err
	}

	user.ID = createdUser.ID
	user.CreatedAt = createdUser.CreatedAt
	user.UpdatedAt = createdUser.UpdatedAt

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, id int32) (*entity.User, error) {
	queries := r.getQueries(ctx) // 智能選擇

	sqlcUser, err := queries.GetUserByID(ctx, id)
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
	queries := r.getQueries(ctx) // 智能選擇

	sqlcUser, err := queries.GetUserByEmail(ctx, email)
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
	queries := r.getQueries(ctx) // 智能選擇

	sqlcUser, err := queries.GetUserByUsername(ctx, username)
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
	queries := r.getQueries(ctx) // 智能選擇

	params := sqlc.ListUsersParams{
		Limit:  limit,
		Offset: offset,
	}

	sqlcUsers, err := queries.ListUsers(ctx, params)
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
	queries := r.getQueries(ctx) // 智能選擇

	params := sqlc.UpdateUserParams{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	}

	updatedUser, err := queries.UpdateUser(ctx, params)
	if err != nil {
		return err
	}

	user.UpdatedAt = updatedUser.UpdatedAt
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id int32) error {
	queries := r.getQueries(ctx) // 智能選擇
	return queries.DeleteUser(ctx, id)
}
