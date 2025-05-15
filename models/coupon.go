package models

import "time"

type TimeWindow struct {
	Start time.Time `json:"valid_start" validate:"required"`
	End   time.Time `json:"valid_end" validate:"required,gtfield=Start"`
}

type Coupon struct {
	CouponCode            string         `json:"coupon_code" validate:"required"`
	DiscountType          DiscountType   `json:"discount_type" validate:"required"`
	DiscountValue         float64        `json:"discount_value" validate:"required,gt=0"`
	DiscountTarget        DiscountTarget `json:"discount_target" validate:"required"`
	MinOrderValue         float64        `json:"min_order_value" validate:"gte=0"`
	MaxUsagePerUser       int            `json:"max_usage_per_user" validate:"gte=0"`
	ExpiryDate            time.Time      `json:"expiry_date" validate:"required"`
	ApplicableMedicineIDs []string       `json:"applicable_medicine_ids" validate:"dive"`
	ApplicableCategories  []string       `json:"applicable_categories" validate:"dive"`
	UsageType             UsageType      `json:"usage_type" validate:"required"`
	ValidTimeWindow       TimeWindow     `json:"valid_time_window" validate:"required"`
	TermsAndConditions    string         `json:"terms_and_conditions"`
	MaxDiscountAmount     float64        `json:"max_discount_amount" validate:"gte=0"`
}