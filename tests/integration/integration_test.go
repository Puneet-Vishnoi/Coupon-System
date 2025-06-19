package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Puneet-Vishnoi/Coupon-System/models"
	"github.com/Puneet-Vishnoi/Coupon-System/routes"
	"github.com/Puneet-Vishnoi/Coupon-System/tests/mockdb"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert"
)

func TestCouponIntegration(t *testing.T) {
	// Setup test dependencies
	testDeps := mockdb.GetTestInstance()
	defer testDeps.Cleanup()

	// Start the test server
	router := gin.Default()
	routes.RegisterRoutes(router, testDeps.Service)
	server := httptest.NewServer(router)
	defer server.Close()

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

	couponJSON, _ := json.Marshal(coupon)
	resp, err := http.Post(server.URL+"/api/coupons", "application/json", bytes.NewBuffer(couponJSON))
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Failed to create coupon: %v", err)
	}
	defer resp.Body.Close()

	// 2. Check applicable
	cartItems := []models.CartItem{
		{ID: "med001", Category: "pain_relief", Price: 100.78},
		{ID: "med004", Category: "fever", Price: 342.89},
	}
	applicableReq := models.ApplicableCouponsRequest{
		CartItems:  cartItems,
		OrderTotal: 120.5,
		Timestamp:  time.Date(2025, time.June, 1, 10, 30, 0, 0, time.UTC),
	}
	applicableJSON, _ := json.Marshal(applicableReq)
	resp, err = http.Post(server.URL+"/api/coupons/applicable", "application/json", bytes.NewBuffer(applicableJSON))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to check applicable coupons: %v", err)
	}
	defer resp.Body.Close()

	// 3. Validate the coupon
	validateReq := models.ValidateCouponRequest{
		UserID:     "Puneet001",
		CouponCode: "SAVE20",
		CartItems:  cartItems,
		OrderTotal: 120.5,
		Timestamp:  time.Date(2025, time.June, 1, 10, 30, 0, 0, time.UTC),
	}
	validateJSON, _ := json.Marshal(validateReq)
	resp, err = http.Post(server.URL+"/api/coupons/validate", "application/json", bytes.NewBuffer(validateJSON))
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Failed to validate coupon: %v", err)
	}
	defer resp.Body.Close()

	var result models.ValidateCouponResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	assert.Equal(t, result.IsValid, true)
	assert.Equal(t, result.Discount["total_order_value"], 24.1)
}
