package usecase

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/dinosaur1258/GolangFramework/internal/domain/dto/request"
	"github.com/dinosaur1258/GolangFramework/internal/domain/entity"
	"github.com/dinosaur1258/GolangFramework/internal/repository/mock"
	customerrors "github.com/dinosaur1258/GolangFramework/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// =============================================================================
// GetUserByID Tests
// =============================================================================

func TestGetUserByID(t *testing.T) {
	testCases := []struct {
		name          string
		userID        int32
		mockUser      *entity.User
		mockError     error
		expectError   error
		expectNilUser bool
	}{
		{
			name:   "Success",
			userID: 1,
			mockUser: &entity.User{
				ID:           1,
				Username:     "testuser",
				Email:        "test@example.com",
				PasswordHash: "hash123",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			mockError:     nil,
			expectError:   nil,
			expectNilUser: false,
		},
		{
			name:          "UserNotFound",
			userID:        999,
			mockUser:      nil,
			mockError:     sql.ErrNoRows,
			expectError:   customerrors.ErrUserNotFound,
			expectNilUser: true,
		},
		{
			name:          "DatabaseError",
			userID:        1,
			mockUser:      nil,
			mockError:     sql.ErrConnDone,
			expectError:   sql.ErrConnDone,
			expectNilUser: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mock.SimpleMockUserRepository{
				User:  tc.mockUser,
				Error: tc.mockError,
			}
			usecase := NewUserUseCase(mockRepo)

			result, err := usecase.GetUserByID(context.Background(), tc.userID)

			if tc.expectError != nil {
				if err != tc.expectError {
					t.Errorf("Expected error %v, got %v", tc.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}

			if tc.expectNilUser {
				if result != nil {
					t.Error("Expected nil user")
				}
			} else {
				if result == nil {
					t.Error("Expected user, got nil")
				} else if result.ID != tc.userID {
					t.Errorf("Expected user ID %d, got %d", tc.userID, result.ID)
				}
			}
		})
	}
}

// =============================================================================
// UpdateUser Tests
// =============================================================================

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name        string
		userID      int32
		request     request.UpdateUserRequest
		setupMock   func() *mock.SimpleMockUserRepository
		expectError error
	}{
		{
			name:   "Success",
			userID: 1,
			request: request.UpdateUserRequest{
				Username: "newname",
				Email:    "new@example.com",
			},
			setupMock: func() *mock.SimpleMockUserRepository {
				existingUser := &entity.User{
					ID:       1,
					Username: "oldname",
					Email:    "old@example.com",
				}

				return &mock.SimpleMockUserRepository{
					GetByIDFunc: func(ctx context.Context, id int32) (*entity.User, error) {
						if id == 1 {
							return existingUser, nil
						}
						return nil, sql.ErrNoRows
					},
					GetByEmailFunc: func(ctx context.Context, email string) (*entity.User, error) {
						if email == "new@example.com" {
							return nil, sql.ErrNoRows // 新 Email 不存在
						}
						return existingUser, nil
					},
					GetByUsernameFunc: func(ctx context.Context, username string) (*entity.User, error) {
						if username == "newname" {
							return nil, sql.ErrNoRows // 新 Username 不存在
						}
						return existingUser, nil
					},
					Error: nil,
				}
			},
			expectError: nil,
		},
		{
			name:   "UserNotFound",
			userID: 999,
			request: request.UpdateUserRequest{
				Username: "newname",
			},
			setupMock: func() *mock.SimpleMockUserRepository {
				return &mock.SimpleMockUserRepository{
					User:  nil,
					Error: sql.ErrNoRows,
				}
			},
			expectError: customerrors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := tc.setupMock()
			usecase := NewUserUseCase(mockRepo)

			result, err := usecase.UpdateUser(context.Background(), tc.userID, tc.request)

			if tc.expectError != nil {
				if err != tc.expectError {
					t.Errorf("Expected error %v, got %v", tc.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result == nil {
					t.Error("Expected result, got nil")
				}
			}
		})
	}
}

// =============================================================================
// DeleteUser Tests
// =============================================================================

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name        string
		userID      int32
		mockUser    *entity.User
		mockError   error
		expectError error
	}{
		{
			name:   "Success",
			userID: 1,
			mockUser: &entity.User{
				ID:       1,
				Username: "testuser",
			},
			mockError:   nil,
			expectError: nil,
		},
		{
			name:        "UserNotFound",
			userID:      999,
			mockUser:    nil,
			mockError:   sql.ErrNoRows,
			expectError: customerrors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mock.SimpleMockUserRepository{
				User:  tc.mockUser,
				Error: tc.mockError,
			}
			usecase := NewUserUseCase(mockRepo)

			err := usecase.DeleteUser(context.Background(), tc.userID)

			if tc.expectError != nil {
				if err != tc.expectError {
					t.Errorf("Expected error %v, got %v", tc.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}

// =============================================================================
// ListUsers Tests
// =============================================================================

func TestListUsers(t *testing.T) {
	testCases := []struct {
		name        string
		page        int
		limit       int
		mockUser    *entity.User
		mockError   error
		expectError error
		expectCount int
	}{
		{
			name:  "Success",
			page:  1,
			limit: 10,
			mockUser: &entity.User{
				ID:       1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			mockError:   nil,
			expectError: nil,
			expectCount: 1,
		},
		{
			name:        "DefaultValues",
			page:        0,
			limit:       0,
			mockUser:    &entity.User{ID: 1},
			mockError:   nil,
			expectError: nil,
			expectCount: 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mock.SimpleMockUserRepository{
				User:  tc.mockUser,
				Error: tc.mockError,
			}
			usecase := NewUserUseCase(mockRepo)

			result, err := usecase.ListUsers(context.Background(), tc.page, tc.limit)

			if tc.expectError != nil {
				if err != tc.expectError {
					t.Errorf("Expected error %v, got %v", tc.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
				if result != nil && len(result) != tc.expectCount {
					t.Errorf("Expected %d users, got %d", tc.expectCount, len(result))
				}
			}
		})
	}
}

// =============================================================================
// ChangePassword Tests
// =============================================================================

func TestChangePassword(t *testing.T) {
	oldPasswordHash, _ := bcrypt.GenerateFromPassword([]byte("oldpass123"), bcrypt.DefaultCost)

	testCases := []struct {
		name        string
		userID      int32
		request     request.ChangePasswordRequest
		mockUser    *entity.User
		mockError   error
		expectError error
	}{
		{
			name:   "Success",
			userID: 1,
			request: request.ChangePasswordRequest{
				OldPassword: "oldpass123",
				NewPassword: "newpass456",
			},
			mockUser: &entity.User{
				ID:           1,
				Username:     "testuser",
				PasswordHash: string(oldPasswordHash),
			},
			mockError:   nil,
			expectError: nil,
		},
		{
			name:   "WrongOldPassword",
			userID: 1,
			request: request.ChangePasswordRequest{
				OldPassword: "wrongpass",
				NewPassword: "newpass456",
			},
			mockUser: &entity.User{
				ID:           1,
				Username:     "testuser",
				PasswordHash: string(oldPasswordHash),
			},
			mockError:   nil,
			expectError: customerrors.ErrInvalidCredentials,
		},
		{
			name:   "UserNotFound",
			userID: 999,
			request: request.ChangePasswordRequest{
				OldPassword: "oldpass",
				NewPassword: "newpass",
			},
			mockUser:    nil,
			mockError:   sql.ErrNoRows,
			expectError: customerrors.ErrUserNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := &mock.SimpleMockUserRepository{
				User:  tc.mockUser,
				Error: tc.mockError,
			}
			usecase := NewUserUseCase(mockRepo)

			err := usecase.ChangePassword(context.Background(), tc.userID, tc.request)

			if tc.expectError != nil {
				if err != tc.expectError {
					t.Errorf("Expected error %v, got %v", tc.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got %v", err)
				}
			}
		})
	}
}
