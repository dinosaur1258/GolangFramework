package utils

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// 設定 Gin 為測試模式
func init() {
	gin.SetMode(gin.TestMode)
}

// TestSuccessResponse 測試成功響應
func TestSuccessResponse(t *testing.T) {
	// 1. 準備測試環境
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 2. 執行要測試的函數
	testData := map[string]string{"id": "123", "name": "Test"}
	SuccessResponse(c, http.StatusOK, "success", testData)

	// 3. 驗證結果
	// 3.1 檢查 HTTP 狀態碼
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// 3.2 解析響應 JSON
	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// 3.3 驗證響應結構
	if !response.Success {
		t.Error("Expected success to be true")
	}

	if response.Message != "success" {
		t.Errorf("Expected message 'success', got '%s'", response.Message)
	}

	if response.Data == nil {
		t.Error("Expected data to be present")
	}
}

// TestErrorResponse 測試錯誤響應
func TestErrorResponse(t *testing.T) {
	// 準備
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// 執行
	ErrorResponse(c, http.StatusBadRequest, "INVALID_INPUT", "Invalid data", "Field 'email' is required")

	// 驗證
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	// 驗證錯誤響應
	if response.Success {
		t.Error("Expected success to be false")
	}

	if response.Error == nil {
		t.Fatal("Expected error to be present")
	}

	if response.Error.Code != "INVALID_INPUT" {
		t.Errorf("Expected error code 'INVALID_INPUT', got '%s'", response.Error.Code)
	}

	if response.Error.Message != "Invalid data" {
		t.Errorf("Expected error message 'Invalid data', got '%s'", response.Error.Message)
	}

	if response.Error.Details != "Field 'email' is required" {
		t.Errorf("Expected error details 'Field 'email' is required', got '%s'", response.Error.Details)
	}
}

// TestErrorResponseWithoutDetails 測試沒有詳細信息的錯誤響應
func TestErrorResponseWithoutDetails(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	ErrorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Something went wrong")

	var response Response
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Error.Details != "" {
		t.Errorf("Expected empty details, got '%s'", response.Error.Details)
	}
}
