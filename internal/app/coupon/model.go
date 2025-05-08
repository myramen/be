package coupon

import (
	"time"
)

type Coupon struct {
	CouponID   string    `json:"couponId"`
	Discount   int       `json:"discount"`
	ExpiryDate time.Time `json:"expiryDate"`
	IsUsed     bool      `json:"isUsed"`
	IssuedAt   time.Time `json:"issuedAt"`
}

const (
	DefaultDiscount = 200
	ExpiryDuration  = 30 * 24 * time.Hour // 30일
)
