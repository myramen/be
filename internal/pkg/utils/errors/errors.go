package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CustomError 커스텀 에러 구조체
type CustomError struct {
	Status  int         `json:"-"`
	Code    string      `json:"-"`
	Message string      `json:"-"`
	Details interface{} `json:"-"`
}

// Error 에러 메시지 반환
func (e CustomError) Error() string {
	return e.Message
}

// ErrorResponse 에러 응답 구조체
type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

// 에러 상수
const (
	StatusBadRequest     = http.StatusBadRequest
	StatusUnauthorized   = http.StatusUnauthorized
	StatusForbidden      = http.StatusForbidden
	StatusNotFound       = http.StatusNotFound
	StatusInternalServer = http.StatusInternalServerError
)

// NewError 새 에러 생성
func NewError(status int, code, message string, details interface{}) error {
	return CustomError{
		Status:  status,
		Code:    code,
		Message: message,
		Details: details,
	}
}

// BadRequest 잘못된 요청 에러
func BadRequest(code, message string) error {
	return NewError(StatusBadRequest, code, message, nil)
}

// Unauthorized 인증 실패 에러
func Unauthorized(code, message string) error {
	return NewError(StatusUnauthorized, code, message, nil)
}

// Forbidden 권한 없음 에러
func Forbidden(code, message string) error {
	return NewError(StatusForbidden, code, message, nil)
}

// NotFound 리소스 없음 에러
func NotFound(code, message string) error {
	return NewError(StatusNotFound, code, message, nil)
}

// Internal 내부 서버 에러
func Internal(code, message string) error {
	return NewError(StatusInternalServer, code, message, nil)
}

// HandleError 에러 처리 함수
func HandleError(c *gin.Context, err error) {
	if customError, ok := err.(CustomError); ok {
		c.JSON(customError.Status, ErrorResponse{
			Error:   customError.Code,
			Message: customError.Message,
			Details: customError.Details,
		})
		return
	}

	// 기본 에러 처리
	c.JSON(http.StatusInternalServerError, ErrorResponse{
		Error:   "INTERNAL_ERROR",
		Message: "Internal server error",
	})
}
