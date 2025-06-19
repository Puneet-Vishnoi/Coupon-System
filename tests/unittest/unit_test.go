package unittest

import (
	"context"
	"testing"
	"time"

	"github.com/Puneet-Vishnoi/Coupon-System/models"
	"github.com/Puneet-Vishnoi/Coupon-System/tests/mockdb"
	"github.com/go-playground/assert"
)

func setupTest(t *testing.T) *mockdb.TestDeps {
	test := mockdb.GetTestInstance()
	t.Cleanup(func() {
		if err := test.PostgresClient.ClearTestData(); err != nil {
			t.Fatalf("failed to clear test data: %v", err)
		}
		test.Cleanup()
	})
	return test
}

func TestValidateCoupon(t *testing.T) {
	now := time.Date(2025, 5, 7, 12, 0, 0, 0, time.UTC)
	baseCoupon := models.Coupon{
		DiscountType:          "flat",
		DiscountValue:         50,
		DiscountTarget:        "total_order_value",
		MinOrderValue:         100,
		MaxUsagePerUser:       1,
		UsageType:             "single_use",
		TermsAndConditions:    "Valid on medicine only",
		ApplicableCategories:  []string{"painkillers"},
		ApplicableMedicineIDs: []string{"med001", "med002"},
	}

	tests := []struct {
		name    string
		setup   func(t *testing.T, test *mockdb.TestDeps) string
		request models.ValidateCouponRequest
		wantErr string
	}{
		{
			name: "Expired Coupon",
			setup: func(t *testing.T, test *mockdb.TestDeps) string {
				c := baseCoupon
				c.CouponCode = "EXPIRED"
				c.ExpiryDate = now.Add(-24 * time.Hour)
				if err := test.Service.CreateCoupon(context.Background(), &c); err != nil {
					t.Fatalf("failed to insert coupon: %v", err)
				}
				return c.CouponCode
			},
			request: models.ValidateCouponRequest{
				UserID:     "user1",
				OrderTotal: 200,
				Timestamp:  now,
			},
			wantErr: "coupon expired",
		},
		{
			name: "Coupon Not Valid in Time Window",
			setup: func(t *testing.T, test *mockdb.TestDeps) string {
				c := baseCoupon
				c.CouponCode = "TIMEFAIL"
				c.ValidTimeWindow = models.TimeWindow{
					Start: now.Add(1 * time.Hour),
					End:   now.Add(2 * time.Hour),
				}
				c.ExpiryDate = now.Add(24 * time.Hour)
				if err := test.Service.CreateCoupon(context.Background(), &c); err != nil {
					t.Fatalf("failed to insert coupon: %v", err)
				}
				return c.CouponCode
			},
			request: models.ValidateCouponRequest{
				UserID:     "user2",
				OrderTotal: 200,
				Timestamp:  now,
			},
			wantErr: "coupon not valid at this time",
		},
		{
			name: "Minimum Order Not Met",
			setup: func(t *testing.T, test *mockdb.TestDeps) string {
				c := baseCoupon
				c.CouponCode = "MINFAIL"
				c.MinOrderValue = 500
				c.ExpiryDate = now.Add(24 * time.Hour)
				c.ValidTimeWindow = models.TimeWindow{
					Start: now.Add(-24 * time.Hour),
					End:   now.Add(24 * time.Hour),
				}
				if err := test.Service.CreateCoupon(context.Background(), &c); err != nil {
					t.Fatalf("failed to insert coupon: %v", err)
				}
				return c.CouponCode
			},
			request: models.ValidateCouponRequest{
				UserID:     "user3",
				OrderTotal: 100,
				CartItems:  []models.CartItem{{Category: "painkillers"}},
				Timestamp:  now,
			},
			wantErr: "order total does not meet minimum requirement",
		},
		{
			name: "Coupon Not Found",
			setup: func(t *testing.T, test *mockdb.TestDeps) string {
				return "INVALIDCODE"
			},
			request: models.ValidateCouponRequest{
				UserID:     "user4",
				OrderTotal: 200,
				Timestamp:  now,
			},
			wantErr: "coupon not found",
		},
		{
			name: "Valid Coupon Use Case",
			setup: func(t *testing.T, test *mockdb.TestDeps) string {
				c := baseCoupon
				c.CouponCode = "VALID1"
				c.ExpiryDate = now.Add(24 * time.Hour)
				c.ValidTimeWindow = models.TimeWindow{
					Start: now.Add(-1 * time.Hour),
					End:   now.Add(2 * time.Hour),
				}
				if err := test.Service.CreateCoupon(context.Background(), &c); err != nil {
					t.Fatalf("failed to insert coupon: %v", err)
				}
				return c.CouponCode
			},
			request: models.ValidateCouponRequest{
				UserID:     "user5",
				OrderTotal: 200,
				Timestamp:  now,
				CartItems:  []models.CartItem{{ID: "med001", Category: "painkillers"}},
			},
			wantErr: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			test := setupTest(t)
			code := tc.setup(t, test)
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
	test := setupTest(t)
	now := time.Now()

	valid := &models.Coupon{
		CouponCode:            "CREATE1",
		ExpiryDate:            now.Add(48 * time.Hour),
		UsageType:             "single_use",
		ApplicableMedicineIDs: []string{"med101"},
		MinOrderValue:         150,
		ValidTimeWindow:       models.TimeWindow{Start: now.Add(-1 * time.Hour), End: now.Add(2 * time.Hour)},
		DiscountType:          "flat",
		DiscountValue:         30,
		MaxUsagePerUser:       1,
		DiscountTarget:        "total_order_value",
		TermsAndConditions:    "Valid only on med101",
	}
	assert.Equal(t, nil, test.Service.CreateCoupon(context.Background(), valid))

	invalid := &models.Coupon{
		CouponCode:      "INVALID1",
		DiscountType:    "flat",
		DiscountValue:   -10,
		DiscountTarget:  "total_order_value",
		MinOrderValue:   100,
		ExpiryDate:      now.Add(1 * time.Hour),
		MaxUsagePerUser: 1,
		UsageType:       "single_use",
		ValidTimeWindow: models.TimeWindow{Start: now.Add(-1 * time.Hour), End: now.Add(2 * time.Hour)},
	}
	assert.NotEqual(t, nil, test.Service.CreateCoupon(context.Background(), invalid))
}

func TestGetApplicableCoupons(t *testing.T) {
	test := setupTest(t)
	now := time.Now()

	coupon1 := &models.Coupon{
		CouponCode:           "APPLICABLE1",
		ExpiryDate:           now.Add(24 * time.Hour),
		UsageType:            "multi_use",
		MinOrderValue:        100,
		ApplicableCategories: []string{"diabetes"},
		ValidTimeWindow:      models.TimeWindow{Start: now.Add(-1 * time.Hour), End: now.Add(2 * time.Hour)},
		DiscountType:         "flat",
		DiscountValue:        20,
		DiscountTarget:       "total_order_value",
		MaxUsagePerUser:      5,
		TermsAndConditions:   "Applicable on diabetes category",
	}
	assert.Equal(t, nil, test.Service.CreateCoupon(context.Background(), coupon1))

	coupon2 := &models.Coupon{
		CouponCode:           "NOTMATCHING",
		ExpiryDate:           now.Add(24 * time.Hour),
		UsageType:            "multi_use",
		MinOrderValue:        100,
		ApplicableCategories: []string{"painkillers"},
		ValidTimeWindow:      models.TimeWindow{Start: now.Add(-1 * time.Hour), End: now.Add(2 * time.Hour)},
		DiscountType:         "percentage",
		DiscountValue:        10,
		DiscountTarget:       "total_order_value",
		MaxUsagePerUser:      5,
		TermsAndConditions:   "Applicable on painkillers",
	}
	assert.Equal(t, nil, test.Service.CreateCoupon(context.Background(), coupon2))

	req := models.ApplicableCouponsRequest{
		OrderTotal: 150,
		Timestamp:  now,
		CartItems:  []models.CartItem{{ID: "med301", Category: "diabetes", Price: 150}},
	}

	coupons, err := test.Service.GetApplicableCoupons(context.Background(), req)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(coupons))
	assert.Equal(t, "APPLICABLE1", coupons[0].CouponCode)
}
