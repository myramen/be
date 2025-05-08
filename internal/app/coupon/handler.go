package coupon

import (
	"net/http"

	"github.com/myramen/be/internal/pkg/middleware"
	"github.com/myramen/be/internal/pkg/utils/errors"

	"github.com/gin-gonic/gin"
)

// Handler 쿠폰 핸들러
type Handler struct {
	service *Service
}

// NewHandler 쿠폰 핸들러 생성
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes 라우트 등록
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	coupons := r.Group("/coupons")
	{
		coupons.GET("/:couponId", h.GetCouponByID)
	}

	admin := r.Group("/admin")
	{
		admin.Use(middleware.AdminAuth())
		admin.GET("/coupons", h.GetAllCoupons)
	}
}

// GetCouponByID 쿠폰 조회 핸들러
func (h *Handler) GetCouponByID(c *gin.Context) {
	couponID := c.Param("couponId")
	if couponID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "쿠폰 ID가 필요합니다.",
		})
		return
	}

	result, err := h.service.GetCouponByID(c, couponID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAllCoupons 모든 유효한 쿠폰 조회 핸들러
func (h *Handler) GetAllCoupons(c *gin.Context) {
	result, err := h.service.GetAllCoupons(c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
