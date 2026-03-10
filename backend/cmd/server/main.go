package main

import (
	"net/http"

	"coupon-marketplace/internal/database"
	"coupon-marketplace/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	couponHandler := handlers.NewCouponHandler()
	router.POST("/admin/coupons", couponHandler.CreateCoupon)

	router.Run(":8080")
}
