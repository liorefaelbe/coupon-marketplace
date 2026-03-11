package main

import (
	"net/http"

	"coupon-marketplace/internal/database"
	"coupon-marketplace/internal/handlers"
	"coupon-marketplace/internal/middleware"

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

	apiV1 := router.Group("/api/v1")
	apiV1.Use(middleware.ResellerAuth())
	{
		apiV1.GET("/products", couponHandler.GetAvailableProducts)
		apiV1.GET("/products/:id", couponHandler.GetProductByID)
		apiV1.POST("/products/:id/purchase", couponHandler.Purchase)
	}

	router.Run("127.0.0.1:8080")
}
