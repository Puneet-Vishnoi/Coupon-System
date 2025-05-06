package models

type ValidateCouponResponse struct {
	IsValid  bool               `json:"is_valid"`
	Discount map[string]float64 `json:"discount"` // e.g., {"medicine": 25.0}
	Message  string             `json:"message"`
}
