package usecase

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dinosaur1258/GolangFramework/internal/domain/entity"
	"github.com/dinosaur1258/GolangFramework/internal/repository/mock"
	customerrors "github.com/dinosaur1258/GolangFramework/pkg/errors"
)

func TestGetUserByID_Success(t *testing.T) {
	fakeUser := &entity.User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hash123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mockRepo := &mock.SimpleMockUserRepository{
		User:  fakeUser,
		Error: nil,
	}

	usecase := NewUserUseCase(mockRepo)
	result, err := usecase.GetUserByID(context.Background(), 1)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result == nil {
		t.Fatal("Expected result, got nil")
	}

	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}

	t.Log("✅ GetUserByID Success test passed!")
}

func TestGetUserByID_NotFound(t *testing.T) {
	mockRepo := &mock.SimpleMockUserRepository{
		User:  nil,
		Error: sql.ErrNoRows,
	}

	usecase := NewUserUseCase(mockRepo)
	result, err := usecase.GetUserByID(context.Background(), 999)

	if err != customerrors.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got: %v", err)
	}

	if result != nil {
		t.Error("Expected nil result")
	}

	t.Log("✅ GetUserByID NotFound test passed!")
}

func TestGetUserByID_DatabaseError(t *testing.T) {
	mockRepo := &mock.SimpleMockUserRepository{
		User:  nil,
		Error: sql.ErrConnDone,
	}

	usecase := NewUserUseCase(mockRepo)
	result, err := usecase.GetUserByID(context.Background(), 1)

	if err != sql.ErrConnDone {
		t.Errorf("Expected sql.ErrConnDone, got: %v", err)
	}

	if result != nil {
		t.Error("Expected nil result")
	}

	t.Log("✅ GetUserByID DatabaseError test passed!")
}

func TestDeleteUser_Success(t *testing.T) {
	fakeUser := &entity.User{
		ID:       1,
		Username: "testuser",
	}

	mockRepo := &mock.SimpleMockUserRepository{
		User:  fakeUser,
		Error: nil,
	}

	usecase := NewUserUseCase(mockRepo)
	err := usecase.DeleteUser(context.Background(), 1)

	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	t.Log("✅ DeleteUser Success test passed!")
}
