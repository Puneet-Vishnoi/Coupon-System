package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	redisProvider "github.com/Puneet-Vishnoi/Coupon-System/cache/redis/providers"
	"github.com/Puneet-Vishnoi/Coupon-System/models"
	"github.com/Puneet-Vishnoi/Coupon-System/repository"
)

type CouponService struct {
	Repo        *repository.CouponRepository
	RedisHelper *redisProvider.RedisHelper
}

func NewCouponService(repo *repository.CouponRepository, redis *redisProvider.RedisHelper) *CouponService {
	return &CouponService{Repo: repo, RedisHelper: redis}
}

func (s *CouponService) CreateCoupon(ctx context.Context, coupon *models.Coupon) (err error) {
	tx, err := s.Repo.DBHelper.PostgresClient.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	c, err := s.Repo.GetCouponByCode(ctx, tx, coupon.CouponCode)
	if err == nil && c.CouponCode == coupon.CouponCode {
		return errors.New("coupon already exists")
	}
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	err = s.Repo.CreateCoupon(ctx, tx, coupon)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	// Invalidate any cached data related to coupon list
	s.RedisHelper.Delete(ctx, "valid_coupons")
	return nil
}

func (s *CouponService) GetApplicableCoupons(ctx context.Context, req models.ApplicableCouponsRequest) ([]models.Coupon, error) {
	var allCoupons []models.Coupon
	cacheHit, err := s.RedisHelper.GetJSON(ctx, "valid_coupons", &allCoupons)
	if err != nil {
		return nil, err
	}
	if !cacheHit {
		allCoupons, err = s.Repo.GetValidCoupons(ctx, req.Timestamp)
		if err != nil {
			return nil, err
		}
		s.RedisHelper.SetJSON(ctx, "valid_coupons", allCoupons, 10*time.Minute)
	}

	var applicable []models.Coupon
	for _, c := range allCoupons {
		if req.OrderTotal < c.MinOrderValue {
			continue
		}

		if len(c.ApplicableMedicineIDs) == 0 && len(c.ApplicableCategories) == 0 {
			applicable = append(applicable, c)
			continue
		}

		for _, item := range req.CartItems {
			if (len(c.ApplicableMedicineIDs) > 0 && contains(c.ApplicableMedicineIDs, item.ID)) ||
				(len(c.ApplicableCategories) > 0 && contains(c.ApplicableCategories, item.Category)) {
				applicable = append(applicable, c)
				break
			}
		}
	}

	return applicable, nil
}

func (s *CouponService) ValidateCoupon(ctx context.Context, req models.ValidateCouponRequest) (resp models.ValidateCouponResponse, err error) {
	tx, err := s.Repo.DBHelper.PostgresClient.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return resp, errors.New("failed to start transaction")
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	coupon, err := s.fetchCouponFromDB(ctx, tx, req.CouponCode)
	if err != nil {
		return resp, err
	}

	if req.Timestamp.After(coupon.ExpiryDate) {
		return resp, errors.New("coupon expired")
	}

	if req.Timestamp.Before(coupon.ValidTimeWindow.Start) || req.Timestamp.After(coupon.ValidTimeWindow.End) {
		return resp, errors.New("coupon not valid at this time")
	}

	usageCount, err := s.Repo.GetUserUsageCount(ctx, tx, req.UserID, req.CouponCode)
	if err != nil {
		return resp, err
	}
	if usageCount >= coupon.MaxUsagePerUser {
		return resp, errors.New("usage limit reached")
	}

	applicable := false
	for _, item := range req.CartItems {
		if contains(coupon.ApplicableMedicineIDs, item.ID) || contains(coupon.ApplicableCategories, item.Category) {
			applicable = true
			break
		}
	}
	if !applicable {
		return resp, errors.New("coupon not applicable to cart items")
	}

	if req.OrderTotal < coupon.MinOrderValue {
		return resp, errors.New("order total does not meet minimum requirement")
	}

	discount := calculateDiscount(coupon, req.OrderTotal)

	err = s.Repo.RecordUsage(ctx, tx, req.UserID, coupon.CouponCode, req.Timestamp)
	if err != nil {
		return resp, err
	}

	if err := tx.Commit(); err != nil {
		return resp, errors.New("failed to commit transaction")
	}

	// Invalidate coupon cache as the usage may impact validity
	s.RedisHelper.Delete(ctx, "valid_coupons")

	resp = models.ValidateCouponResponse{
		IsValid:  true,
		Discount: discount,
		Message:  "coupon applied successfully",
	}
	return resp, nil
}

func (s *CouponService) fetchCouponFromDB(ctx context.Context, tx *sql.Tx, couponCode string) (models.Coupon, error) {
	coupon, err := s.Repo.GetCouponByCode(ctx, tx, couponCode)
	if err != nil {
		return models.Coupon{}, errors.New("coupon not found")
	}
	return coupon, nil
}

// Helpers
func contains(slice []string, target string) bool {
	for _, item := range slice {
		if item == target {
			return true
		}
	}
	return false
}

func calculateDiscount(coupon models.Coupon, orderTotal float64) map[string]float64 {
	discount := make(map[string]float64)

	switch coupon.DiscountType {
	case models.DiscountTypeFlat:
		discount[string(coupon.DiscountTarget)] = min(coupon.DiscountValue, coupon.MaxDiscountAmount)
	case models.DiscountTypePercentage:
		calculated := (orderTotal * coupon.DiscountValue) / 100
		discount[string(coupon.DiscountTarget)] = min(calculated, coupon.MaxDiscountAmount)
	}
	return discount
}

func min(a, b float64) float64 {
	if b <= 0 {
		return a
	}
	if a < b {
		return a
	}
	return b
}
