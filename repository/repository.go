package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Puneet-Vishnoi/Coupon-System/db/postgres/providers"
	"github.com/Puneet-Vishnoi/Coupon-System/models"
)

type CouponRepository struct {
	DBHelper *providers.DBHelper
}

func NewCouponRepository(db *providers.DBHelper) *CouponRepository {
	return &CouponRepository{DBHelper: db}
}

func (r *CouponRepository) CreateCoupon(ctx context.Context, c *models.Coupon) error {
	meds, err := json.Marshal(c.ApplicableMedicineIDs)
	if err != nil {
		return fmt.Errorf("failed to marshal medicine IDs: %w", err)
	}

	if c.DiscountValue<0{
		return fmt.Errorf("discount cant be negetive")

	}

	cats, err := json.Marshal(c.ApplicableCategories)
	if err != nil {
		return fmt.Errorf("failed to marshal categories: %w", err)
	}

	_, err = r.DBHelper.PostgresClient.ExecContext(ctx, `
	INSERT INTO coupons (
		coupon_code,
		discount_type,
		discount_value,
		discount_target,
		min_order_value,
		max_usage_per_user,
		expiry_date,
		applicable_medicine_ids,
		applicable_categories,
		usage_type,
		valid_start,
		valid_end,
		terms_and_conditions,
		max_discount_amount
	) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
`,
		c.CouponCode,
		c.DiscountType,
		c.DiscountValue,
		c.DiscountTarget,
		c.MinOrderValue,
		c.MaxUsagePerUser,
		c.ExpiryDate,
		meds,
		cats,
		c.UsageType,
		c.ValidTimeWindow.Start,
		c.ValidTimeWindow.End,
		c.TermsAndConditions,
		c.MaxDiscountAmount,
	)

	if err != nil {
		return fmt.Errorf("failed to insert coupon: %w", err)
	}
	return nil
}

func (r *CouponRepository) GetAllCoupons(ctx context.Context) ([]*models.Coupon, error) {
	rows, err := r.DBHelper.PostgresClient.QueryContext(ctx, `
        SELECT 
            coupon_code, expiry_date, usage_type, 
            applicable_medicine_ids, applicable_categories, 
            min_order_value, valid_start, valid_end, 
            terms_and_conditions, discount_type, discount_value, 
            max_usage_per_user, discount_target, max_discount_amount
        FROM coupons
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []*models.Coupon

	for rows.Next() {
		var c models.Coupon
		var meds, cats []byte

		err := rows.Scan(
			&c.CouponCode, &c.ExpiryDate, &c.UsageType,
			&meds, &cats,
			&c.MinOrderValue, &c.ValidTimeWindow.Start, &c.ValidTimeWindow.End,
			&c.TermsAndConditions, &c.DiscountType, &c.DiscountValue,
			&c.MaxUsagePerUser, &c.DiscountTarget, &c.MaxDiscountAmount,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(meds, &c.ApplicableMedicineIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal medicine IDs: %w", err)
		}
		if err := json.Unmarshal(cats, &c.ApplicableCategories); err != nil {
			return nil, fmt.Errorf("failed to unmarshal categories: %w", err)
		}

		coupons = append(coupons, &c)
	}
	return coupons, nil
}

func (r *CouponRepository) GetCouponByCode(ctx context.Context, tx *sql.Tx, code string) (models.Coupon, error) {
	var c models.Coupon
	var meds, cats []byte

	err := r.DBHelper.PostgresClient.QueryRowContext(ctx, `
        SELECT 
            coupon_code, expiry_date, usage_type, 
            applicable_medicine_ids, applicable_categories, 
            min_order_value, valid_start, valid_end, 
            terms_and_conditions, discount_type, discount_value, 
            max_usage_per_user, discount_target, max_discount_amount
        FROM coupons
        WHERE coupon_code = $1
				For UPDATE
    `, code).Scan(
		&c.CouponCode, &c.ExpiryDate, &c.UsageType,
		&meds, &cats,
		&c.MinOrderValue, &c.ValidTimeWindow.Start, &c.ValidTimeWindow.End,
		&c.TermsAndConditions, &c.DiscountType, &c.DiscountValue,
		&c.MaxUsagePerUser, &c.DiscountTarget, &c.MaxDiscountAmount,
	)
	if err != nil {
		return c, err
	}

	if err := json.Unmarshal(meds, &c.ApplicableMedicineIDs); err != nil {
		return c, fmt.Errorf("failed to unmarshal medicine IDs: %w", err)
	}
	if err := json.Unmarshal(cats, &c.ApplicableCategories); err != nil {
		return c, fmt.Errorf("failed to unmarshal categories: %w", err)
	}

	return c, nil
}

func (r *CouponRepository) GetUserUsageCount(ctx context.Context, tx *sql.Tx, userID, couponCode string) (int, error) {
	var count int
	// Removed FOR UPDATE as it's unnecessary for COUNT
	err := tx.QueryRowContext(ctx, `
        SELECT COUNT(*) FROM coupon_usages
        WHERE user_id = $1 AND coupon_code = $2
    `, userID, couponCode).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *CouponRepository) RecordUsage(ctx context.Context, tx *sql.Tx, userID, couponCode string, usedAt time.Time) error {
	_, err := tx.ExecContext(ctx, `
        INSERT INTO coupon_usages (user_id, coupon_code, used_at)
        VALUES ($1, $2, $3)
    `, userID, couponCode, usedAt)
	return err
}

func (r *CouponRepository) GetValidCoupons(ctx context.Context, currentTime time.Time) ([]models.Coupon, error) {
	rows, err := r.DBHelper.PostgresClient.QueryContext(ctx, `
        SELECT 
            coupon_code, expiry_date, usage_type, 
            applicable_medicine_ids, applicable_categories, 
            min_order_value, valid_start, valid_end, 
            terms_and_conditions, discount_type, discount_value, 
            max_usage_per_user, discount_target, max_discount_amount
        FROM coupons
        WHERE expiry_date >= $1
    `, currentTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coupons []models.Coupon

	for rows.Next() {
		var c models.Coupon
		var meds, cats []byte

		err := rows.Scan(
			&c.CouponCode, &c.ExpiryDate, &c.UsageType,
			&meds, &cats,
			&c.MinOrderValue, &c.ValidTimeWindow.Start, &c.ValidTimeWindow.End,
			&c.TermsAndConditions, &c.DiscountType, &c.DiscountValue,
			&c.MaxUsagePerUser, &c.DiscountTarget, &c.MaxDiscountAmount,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(meds, &c.ApplicableMedicineIDs); err != nil {
			return nil, fmt.Errorf("failed to unmarshal medicine IDs: %w", err)
		}
		if err := json.Unmarshal(cats, &c.ApplicableCategories); err != nil {
			return nil, fmt.Errorf("failed to unmarshal categories: %w", err)
		}

		coupons = append(coupons, c)
	}

	return coupons, nil
}
