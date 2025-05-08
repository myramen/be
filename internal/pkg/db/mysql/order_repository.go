package mysql

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/myramen/be/internal/app/order"
	"github.com/myramen/be/internal/pkg/utils/errors"
)

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) order.Repository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(ctx context.Context, order *order.Order) error {
	// Options JSON으로 변환
	optionsJSON, err := json.Marshal(order.Options)
	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "주문 옵션을 JSON으로 변환하는데 실패했습니다.")
	}

	// 적용된 쿠폰 JSON으로 변환
	var appliedCouponJSON []byte
	if order.AppliedCoupon != nil {
		appliedCouponJSON, err = json.Marshal(order.AppliedCoupon)
		if err != nil {
			return errors.Internal("INTERNAL_ERROR", "적용된 쿠폰을 JSON으로 변환하는데 실패했습니다.")
		}
	}

	// 새 쿠폰 JSON으로 변환
	var newCouponJSON []byte
	if order.NewCoupon != nil {
		newCouponJSON, err = json.Marshal(order.NewCoupon)
		if err != nil {
			return errors.Internal("INTERNAL_ERROR", "새 쿠폰을 JSON으로 변환하는데 실패했습니다.")
		}
	}

	query := `
		INSERT INTO orders (
			order_id, name, account_number, quantity, spicy_level, 
			delivery_option, options, total_price, status, 
			applied_coupon, new_coupon, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err = r.db.ExecContext(
		ctx, query,
		order.OrderID, order.Name, order.AccountNumber, order.Quantity, order.SpicyLevel,
		order.DeliveryOption, optionsJSON, order.TotalPrice, order.Status,
		appliedCouponJSON, newCouponJSON, order.CreatedAt, order.UpdatedAt,
	)

	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "주문을 저장하는데 실패했습니다.")
	}

	return nil
}

func (r *orderRepository) FindByID(ctx context.Context, orderID string) (*order.Order, error) {
	query := `
		SELECT 
			order_id, name, account_number, quantity, spicy_level, 
			delivery_option, options, total_price, status, 
			applied_coupon, new_coupon, created_at, updated_at
		FROM orders 
		WHERE order_id = ?
	`

	var (
		orderResult       order.Order
		optionsJSON       []byte
		appliedCouponJSON sql.NullString
		newCouponJSON     sql.NullString
	)

	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&orderResult.OrderID, &orderResult.Name, &orderResult.AccountNumber,
		&orderResult.Quantity, &orderResult.SpicyLevel, &orderResult.DeliveryOption,
		&optionsJSON, &orderResult.TotalPrice, &orderResult.Status,
		&appliedCouponJSON, &newCouponJSON, &orderResult.CreatedAt, &orderResult.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Internal("INTERNAL_ERROR", "주문을 조회하는데 실패했습니다.")
	}

	// Options 파싱
	if err := json.Unmarshal(optionsJSON, &orderResult.Options); err != nil {
		return nil, errors.Internal("INTERNAL_ERROR", "주문 옵션을 파싱하는데 실패했습니다.")
	}

	// 적용된 쿠폰 파싱
	if appliedCouponJSON.Valid && appliedCouponJSON.String != "" {
		var appliedCoupon order.Coupon
		if err := json.Unmarshal([]byte(appliedCouponJSON.String), &appliedCoupon); err != nil {
			return nil, errors.Internal("INTERNAL_ERROR", "적용된 쿠폰을 파싱하는데 실패했습니다.")
		}
		orderResult.AppliedCoupon = &appliedCoupon
	}

	// 새 쿠폰 파싱
	if newCouponJSON.Valid && newCouponJSON.String != "" {
		var newCoupon order.Coupon
		if err := json.Unmarshal([]byte(newCouponJSON.String), &newCoupon); err != nil {
			return nil, errors.Internal("INTERNAL_ERROR", "새 쿠폰을 파싱하는데 실패했습니다.")
		}
		orderResult.NewCoupon = &newCoupon
	}

	return &orderResult, nil
}

func (r *orderRepository) FindAll(ctx context.Context) ([]order.Order, error) {
	query := `
		SELECT 
			order_id, name, account_number, quantity, spicy_level, 
			delivery_option, options, total_price, status, 
			applied_coupon, new_coupon, created_at, updated_at
		FROM orders
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Internal("INTERNAL_ERROR", "주문을 조회하는데 실패했습니다.")
	}
	defer rows.Close()

	var orders []order.Order

	for rows.Next() {
		var (
			orderItem         order.Order
			optionsJSON       []byte
			appliedCouponJSON sql.NullString
			newCouponJSON     sql.NullString
		)

		if err := rows.Scan(
			&orderItem.OrderID, &orderItem.Name, &orderItem.AccountNumber,
			&orderItem.Quantity, &orderItem.SpicyLevel, &orderItem.DeliveryOption,
			&optionsJSON, &orderItem.TotalPrice, &orderItem.Status,
			&appliedCouponJSON, &newCouponJSON, &orderItem.CreatedAt, &orderItem.UpdatedAt,
		); err != nil {
			return nil, errors.Internal("INTERNAL_ERROR", "주문 정보를 파싱하는데 실패했습니다.")
		}

		// Options 파싱
		if err := json.Unmarshal(optionsJSON, &orderItem.Options); err != nil {
			return nil, errors.Internal("INTERNAL_ERROR", "주문 옵션을 파싱하는데 실패했습니다.")
		}

		// 적용된 쿠폰 파싱
		if appliedCouponJSON.Valid && appliedCouponJSON.String != "" {
			var appliedCoupon order.Coupon
			if err := json.Unmarshal([]byte(appliedCouponJSON.String), &appliedCoupon); err != nil {
				return nil, errors.Internal("INTERNAL_ERROR", "적용된 쿠폰을 파싱하는데 실패했습니다.")
			}
			orderItem.AppliedCoupon = &appliedCoupon
		}

		// 새 쿠폰 파싱
		if newCouponJSON.Valid && newCouponJSON.String != "" {
			var newCoupon order.Coupon
			if err := json.Unmarshal([]byte(newCouponJSON.String), &newCoupon); err != nil {
				return nil, errors.Internal("INTERNAL_ERROR", "새 쿠폰을 파싱하는데 실패했습니다.")
			}
			orderItem.NewCoupon = &newCoupon
		}

		orders = append(orders, orderItem)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Internal("INTERNAL_ERROR", "주문 데이터를 처리하는데 실패했습니다.")
	}

	return orders, nil
}

func (r *orderRepository) UpdateStatus(ctx context.Context, orderID string, status string) error {
	query := `
		UPDATE orders 
		SET status = ?, updated_at = NOW() 
		WHERE order_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, status, orderID)
	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "주문 상태를 업데이트하는데 실패했습니다.")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "영향받은 행 수를 확인하는데 실패했습니다.")
	}

	if rows == 0 {
		return errors.NotFound("NOT_FOUND", "해당 주문을 찾을 수 없습니다.")
	}

	return nil
}

func (r *orderRepository) Delete(ctx context.Context, orderID string) error {
	query := "DELETE FROM orders WHERE order_id = ?"

	result, err := r.db.ExecContext(ctx, query, orderID)
	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "주문을 삭제하는데 실패했습니다.")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "영향받은 행 수를 확인하는데 실패했습니다.")
	}

	if rows == 0 {
		return errors.NotFound("NOT_FOUND", "해당 주문을 찾을 수 없습니다.")
	}

	return nil
}
