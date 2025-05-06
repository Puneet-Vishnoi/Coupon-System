package controller

// import (
// 	"net/http"

// 	"github.com/Puneet-Vishnoi/Coupon-System/models"
// 	"github.com/Puneet-Vishnoi/Coupon-System/service"
// 	"github.com/gin-gonic/gin"
// )

// type CouponController struct {
// 	Service *service.CouponService
// }

// func NewCouponController(service *service.CouponService) *CouponController {
// 	return &CouponController{Service: service}
// }

// // POST /admin/coupons
// func (cc *CouponController) CreateCoupon(c *gin.Context) {
// 	var req models.Coupon
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	err := cc.Service.CreateCoupon(&req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "coupon created successfully"})
// }

// // GET /coupons/applicable
// func (cc *CouponController) GetApplicableCoupons(c *gin.Context) {
// 	var req models.ApplicableCouponsRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	coupons, err := cc.Service.GetApplicableCoupons(req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"applicable_coupons": coupons})
// }

// // POST /coupons/validate
// func (cc *CouponController) ValidateCoupon(c *gin.Context) {
// 	var req models.ValidateCouponRequest
// 	if err := c.ShouldBindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	resp, err := cc.Service.ValidateCoupon(req)
// 	if err != nil {
// 		c.JSON(http.StatusOK, gin.H{"is_valid": false, "reason": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, resp)
// }
