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
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	couponHandler := handlers.NewCouponHandler()

	admin := router.Group("/admin")
	{
		admin.POST("/coupons", couponHandler.CreateCoupon)
		admin.GET("/products", couponHandler.GetAdminProducts)
		admin.GET("/products/:id", couponHandler.GetAdminProductByID)
		admin.PUT("/products/:id", couponHandler.UpdateAdminProduct)
		admin.DELETE("/products/:id", couponHandler.DeleteAdminProduct)
	}

	apiV1 := router.Group("/api/v1")
	apiV1.Use(middleware.ResellerAuth())
	{
		apiV1.GET("/products", couponHandler.GetAvailableProducts)
		apiV1.GET("/products/:id", couponHandler.GetProductByID)
		apiV1.POST("/products/:id/purchase", couponHandler.Purchase)
	}

	router.Run("127.0.0.1:8080")
}
