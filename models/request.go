package models

import "time"

type ApplicableCouponsRequest struct {
	CartItems  []CartItem `json:"cart_items" validate:"required,dive"`
	OrderTotal float64    `json:"order_total" validate:"required,gt=0"`
	Timestamp  time.Time  `json:"timestamp" validate:"required"`
}

type ValidateCouponRequest struct {
	UserID     string     `json:"user_id" validate:"required"`
	CouponCode string     `json:"coupon_code" validate:"required"`
	CartItems  []CartItem `json:"cart_items" validate:"required,dive"`
	OrderTotal float64    `json:"order_total" validate:"required,gt=0"`
	Timestamp  time.Time  `json:"timestamp" validate:"required"`
}

type CartItem struct {
	ID       string  `json:"medicine_id" validate:"required"`
	Category string  `json:"category" validate:"required"`
	Price    float64 `json:"price" validate:"required,gt=0"`
}
