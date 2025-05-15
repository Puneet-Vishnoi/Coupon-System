package models

import "time"

type TimeWindow struct {
	Start time.Time `json:"valid_start"`
	End   time.Time `json:"valid_end"`
}

type Coupon struct {
	CouponCode            string         `json:"coupon_code"`
	DiscountType          DiscountType   `json:"discount_type"`
	DiscountValue         float64        `json:"discount_value"`
	DiscountTarget        DiscountTarget `json:"discount_target"`
	MinOrderValue         float64        `json:"min_order_value"`
	MaxUsagePerUser       int            `json:"max_usage_per_user"`
	ExpiryDate            time.Time      `json:"expiry_date"`
	ApplicableMedicineIDs []string       `json:"applicable_medicine_ids"`
	ApplicableCategories  []string       `json:"applicable_categories"`
	UsageType             UsageType      `json:"usage_type"`
	ValidTimeWindow       TimeWindow     `json:"valid_time_window"`
	TermsAndConditions    string         `json:"terms_and_conditions"`
	MaxDiscountAmount     float64        `json:"max_discount_amount"`
}
