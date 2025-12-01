package errors

import "errors"

// 錯誤定義
var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidInput       = errors.New("invalid input")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrForbidden          = errors.New("forbidden")
	ErrInternalServer     = errors.New("internal server error")
)

// 錯誤代碼（用於 API 響應）
const (
	CodeUserNotFound       = "USER_NOT_FOUND"
	CodeUserAlreadyExists  = "USER_ALREADY_EXISTS"
	CodeInvalidCredentials = "INVALID_CREDENTIALS"
	CodeInvalidInput       = "INVALID_INPUT"
	CodeUnauthorized       = "UNAUTHORIZED"
	CodeForbidden          = "FORBIDDEN"
	CodeInternalServer     = "INTERNAL_SERVER_ERROR"
	CodeValidationFailed   = "VALIDATION_FAILED"
)

// 錯誤訊息
const (
	MsgUserNotFound       = "User not found"
	MsgUserAlreadyExists  = "User already exists"
	MsgInvalidCredentials = "Invalid email or password"
	MsgInvalidInput       = "Invalid input data"
	MsgUnauthorized       = "Unauthorized access"
	MsgForbidden          = "Access forbidden"
	MsgInternalServer     = "Internal server error"
	MsgValidationFailed   = "Validation failed"
)
