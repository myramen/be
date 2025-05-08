package coupon

import (
	"time"
)

// CouponResponse 쿠폰 응답 DTO
type CouponResponse struct {
	CouponID   string    `json:"couponId"`
	Discount   int       `json:"discount"`
	ExpiryDate time.Time `json:"expiryDate"`
	IsUsed     bool      `json:"isUsed"`
	IssuedAt   time.Time `json:"issuedAt"`
}

// CouponListResponse 쿠폰 목록 응답 DTO
type CouponListResponse struct {
	Coupons []CouponResponse `json:"coupons"`
}

// ErrorResponse 에러 응답 DTO
type ErrorResponse struct {
	Error   string      `json:"error"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}
