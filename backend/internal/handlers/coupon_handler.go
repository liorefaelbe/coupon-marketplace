package handlers

import (
	"errors"
	"net/http"

	"coupon-marketplace/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (h *CouponHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "INVALID_PRODUCT_ID",
			"message":    "Invalid product ID format",
		})
		return
	}

	product, err := h.service.GetProductByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code": "PRODUCT_NOT_FOUND",
				"message":    "Product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *CouponHandler) Purchase(c *gin.Context) {

	id := c.Param("id")

	var input services.PurchaseInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request",
		})
		return
	}

	coupon, err := h.service.Purchase(id, input.ResellerPrice)
	if err != nil {

		switch err.Error() {

		case "PRODUCT_ALREADY_SOLD":
			c.JSON(http.StatusConflict, gin.H{
				"error_code": "PRODUCT_ALREADY_SOLD",
				"message":    "Product already sold",
			})
			return

		case "RESELLER_PRICE_TOO_LOW":
			c.JSON(http.StatusBadRequest, gin.H{
				"error_code": "RESELLER_PRICE_TOO_LOW",
				"message":    "Reseller price below minimum",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "purchase failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product_id":  coupon.ProductID,
		"final_price": input.ResellerPrice,
		"value_type":  coupon.ValueType,
		"value":       coupon.Value,
	})
}
