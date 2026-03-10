package handlers

import (
	"net/http"

	"coupon-marketplace/internal/services"

	"github.com/gin-gonic/gin"
)

type CouponHandler struct {
	service *services.CouponService
}

func NewCouponHandler() *CouponHandler {
	return &CouponHandler{
		service: services.NewCouponService(),
	}
}

func (h *CouponHandler) CreateCoupon(c *gin.Context) {

	var input services.CreateCouponInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	product, coupon, err := h.service.CreateCoupon(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"product": product,
		"coupon":  coupon,
	})
}

func (h *CouponHandler) GetAvailableProducts(c *gin.Context) {
	products, err := h.service.GetAvailableProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch products",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}
