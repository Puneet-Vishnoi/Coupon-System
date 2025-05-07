package unittest

import (
	"context"
	"testing"
	"time"

	"github.com/Puneet-Vishnoi/Coupon-System/models"
	"github.com/Puneet-Vishnoi/Coupon-System/tests/mockdb"
	"github.com/go-playground/assert"
)

var test = mockdb.GetTestInstance()

func TestValidateCoupon(t *testing.T) {
	now := time.Date(2025, 5, 7, 12, 0, 0, 0, time.UTC)

	baseCoupon := models.Coupon{
		DiscountType:          "flat",
		DiscountValue:         50,
		DiscountTarget:        "medicine",
		MinOrderValue:         100,
		MaxUsagePerUser:       1,
		UsageType:             "single_use",
		TermsAndConditions:    "Valid on medicine only",
		ApplicableCategories:  []string{"painkillers"},
		ApplicableMedicineIDs: []string{"med001", "med002"},
	}

	tests := []struct {
		name    string
		setup   func() string
		request models.ValidateCouponRequest
		wantErr string
	}{
		{
			name: "Expired Coupon",
			setup: func() string {
				c := baseCoupon
				c.CouponCode = "EXPIRED"
				c.ExpiryDate = now.Add(-24 * time.Hour)
				test.Service.CreateCoupon(context.Background(), &c)
				// log.Println(c)
				return c.CouponCode
			},
			request: models.ValidateCouponRequest{
				UserID:     "user1",
				OrderTotal: 200,
				Timestamp:  time.Now(),
			},
			wantErr: "coupon expired",
		},
		{
			name: "Coupon Not Valid in Time Window",
			setup: func() string {
				c := baseCoupon
				c.CouponCode = "TIMEFAIL"
				c.ValidTimeWindow = models.TimeWindow{
					Start: now.Add(1 * time.Hour),
					End:   now.Add(2 * time.Hour),
				}
				c.ExpiryDate = now.Add(24 * time.Hour)
				test.Service.CreateCoupon(context.Background(), &c)
				return c.CouponCode
			},
			request: models.ValidateCouponRequest{
				UserID:     "user2",
				OrderTotal: 200,
			},
			wantErr: "coupon not valid at this time",
		},
		{
			name: "Minimum Order Not Met",
			setup: func() string {
				c := baseCoupon
				c.CouponCode = "MINFAIL"
				c.MinOrderValue = 500
				c.ExpiryDate = now.Add(24 * time.Hour)
				test.Service.CreateCoupon(context.Background(), &c)
				//log.Println(c)

				return c.CouponCode
			},
			request: models.ValidateCouponRequest{
				UserID:     "user3",
				OrderTotal: 100,
				CartItems:  []models.CartItem{{Category: "painkillers"}},
			},
			wantErr: "order total does not meet minimum requirement",
		},
		{
			name: "Coupon Not Found",
			setup: func() string {
				return "INVALIDCODE"
			},
			request: models.ValidateCouponRequest{
				UserID:     "user4",
				OrderTotal: 200,
			},
			wantErr: "coupon not found",
		},
		{
			name: "Valid Coupon Use Case",
			setup: func() string {
				c := baseCoupon
				c.CouponCode = "VALID1"
				c.ExpiryDate = now.Add(24 * time.Hour)
				c.ValidTimeWindow = models.TimeWindow{
					Start: now.Add(-1 * time.Hour),
					End:   now.Add(2 * time.Hour),
				}
				test.Service.CreateCoupon(context.Background(), &c)
				return c.CouponCode
			},
			request: models.ValidateCouponRequest{
				UserID:     "user5",
				OrderTotal: 200,
				Timestamp:  now,
				CartItems: []models.CartItem{
					{ID: "med001", Category: "painkillers"},
				},
			},
			wantErr: "",
		},
	}
	t.Cleanup(func() { test.Cleanup() })
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code := tc.setup()
			tc.request.CouponCode = code
			resp, err := test.Service.ValidateCoupon(context.Background(), tc.request)
			if tc.wantErr != "" {
				assert.Equal(t, tc.wantErr, err.Error())
			} else {
				assert.Equal(t, nil, err)
				assert.NotEqual(t, 0.0, resp.Discount)
			}
		})
	}

}

func TestCreateCoupon(t *testing.T) {
	t.Cleanup(func() { test.Cleanup() })
	now := time.Now()

	// Valid Coupon
	coupon := &models.Coupon{
		CouponCode:            "CREATE1",
		ExpiryDate:            now.Add(48 * time.Hour),
		UsageType:             "single_use",
		ApplicableMedicineIDs: []string{"med101"},
		MinOrderValue:         150,
		ValidTimeWindow: models.TimeWindow{
			Start: now.Add(-1 * time.Hour),
			End:   now.Add(2 * time.Hour),
		},
		DiscountType:       "flat",
		DiscountValue:      30,
		MaxUsagePerUser:    1,
		DiscountTarget:     "medicine",
		TermsAndConditions: "Valid only on med101",
	}
	err := test.Service.CreateCoupon(context.Background(), coupon)
	assert.Equal(t, nil, err)

	// Invalid Coupon (negative discount)
	invalid := &models.Coupon{
		CouponCode:      "INVALID1",
		DiscountType:    "flat",
		DiscountValue:   -10,
		DiscountTarget:  "medicine",
		MinOrderValue:   100,
		ExpiryDate:      now.Add(1 * time.Hour),
		MaxUsagePerUser: 1,
		UsageType:       "single_use",
		ValidTimeWindow: models.TimeWindow{Start: now.Add(-1 * time.Hour), End: now.Add(2 * time.Hour)},
	}
	err = test.Service.CreateCoupon(context.Background(), invalid)
	assert.NotEqual(t, nil, err)
}

func TestGetApplicableCoupons(t *testing.T) {
	t.Cleanup(func() { test.Cleanup() })
	now := time.Now()

	coupon1 := &models.Coupon{
		CouponCode:           "APPLICABLE1",
		ExpiryDate:           now.Add(24 * time.Hour),
		UsageType:            "multi_use",
		MinOrderValue:        100,
		ApplicableCategories: []string{"diabetes"},
		ValidTimeWindow: models.TimeWindow{
			Start: now.Add(-1 * time.Hour),
			End:   now.Add(2 * time.Hour),
		},
		DiscountType:       "flat",
		DiscountValue:      20,
		DiscountTarget:     "medicine",
		MaxUsagePerUser:    5,
		TermsAndConditions: "Applicable on diabetes category",
	}
	_ = test.Service.CreateCoupon(context.Background(), coupon1)

	coupon2 := &models.Coupon{
		CouponCode:           "NOTMATCHING",
		ExpiryDate:           now.Add(24 * time.Hour),
		UsageType:            "multi_use",
		MinOrderValue:        100,
		ApplicableCategories: []string{"painkillers"},
		ValidTimeWindow: models.TimeWindow{
			Start: now.Add(-1 * time.Hour),
			End:   now.Add(2 * time.Hour),
		},
		DiscountType:       "percentage",
		DiscountValue:      10,
		DiscountTarget:     "medicine",
		MaxUsagePerUser:    5,
		TermsAndConditions: "Applicable on painkillers",
	}
	_ = test.Service.CreateCoupon(context.Background(), coupon2)

	req := models.ApplicableCouponsRequest{
		OrderTotal: 150,
		Timestamp:  now,
		CartItems: []models.CartItem{
			{ID: "med301", Category: "diabetes", Price: 150},
		},
	}

	coupons, err := test.Service.GetApplicableCoupons(context.Background(), req)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(coupons))
	assert.Equal(t, "APPLICABLE1", coupons[0].CouponCode)
}
