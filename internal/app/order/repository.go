package order

import (
	"context"
)

// Repository 주문 리포지토리 인터페이스
type Repository interface {
	// Create 새로운 주문 생성
	Create(ctx context.Context, order *Order) error
	
	// FindByID 주문 ID로 주문 조회
	FindByID(ctx context.Context, orderID string) (*Order, error)
	
	// FindAll 모든 주문 조회
	FindAll(ctx context.Context) ([]Order, error)
	
	// UpdateStatus 주문 상태 업데이트
	UpdateStatus(ctx context.Context, orderID string, status string) error
	
	// Delete 주문 삭제
	Delete(ctx context.Context, orderID string) error
}
