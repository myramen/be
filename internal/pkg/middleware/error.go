package middleware

import (
	"github.com/gin-gonic/gin"
)

// ErrorHandler 에러 처리 미들웨어
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
