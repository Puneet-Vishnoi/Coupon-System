package routes

import (
	"github.com/Puneet-Vishnoi/Coupon-System/handlers"
	"github.com/Puneet-Vishnoi/Coupon-System/service"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, service *service.CouponService) {
	couponHandler := handlers.NewCouponHandler(service)

	api := router.Group("/api")
	{
		api.POST("/coupons", couponHandler.CreateCoupon)
		api.GET("/coupons/applicable", couponHandler.GetApplicableCoupons)
		api.POST("/coupons/validate", couponHandler.ValidateCoupon)
	}
}
