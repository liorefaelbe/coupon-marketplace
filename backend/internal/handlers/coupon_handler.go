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
			"error_code": "INVALID_REQUEST",
			"message":    "Invalid request body",
		})
		return
	}

	product, coupon, err := h.service.CreateCoupon(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
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
			"error_code": "INTERNAL_ERROR",
			"message":    "Failed to fetch products",
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
			"error_code": "INTERNAL_ERROR",
			"message":    "Failed to fetch product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *CouponHandler) Purchase(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "INVALID_PRODUCT_ID",
			"message":    "Invalid product ID format",
		})
		return
	}

	var input services.PurchaseInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "INVALID_REQUEST",
			"message":    "Invalid request body",
		})
		return
	}

	coupon, err := h.service.Purchase(id, input.ResellerPrice)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code": "PRODUCT_NOT_FOUND",
				"message":    "Product not found",
			})
			return
		}

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
			"error_code": "INTERNAL_ERROR",
			"message":    "Purchase failed",
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

func (h *CouponHandler) GetAdminProducts(c *gin.Context) {
	products, err := h.service.GetAllAdminProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": "INTERNAL_ERROR",
			"message":    "Failed to fetch admin products",
		})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *CouponHandler) GetAdminProductByID(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "INVALID_PRODUCT_ID",
			"message":    "Invalid product ID format",
		})
		return
	}

	product, err := h.service.GetAdminProductByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code": "PRODUCT_NOT_FOUND",
				"message":    "Product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": "INTERNAL_ERROR",
			"message":    "Failed to fetch admin product",
		})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *CouponHandler) UpdateAdminProduct(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "INVALID_PRODUCT_ID",
			"message":    "Invalid product ID format",
		})
		return
	}

	var input services.UpdateCouponInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "INVALID_REQUEST",
			"message":    "Invalid request body",
		})
		return
	}

	err := h.service.UpdateCoupon(id, input)
	if err != nil {
		if err.Error() == "PRODUCT_NOT_FOUND" {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code": "PRODUCT_NOT_FOUND",
				"message":    "Product not found",
			})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "VALIDATION_ERROR",
			"message":    err.Error(),
		})
		return
	}

	updated, err := h.service.GetAdminProductByID(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "Product updated successfully",
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *CouponHandler) DeleteAdminProduct(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error_code": "INVALID_PRODUCT_ID",
			"message":    "Invalid product ID format",
		})
		return
	}

	existing, err := h.service.GetAdminProductByID(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{
				"error_code": "PRODUCT_NOT_FOUND",
				"message":    "Product not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": "INTERNAL_ERROR",
			"message":    "Failed to fetch product before delete",
		})
		return
	}

	if existing.IsSold {
		c.JSON(http.StatusConflict, gin.H{
			"error_code": "PRODUCT_ALREADY_SOLD",
			"message":    "Cannot delete sold product",
		})
		return
	}

	if err := h.service.DeleteProduct(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error_code": "INTERNAL_ERROR",
			"message":    "Failed to delete product",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
