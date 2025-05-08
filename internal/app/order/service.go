package order

import (
	"context"
	"fmt"
	"time"

	"github.com/myramen/be/internal/app/coupon"
	"github.com/myramen/be/internal/pkg/utils/errors"
)

// Service 주문 서비스 인터페이스
type Service struct {
	orderRepo  Repository
	couponRepo coupon.Repository
}

// NewService 주문 서비스 생성
func NewService(orderRepo Repository, couponRepo coupon.Repository) *Service {
	return &Service{
		orderRepo:  orderRepo,
		couponRepo: couponRepo,
	}
}

// CreateOrder 새로운 주문 생성
func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*OrderResponse, error) {
	// 기본 주문 정보 설정
	newOrder := &Order{
		OrderID:        generateOrderID(),
		Name:           req.Name,
		AccountNumber:  req.AccountNumber,
		Quantity:       req.Quantity,
		SpicyLevel:     req.SpicyLevel,
		DeliveryOption: req.DeliveryOption,
		Options:        req.Options,
		Status:         StatusPending,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// 가격 계산
	newOrder.TotalPrice = calculateTotalPrice(newOrder)

	// 쿠폰 적용 처리
	if req.CouponID != "" {
		couponData, err := s.couponRepo.FindByID(ctx, req.CouponID)
		if err != nil {
			return nil, err
		}

		if couponData == nil {
			return nil, errors.BadRequest("INVALID_COUPON", "사용할 수 없는 쿠폰입니다.")
		}

		if couponData.IsUsed {
			return nil, errors.BadRequest("INVALID_COUPON", "이미 사용된 쿠폰입니다.")
		}

		if time.Now().After(couponData.ExpiryDate) {
			return nil, errors.BadRequest("INVALID_COUPON", "만료된 쿠폰입니다.")
		}

		// 쿠폰 할인 적용
		newOrder.AppliedCoupon = &Coupon{
			CouponID: couponData.CouponID,
			Discount: couponData.Discount,
		}

		newOrder.TotalPrice -= couponData.Discount

		// 쿠폰 사용 처리
		if err := s.couponRepo.MarkAsUsed(ctx, req.CouponID); err != nil {
			return nil, err
		}
	}

	// 3개 이상 주문 시 신규 쿠폰 발급
	if req.Quantity >= CouponThreshold {
		newCoupon := &coupon.Coupon{
			CouponID:   generateCouponID(),
			Discount:   DefaultCouponAmount,
			ExpiryDate: time.Now().Add(coupon.ExpiryDuration),
			IsUsed:     false,
			IssuedAt:   time.Now(),
		}

		if err := s.couponRepo.Create(ctx, newCoupon); err != nil {
			return nil, err
		}

		newOrder.NewCoupon = &Coupon{
			CouponID:   newCoupon.CouponID,
			Discount:   newCoupon.Discount,
			ExpiryDate: newCoupon.ExpiryDate,
		}
	}

	// 주문 저장
	if err := s.orderRepo.Create(ctx, newOrder); err != nil {
		return nil, err
	}

	// 응답 생성
	return &OrderResponse{
		OrderID:        newOrder.OrderID,
		Name:           newOrder.Name,
		AccountNumber:  newOrder.AccountNumber,
		Quantity:       newOrder.Quantity,
		SpicyLevel:     newOrder.SpicyLevel,
		DeliveryOption: newOrder.DeliveryOption,
		Options:        newOrder.Options,
		TotalPrice:     newOrder.TotalPrice,
		Status:         newOrder.Status,
		AppliedCoupon:  newOrder.AppliedCoupon,
		NewCoupon:      newOrder.NewCoupon,
	}, nil
}

// GetOrderByID 주문 ID로 주문 조회
func (s *Service) GetOrderByID(ctx context.Context, orderID string) (*OrderResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.NotFound("NOT_FOUND", "해당 주문을 찾을 수 없습니다.")
	}

	return &OrderResponse{
		OrderID:        order.OrderID,
		Name:           order.Name,
		AccountNumber:  order.AccountNumber,
		Quantity:       order.Quantity,
		SpicyLevel:     order.SpicyLevel,
		DeliveryOption: order.DeliveryOption,
		Options:        order.Options,
		TotalPrice:     order.TotalPrice,
		Status:         order.Status,
		AppliedCoupon:  order.AppliedCoupon,
	}, nil
}

// GetAllOrders 모든 주문 조회
func (s *Service) GetAllOrders(ctx context.Context) (*OrderListResponse, error) {
	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	orderResponses := make([]OrderResponse, 0, len(orders))
	for _, order := range orders {
		orderResponses = append(orderResponses, OrderResponse{
			OrderID:        order.OrderID,
			Name:           order.Name,
			AccountNumber:  order.AccountNumber,
			Quantity:       order.Quantity,
			SpicyLevel:     order.SpicyLevel,
			DeliveryOption: order.DeliveryOption,
			Options:        order.Options,
			TotalPrice:     order.TotalPrice,
			Status:         order.Status,
			AppliedCoupon:  order.AppliedCoupon,
		})
	}

	return &OrderListResponse{
		Orders: orderResponses,
	}, nil
}

// UpdateOrderStatus 주문 상태 업데이트
func (s *Service) UpdateOrderStatus(ctx context.Context, orderID string, status string) (*OrderResponse, error) {
	order, err := s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, errors.NotFound("NOT_FOUND", "해당 주문을 찾을 수 없습니다.")
	}

	// 상태 유효성 검사 - 이미 handler에서 확인함

	// 상태 업데이트
	if err := s.orderRepo.UpdateStatus(ctx, orderID, status); err != nil {
		return nil, err
	}

	// 업데이트된 주문 정보 조회
	order, err = s.orderRepo.FindByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	return &OrderResponse{
		OrderID:        order.OrderID,
		Name:           order.Name,
		AccountNumber:  order.AccountNumber,
		Quantity:       order.Quantity,
		SpicyLevel:     order.SpicyLevel,
		DeliveryOption: order.DeliveryOption,
		Options:        order.Options,
		TotalPrice:     order.TotalPrice,
		Status:         order.Status,
		AppliedCoupon:  order.AppliedCoupon,
	}, nil
}

// 가격 계산 함수
func calculateTotalPrice(order *Order) int {
	// 기본 라면 가격
	price := BasePrice * order.Quantity

	// 추가 옵션 가격
	if order.Options.HotWaterDelivery {
		price += HotWaterServiceFee
	}

	if order.Options.CookingService {
		price += CookingServiceFee
	}

	// 쿠폰 할인은 이 함수 밖에서 처리

	return price
}

// 주문 ID 생성 함수
func generateOrderID() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("o%d", timestamp)
}

// 쿠폰 ID 생성 함수
func generateCouponID() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	return fmt.Sprintf("c%d", timestamp)
}
