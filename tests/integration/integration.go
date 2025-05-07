package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/Puneet-Vishnoi/Coupon-System/models"
	"github.com/Puneet-Vishnoi/Coupon-System/tests/mockdb"
	"github.com/go-playground/assert/v2"
)

func TestCouponIntegration(t *testing.T) {
	// Initialize the test dependencies
	testDeps := mockdb.GetTestInstance()
	defer testDeps.Cleanup()

	// 1. Create a coupon
	coupon := models.Coupon{
		CouponCode:            "SAVE20",
		ExpiryDate:            time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
		UsageType:             "single_use",
		ApplicableMedicineIDs: []string{"med001", "med002", "med003"},
		ApplicableCategories:  []string{"pain_relief", "fever"},
		MinOrderValue:         100.5,
		ValidTimeWindow: models.TimeWindow{
			Start: time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
			End:   time.Date(2025, time.December, 31, 23, 59, 59, 0, time.UTC),
		},
		TermsAndConditions: "This coupon is valid for selected medicines only.",
		DiscountType:       "percentage",
		DiscountValue:      20,
		MaxUsagePerUser:    1,
		DiscountTarget:     "total_order_value",
	}

	// Send POST request to create the coupon
	couponJSON, err := json.Marshal(coupon)
	if err != nil {
		t.Fatalf("failed to marshal coupon: %v", err)
	}

	resp, err := http.Post("http://localhost:8080/api/coupons", "application/json", bytes.NewBuffer(couponJSON))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to create coupon: %v", err)
	}
	defer resp.Body.Close()

	// 2. Check if the coupon is applicable
	cartItems := []models.CartItem{
		{ID: "med001", Category: "pain_relief", Price: 100.78},
		{ID: "med004", Category: "fever", Price: 342.89},
	}

	applicableRequest := models.ApplicableCouponsRequest{
		CartItems:  cartItems,
		OrderTotal: 120.5,
		Timestamp:  time.Date(2025, time.June, 1, 10, 30, 0, 0, time.UTC),
	}

	applicableJSON, err := json.Marshal(applicableRequest)
	if err != nil {
		t.Fatalf("failed to marshal applicable request: %v", err)
	}

	resp, err = http.Post("http://localhost:8080/api/coupons/applicable", "application/json", bytes.NewBuffer(applicableJSON))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to check if coupon is applicable: %v", err)
	}
	defer resp.Body.Close()

	// 3. Validate the coupon for a specific user and cart
	validateRequest := models.ValidateCouponRequest{
		UserID:     "Puneet001",
		CouponCode: "SAVE20",
		CartItems:  cartItems,
		OrderTotal: 120.5,
		Timestamp:  time.Date(2025, time.June, 1, 10, 30, 0, 0, time.UTC),
	}

	validateJSON, err := json.Marshal(validateRequest)
	if err != nil {
		t.Fatalf("failed to marshal validate request: %v", err)
	}

	resp, err = http.Post("http://localhost:8080/api/coupons/validate", "application/json", bytes.NewBuffer(validateJSON))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("failed to validate coupon: %v", err)
	}
	defer resp.Body.Close()

	// 4. Verify the coupon is applied correctly (expected to be applied)
	var result models.ValidateCouponResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	assert.Equal(t, result.IsValid, true)
	assert.Equal(t, result.Discount, 24.1) 
}
