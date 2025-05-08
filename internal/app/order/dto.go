package order

// CreateOrderRequest 주문 생성 요청 DTO
type CreateOrderRequest struct {
	Name           string  `json:"name" binding:"required"`
	AccountNumber  string  `json:"accountNumber" binding:"required"`
	Quantity       int     `json:"quantity" binding:"required,min=1"`
	SpicyLevel     int     `json:"spicyLevel" binding:"min=1,max=5"`
	DeliveryOption string  `json:"deliveryOption" binding:"required"`
	Options        Options `json:"options"`
	CouponID       string  `json:"couponId,omitempty"`
}

// OrderResponse 주문 응답 DTO
type OrderResponse struct {
	OrderID        string  `json:"orderId"`
	Name           string  `json:"name"`
	AccountNumber  string  `json:"accountNumber"`
	Quantity       int     `json:"quantity"`
	SpicyLevel     int     `json:"spicyLevel"`
	DeliveryOption string  `json:"deliveryOption"`
	Options        Options `json:"options"`
	TotalPrice     int     `json:"totalPrice"`
	Status         string  `json:"status"`
	AppliedCoupon  *Coupon `json:"appliedCoupon,omitempty"`
	NewCoupon      *Coupon `json:"newCoupon,omitempty"`
}

// OrderListResponse 주문 목록 응답 DTO
type OrderListResponse struct {
	Orders []OrderResponse `json:"orders"`
}

// UpdateOrderStatusRequest 주문 상태 변경 요청 DTO
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=PENDING PAID COOKING READY DELIVERING DELIVERED"`
}

// ErrorResponse 에러 응답 DTO
type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
