package mock

import (
	"context"

	"github.com/dinosaur1258/GolangFramework/internal/domain/entity"
)

type SimpleMockUserRepository struct {
	User  *entity.User
	Error error
}

func (m *SimpleMockUserRepository) Create(ctx context.Context, user *entity.User) error {
	return m.Error
}

func (m *SimpleMockUserRepository) GetByID(ctx context.Context, id int32) (*entity.User, error) {
	return m.User, m.Error
}

func (m *SimpleMockUserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	return m.User, m.Error
}

func (m *SimpleMockUserRepository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	return m.User, m.Error
}

func (m *SimpleMockUserRepository) List(ctx context.Context, limit, offset int32) ([]*entity.User, error) {
	if m.User != nil {
		return []*entity.User{m.User}, m.Error
	}
	return nil, m.Error
}

func (m *SimpleMockUserRepository) Update(ctx context.Context, user *entity.User) error {
	return m.Error
}

func (m *SimpleMockUserRepository) Delete(ctx context.Context, id int32) error {
	return m.Error
}
