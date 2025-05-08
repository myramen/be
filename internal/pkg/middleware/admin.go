package middleware

import (
	"net/http"

	"github.com/myramen/be/internal/pkg/config"
	"github.com/myramen/be/internal/pkg/utils/errors"

	"github.com/gin-gonic/gin"
)

// AdminAuth 관리자 인증 미들웨어
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		adminPassword := c.GetHeader("X-Admin-Password")
		if adminPassword == "" {
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
				Error:   "UNAUTHORIZED",
				Message: "관리자 인증이 필요합니다.",
			})
			c.Abort()
			return
		}

		if adminPassword != config.AppConfig.AdminPassword {
			c.JSON(http.StatusUnauthorized, errors.ErrorResponse{
				Error:   "UNAUTHORIZED",
				Message: "관리자 인증에 실패했습니다.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
