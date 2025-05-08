package order

import (
	"time"
)

type Order struct {
	OrderID        string    `json:"orderId"`
	Name           string    `json:"name"`
	AccountNumber  string    `json:"accountNumber"`
	Quantity       int       `json:"quantity"`
	SpicyLevel     int       `json:"spicyLevel"`
	DeliveryOption string    `json:"deliveryOption"`
	Options        Options   `json:"options"`
	TotalPrice     int       `json:"totalPrice"`
	Status         string    `json:"status"`
	AppliedCoupon  *Coupon   `json:"appliedCoupon,omitempty"`
	NewCoupon      *Coupon   `json:"newCoupon,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Options struct {
	Chopsticks       bool `json:"chopsticks"`
	HotWaterDelivery bool `json:"hotWaterDelivery"`
	CookingService   bool `json:"cookingService"`
}

type Coupon struct {
	CouponID   string    `json:"couponId"`
	Discount   int       `json:"discount"`
	ExpiryDate time.Time `json:"expiryDate,omitempty"`
}

// Status 상수 정의
const (
	StatusPending    = "PENDING"
	StatusPaid       = "PAID"
	StatusCooking    = "COOKING"
	StatusReady      = "READY"
	StatusDelivering = "DELIVERING"
	StatusDelivered  = "DELIVERED"
)

// DeliveryOption 상수 정의
const (
	DeliveryOptionPickup4F      = "PICKUP_4F"
	DeliveryOptionPickupLaundry = "PICKUP_LAUNDRY"
	DeliveryOptionDelivery      = "DELIVERY"
)

// 가격 상수 정의
const (
	BasePrice           = 4000
	HotWaterServiceFee  = 500
	CookingServiceFee   = 500
	DefaultCouponAmount = 200
	CouponThreshold     = 3
)
