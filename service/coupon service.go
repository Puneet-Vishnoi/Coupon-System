package service

import (
	"context"
	"errors"
	"log"
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

func (s *CouponService) CreateCoupon(ctx context.Context, coupon *models.Coupon) error {
	return s.Repo.CreateCoupon(ctx, coupon)
}

func (s *CouponService) GetApplicableCoupons(ctx context.Context, req models.ApplicableCouponsRequest) ([]models.Coupon, error) {
	allCoupons, err := s.Repo.GetValidCoupons(ctx, req.Timestamp)
	if err != nil {
		return nil, err
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
			log.Println(item)

			if contains(c.ApplicableMedicineIDs, item.ID) || contains(c.ApplicableCategories, item.Category) {
				applicable = append(applicable, c)
				break
			}
		}
	}

	return applicable, nil
}

func (s *CouponService) ValidateCoupon(ctx context.Context, req models.ValidateCouponRequest) (models.ValidateCouponResponse, error) {
	// Start a new transaction
	tx, err := s.Repo.DBHelper.PostgresClient.BeginTx(ctx, nil)
	if err != nil {
		return models.ValidateCouponResponse{}, errors.New("failed to start transaction")
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()
	// Fetch the coupon from cache or DB
	coupon, err := s.fetchCouponFromCacheOrDB(ctx, req.CouponCode)
	if err != nil {
		return models.ValidateCouponResponse{}, err
	}

	log.Println(coupon, "cpn", 	req.Timestamp.Before(coupon.ValidTimeWindow.Start) || req.Timestamp.After(coupon.ValidTimeWindow.End))

	// Check if the coupon is expired
	if req.Timestamp.After(coupon.ExpiryDate) {
		return models.ValidateCouponResponse{}, errors.New("coupon expired")
	}

	// Validate the coupon's time window
	if req.Timestamp.Before(coupon.ValidTimeWindow.Start) || req.Timestamp.After(coupon.ValidTimeWindow.End) {
		return models.ValidateCouponResponse{}, errors.New("coupon not valid at this time")
	}

	// Check if the user has exceeded the usage limit
	usageCount, err := s.Repo.GetUserUsageCount(ctx, tx, req.UserID, req.CouponCode)
	if err != nil {
		return models.ValidateCouponResponse{}, err
	}
	if usageCount >= coupon.MaxUsagePerUser {
		return models.ValidateCouponResponse{}, errors.New("usage limit reached")
	}

	log.Println(req.OrderTotal , coupon.MinOrderValue, req.OrderTotal < coupon.MinOrderValue, coupon, "//xzmmxclc")


	// Validate if the coupon is applicable for the cart items (medicine/category check)
	applicable := false
	for _, item := range req.CartItems {
		if contains(coupon.ApplicableMedicineIDs, item.ID) || contains(coupon.ApplicableCategories, item.Category) {
			applicable = true
			break
		}
	}
	if !applicable {
		return models.ValidateCouponResponse{}, errors.New("coupon not applicable to cart items")
	}
	log.Println(req.OrderTotal , coupon.MinOrderValue, req.OrderTotal < coupon.MinOrderValue)
	// Validate that the order total meets the minimum order value
	if req.OrderTotal < coupon.MinOrderValue {
		return models.ValidateCouponResponse{}, errors.New("order total does not meet minimum requirement")
	}

	// Calculate the discount
	discount := calculateDiscount(coupon, req.OrderTotal)

	// Record the coupon usage for the user
	err = s.Repo.RecordUsage(ctx, tx, req.UserID, coupon.CouponCode, req.Timestamp)
	if err != nil {
		return models.ValidateCouponResponse{}, err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return models.ValidateCouponResponse{}, errors.New("failed to commit transaction")
	}

	// Return the valid response with the calculated discount
	return models.ValidateCouponResponse{
		IsValid:  true,
		Discount: discount,
		Message:  "coupon applied successfully",
	}, nil
}


func (s *CouponService) fetchCouponFromCacheOrDB(ctx context.Context, couponCode string) (models.Coupon, error) {
	var coupon models.Coupon
	cacheHit, err := s.RedisHelper.GetJSON(ctx, couponCode, &coupon)
	if err != nil {
		return models.Coupon{}, err
	}

	if !cacheHit {
		coupon, err = s.Repo.GetCouponByCode(ctx, couponCode)
		if err != nil {
			return models.Coupon{}, errors.New("coupon not found")
		}

		// Cache it in Redis
		ttl := time.Until(coupon.ExpiryDate)
		err = s.RedisHelper.SetJSON(ctx, couponCode, coupon, ttl)
		if err != nil {
			log.Printf("failed to cache coupon: %v", err)
		}
	}
	return coupon, nil
}

// Helper
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
	if coupon.DiscountType == "flat" {
		discount[coupon.DiscountTarget] = coupon.DiscountValue
	} else if coupon.DiscountType == "percentage" {
		discountAmount := (orderTotal * coupon.DiscountValue) / 100
		discount[coupon.DiscountTarget] = discountAmount
	}
	return discount
}
