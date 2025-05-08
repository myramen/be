package order

import (
	"net/http"

	"github.com/myramen/be/internal/pkg/middleware"
	"github.com/myramen/be/internal/pkg/utils/errors"

	"github.com/gin-gonic/gin"
)

// Handler 주문 핸들러
type Handler struct {
	service *Service
}

// NewHandler 주문 핸들러 생성
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// RegisterRoutes 라우트 등록
func (h *Handler) RegisterRoutes(r *gin.RouterGroup) {
	orders := r.Group("/orders")
	{
		orders.POST("", h.CreateOrder)
		orders.GET("/:orderId", h.GetOrderByID)
	}

	admin := r.Group("/admin")
	{
		admin.Use(middleware.AdminAuth())
		admin.GET("/orders", h.GetAllOrders)
		admin.PUT("/orders/:orderId/status", h.UpdateOrderStatus)
	}
}

// CreateOrder 주문 생성 핸들러
func (h *Handler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "구매 요청 정보가 유효하지 않습니다.",
		})
		return
	}

	// 기본값 설정
	if req.SpicyLevel == 0 {
		req.SpicyLevel = 3 // 기본 매운맛 레벨
	}

	result, err := h.service.CreateOrder(c, req)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetOrderByID 주문 조회 핸들러
func (h *Handler) GetOrderByID(c *gin.Context) {
	orderID := c.Param("orderId")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "주문 ID가 필요합니다.",
		})
		return
	}

	result, err := h.service.GetOrderByID(c, orderID)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetAllOrders 모든 주문 조회 핸들러
func (h *Handler) GetAllOrders(c *gin.Context) {
	result, err := h.service.GetAllOrders(c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// UpdateOrderStatus 주문 상태 업데이트 핸들러
func (h *Handler) UpdateOrderStatus(c *gin.Context) {
	orderID := c.Param("orderId")
	if orderID == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_REQUEST",
			Message: "주문 ID가 필요합니다.",
		})
		return
	}

	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "INVALID_STATUS",
			Message: "유효하지 않은 주문 상태입니다.",
		})
		return
	}

	result, err := h.service.UpdateOrderStatus(c, orderID, req.Status)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
