package mysql

import (
	"context"
	"database/sql"

	"github.com/myramen/be/internal/app/coupon"
	"github.com/myramen/be/internal/pkg/utils/errors"
)

type couponRepository struct {
	db *sql.DB
}

func NewCouponRepository(db *sql.DB) coupon.Repository {
	return &couponRepository{db: db}
}

func (r *couponRepository) Create(ctx context.Context, coupon *coupon.Coupon) error {
	query := `
		INSERT INTO coupons (
			coupon_id, discount, expiry_date, is_used, issued_at
		) VALUES (?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(
		ctx, query,
		coupon.CouponID, coupon.Discount, coupon.ExpiryDate, coupon.IsUsed, coupon.IssuedAt,
	)

	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "쿠폰을 저장하는데 실패했습니다.")
	}

	return nil
}

func (r *couponRepository) FindByID(ctx context.Context, couponID string) (*coupon.Coupon, error) {
	query := `
		SELECT coupon_id, discount, expiry_date, is_used, issued_at
		FROM coupons 
		WHERE coupon_id = ?
	`

	var couponResult coupon.Coupon

	err := r.db.QueryRowContext(ctx, query, couponID).Scan(
		&couponResult.CouponID, &couponResult.Discount, &couponResult.ExpiryDate,
		&couponResult.IsUsed, &couponResult.IssuedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, errors.Internal("INTERNAL_ERROR", "쿠폰을 조회하는데 실패했습니다.")
	}

	return &couponResult, nil
}

func (r *couponRepository) FindAll(ctx context.Context) ([]coupon.Coupon, error) {
	query := `
		SELECT coupon_id, discount, expiry_date, is_used, issued_at
		FROM coupons
		WHERE is_used = FALSE AND expiry_date > NOW()
		ORDER BY issued_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.Internal("INTERNAL_ERROR", "쿠폰을 조회하는데 실패했습니다.")
	}
	defer rows.Close()

	var coupons []coupon.Coupon

	for rows.Next() {
		var couponItem coupon.Coupon

		if err := rows.Scan(
			&couponItem.CouponID, &couponItem.Discount, &couponItem.ExpiryDate,
			&couponItem.IsUsed, &couponItem.IssuedAt,
		); err != nil {
			return nil, errors.Internal("INTERNAL_ERROR", "쿠폰 정보를 파싱하는데 실패했습니다.")
		}

		coupons = append(coupons, couponItem)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Internal("INTERNAL_ERROR", "쿠폰 데이터를 처리하는데 실패했습니다.")
	}

	return coupons, nil
}

func (r *couponRepository) Update(ctx context.Context, coupon *coupon.Coupon) error {
	query := `
		UPDATE coupons 
		SET discount = ?, expiry_date = ?, is_used = ?
		WHERE coupon_id = ?
	`

	result, err := r.db.ExecContext(
		ctx, query,
		coupon.Discount, coupon.ExpiryDate, coupon.IsUsed, coupon.CouponID,
	)

	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "쿠폰을 업데이트하는데 실패했습니다.")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "영향받은 행 수를 확인하는데 실패했습니다.")
	}

	if rows == 0 {
		return errors.NotFound("NOT_FOUND", "해당 쿠폰을 찾을 수 없습니다.")
	}

	return nil
}

func (r *couponRepository) MarkAsUsed(ctx context.Context, couponID string) error {
	query := `
		UPDATE coupons 
		SET is_used = TRUE
		WHERE coupon_id = ?
	`

	result, err := r.db.ExecContext(ctx, query, couponID)
	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "쿠폰 사용 처리에 실패했습니다.")
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return errors.Internal("INTERNAL_ERROR", "영향받은 행 수를 확인하는데 실패했습니다.")
	}

	if rows == 0 {
		return errors.NotFound("NOT_FOUND", "해당 쿠폰을 찾을 수 없습니다.")
	}

	return nil
}
