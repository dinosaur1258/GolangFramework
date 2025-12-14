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

// TestGetUserByID_TableDriven 使用表格驅動測試
func TestGetUserByID_TableDriven(t *testing.T) {
	// 定義測試案例表格
	testCases := []struct {
		name          string
		userID        int32
		mockUser      *entity.User
		mockError     error
		expectError   error
		expectNilUser bool
		description   string
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
			description:   "成功取得用戶",
		},
		{
			name:          "UserNotFound",
			userID:        999,
			mockUser:      nil,
			mockError:     sql.ErrNoRows,
			expectError:   customerrors.ErrUserNotFound,
			expectNilUser: true,
			description:   "用戶不存在",
		},
		{
			name:          "DatabaseError",
			userID:        1,
			mockUser:      nil,
			mockError:     sql.ErrConnDone,
			expectError:   sql.ErrConnDone,
			expectNilUser: true,
			description:   "資料庫連線錯誤",
		},
	}

	// 遍歷所有測試案例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			mockRepo := &mock.SimpleMockUserRepository{
				User:  tc.mockUser,
				Error: tc.mockError,
			}
			usecase := NewUserUseCase(mockRepo)

			// Act
			result, err := usecase.GetUserByID(context.Background(), tc.userID)

			// Assert
			if tc.expectError != nil {
				if err != tc.expectError {
					t.Errorf("[%s] Expected error %v, got %v", tc.description, tc.expectError, err)
				}
			} else {
				if err != nil {
					t.Errorf("[%s] Expected no error, got %v", tc.description, err)
				}
			}

			if tc.expectNilUser {
				if result != nil {
					t.Errorf("[%s] Expected nil user", tc.description)
				}
			} else {
				if result == nil {
					t.Errorf("[%s] Expected user, got nil", tc.description)
				} else if result.ID != tc.userID {
					t.Errorf("[%s] Expected user ID %d, got %d", tc.description, tc.userID, result.ID)
				}
			}

			t.Logf("✅ [%s] %s", tc.name, tc.description)
		})
	}
}
