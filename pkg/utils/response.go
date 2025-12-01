package utils

import "github.com/gin-gonic/gin"

// Response 統一響應結構
type Response struct {
	Success bool         `json:"success"`
	Message string       `json:"message,omitempty"`
	Data    interface{}  `json:"data,omitempty"`
	Error   *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail 錯誤詳情
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// SuccessResponse 成功響應
func SuccessResponse(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse 錯誤響應
func ErrorResponse(c *gin.Context, statusCode int, code, message string, details ...string) {
	errorDetail := &ErrorDetail{
		Code:    code,
		Message: message,
	}

	if len(details) > 0 {
		errorDetail.Details = details[0]
	}

	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorDetail,
	})
}
