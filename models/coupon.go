package models

import "time"

type Coupon struct {
	CouponCode            string    `json:"coupon_code"`
	ExpiryDate            time.Time `json:"expiry_date"`
	UsageType             string    `json:"usage_type"`               // e.g., "single-use", "multi-use"
	ApplicableMedicineIDs []string  `json:"applicable_medicine_ids"` // Stored as JSONB
	ApplicableCategories  []string  `json:"applicable_categories"`   // Stored as JSONB
	MinOrderValue         float64   `json:"min_order_value"`
	ValidTimeWindow       TimeWindow `json:"valid_time_window"`
	TermsAndConditions    string    `json:"terms_and_conditions"`
	DiscountType          string    `json:"discount_type"`     // "flat" or "percentage"
	DiscountValue         float64   `json:"discount_value"`    // flat amount or percentage
	MaxUsagePerUser       int       `json:"max_usage_per_user"`
	DiscountTarget        string    `json:"discount_target"`   // "medicine" or "delivery"
	MaxDiscountAmount     float64   `json:"max_discount_amount"` // ✅ added missing field
}

type TimeWindow struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}
