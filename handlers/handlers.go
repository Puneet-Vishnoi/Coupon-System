package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/Puneet-Vishnoi/Coupon-System/models"
	"github.com/Puneet-Vishnoi/Coupon-System/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type CouponHandler struct {
	Service   *service.CouponService
	Validator *validator.Validate
}

func NewCouponHandler(s *service.CouponService) *CouponHandler {
	return &CouponHandler{
		Service:   s,
		Validator: validator.New(),
	}
}

// formatValidationError formats validation errors from validator package into a map
func formatValidationError(err error) map[string]string {
	errors := make(map[string]string)
	for _, e := range err.(validator.ValidationErrors) {
		errors[e.Field()] = "failed on tag '" + e.Tag() + "'"
	}
	return errors
}

// POST /coupons
func (h *CouponHandler) CreateCoupon(c *gin.Context) {
	var req models.Coupon
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.Validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationError(err)})
		return
	}

	if err := h.Service.CreateCoupon(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Coupon created successfully"})
}

// POST /coupons/applicable
func (h *CouponHandler) GetApplicableCoupons(c *gin.Context) {
	var req models.ApplicableCouponsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	if err := h.Validator.Struct(req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationError(err)})
		return
	}

	coupons, err := h.Service.GetApplicableCoupons(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, coupons)
}

// POST /coupons/validate
func (h *CouponHandler) ValidateCoupon(c *gin.Context) {
	var req models.ValidateCouponRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if req.Timestamp.IsZero() {
		req.Timestamp = time.Now()
	}

	if err := h.Validator.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"validation_errors": formatValidationError(err)})
		return
	}

	resp, err := h.Service.ValidateCoupon(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}
