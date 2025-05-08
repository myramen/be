package coupon

import (
	"context"

	"github.com/myramen/be/internal/pkg/utils/errors"
)

// Service 쿠폰 서비스
type Service struct {
	repo Repository
}

// NewService 쿠폰 서비스 생성
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// GetCouponByID 쿠폰 ID로 쿠폰 조회
func (s *Service) GetCouponByID(ctx context.Context, couponID string) (*CouponResponse, error) {
	coupon, err := s.repo.FindByID(ctx, couponID)
	if err != nil {
		return nil, err
	}

	if coupon == nil {
		return nil, errors.NotFound("NOT_FOUND", "해당 쿠폰을 찾을 수 없습니다.")
	}

	return &CouponResponse{
		CouponID:   coupon.CouponID,
		Discount:   coupon.Discount,
		ExpiryDate: coupon.ExpiryDate,
		IsUsed:     coupon.IsUsed,
		IssuedAt:   coupon.IssuedAt,
	}, nil
}

// GetAllCoupons 모든 유효한 쿠폰 조회
func (s *Service) GetAllCoupons(ctx context.Context) (*CouponListResponse, error) {
	coupons, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	couponResponses := make([]CouponResponse, 0, len(coupons))
	for _, coupon := range coupons {
		couponResponses = append(couponResponses, CouponResponse{
			CouponID:   coupon.CouponID,
			Discount:   coupon.Discount,
			ExpiryDate: coupon.ExpiryDate,
			IsUsed:     coupon.IsUsed,
			IssuedAt:   coupon.IssuedAt,
		})
	}

	return &CouponListResponse{
		Coupons: couponResponses,
	}, nil
}

// CreateCoupon 쿠폰 생성 (시스템에서만 사용)
func (s *Service) CreateCoupon(ctx context.Context, coupon *Coupon) error {
	return s.repo.Create(ctx, coupon)
}

// UseCoupon 쿠폰 사용 처리
func (s *Service) UseCoupon(ctx context.Context, couponID string) error {
	coupon, err := s.repo.FindByID(ctx, couponID)
	if err != nil {
		return err
	}

	if coupon == nil {
		return errors.NotFound("NOT_FOUND", "해당 쿠폰을 찾을 수 없습니다.")
	}

	if coupon.IsUsed {
		return errors.BadRequest("INVALID_COUPON", "이미 사용된 쿠폰입니다.")
	}

	return s.repo.MarkAsUsed(ctx, couponID)
}
