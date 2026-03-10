package main

import (
	"net/http"

	"coupon-marketplace/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {

	database.Connect()
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})

	router.Run(":8080")
}
