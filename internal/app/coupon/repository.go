package coupon

import (
	"context"
)

// Repository 쿠폰 리포지토리 인터페이스
type Repository interface {
	// Create 새로운 쿠폰 생성
	Create(ctx context.Context, coupon *Coupon) error
	
	// FindByID 쿠폰 ID로 쿠폰 조회
	FindByID(ctx context.Context, couponID string) (*Coupon, error)
	
	// FindAll 모든 유효한 쿠폰 조회
	FindAll(ctx context.Context) ([]Coupon, error)
	
	// Update 쿠폰 정보 업데이트
	Update(ctx context.Context, coupon *Coupon) error
	
	// MarkAsUsed 쿠폰 사용 처리
	MarkAsUsed(ctx context.Context, couponID string) error
}
