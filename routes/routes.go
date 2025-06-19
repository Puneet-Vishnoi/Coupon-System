package routes

import (
	"net/http"

	"github.com/Puneet-Vishnoi/Coupon-System/handlers"
	"github.com/Puneet-Vishnoi/Coupon-System/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, service *service.CouponService) {
	couponHandler := handlers.NewCouponHandler(service)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	api := router.Group("/api")
	{
		api.POST("/coupons", couponHandler.CreateCoupon)
		api.POST("/coupons/applicable", couponHandler.GetApplicableCoupons)
		api.POST("/coupons/validate", couponHandler.ValidateCoupon)
	}
}
