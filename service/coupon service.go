package service

import (
	"context"
	"errors"
	"log"

	redisProvider "github.com/Puneet-Vishnoi/Coupon-System/cache/redis/providers"
	"github.com/Puneet-Vishnoi/Coupon-System/models"
	"github.com/Puneet-Vishnoi/Coupon-System/repository"
)

type CouponService struct {
	Repo *repository.CouponRepository
}

func NewCouponService(repo *repository.CouponRepository, redis *redisProvider.RedisHelper) *CouponService {
	return &CouponService{Repo: repo}
}

func (s *CouponService) CreateCoupon(ctx context.Context, coupon *models.Coupon) error {
	log.Println(coupon, "kkkkk")
	return s.Repo.CreateCoupon(ctx, coupon) 
}

func (s *CouponService) GetApplicableCoupons(ctx context.Context, req models.ApplicableCouponsRequest) ([]models.Coupon, error) {
	allCoupons, err := s.Repo.GetValidCoupons(ctx, req.Timestamp)
	if err != nil {
		return nil, err
	}
	log.Println(allCoupons, "all")
	log.Println(req, "req")


	var applicable []models.Coupon
	for _, c := range allCoupons {
		log.Println(c.ApplicableCategories, "acata")
		log.Println(c.ApplicableMedicineIDs, "amedi")

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
	coupon, err := s.Repo.GetCouponByCode(ctx, req.CouponCode)
	if err != nil {
		return models.ValidateCouponResponse{}, errors.New("coupon not found")
	}

	if req.Timestamp.After(coupon.ExpiryDate) {
		return models.ValidateCouponResponse{}, errors.New("coupon expired")
	}

	usageCount, err := s.Repo.GetUserUsageCount(ctx, req.UserID, req.CouponCode)
	if err != nil {
		return models.ValidateCouponResponse{}, err
	}
	if usageCount >= coupon.MaxUsagePerUser {
		return models.ValidateCouponResponse{}, errors.New("usage limit reached")
	}

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

	if req.OrderTotal < coupon.MinOrderValue {
		return models.ValidateCouponResponse{}, errors.New("order total does not meet minimum requirement")
	}

	discount := calculateDiscount(coupon, req.OrderTotal)

	_ = s.Repo.RecordUsage(ctx, req.UserID, coupon.CouponCode, req.Timestamp)

	return models.ValidateCouponResponse{
		IsValid:  true,
		Discount: discount,
		Message:  "coupon applied successfully",
	}, nil
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
