package models

import "time"

type ApplicableCouponsRequest struct {
	CartItems  []CartItem `json:"cart_items"`
	OrderTotal float64    `json:"order_total"`
	Timestamp  time.Time  `json:"timestamp"`
}

type ValidateCouponRequest struct {
	UserID     string     `json:"user_id"`
	CouponCode string     `json:"coupon_code"`
	CartItems  []CartItem `json:"cart_items"`
	OrderTotal float64    `json:"order_total"`
	Timestamp  time.Time  `json:"timestamp"`
}

type CartItem struct {
	ID       string  `json:"medicine_id"`       // Medicine ID
	Category string  `json:"category"` // e.g., "painkillers", "diabetes"
	Price    float64 `json:"price"`
}
