package services

import (
	"errors"
	"strings"

	"coupon-marketplace/internal/models"
	"coupon-marketplace/internal/repository"

	"github.com/google/uuid"
)

type CouponService struct {
	productRepo *repository.ProductRepository
	couponRepo  *repository.CouponRepository
}

func NewCouponService() *CouponService {
	return &CouponService{
		productRepo: repository.NewProductRepository(),
		couponRepo:  repository.NewCouponRepository(),
	}
}

type CreateCouponInput struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ImageURL         string  `json:"image_url"`
	CostPrice        float64 `json:"cost_price"`
	MarginPercentage float64 `json:"margin_percentage"`
	ValueType        string  `json:"value_type"`
	Value            string  `json:"value"`
}

func (s *CouponService) CreateCoupon(input CreateCouponInput) (*models.Product, *models.Coupon, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, nil, errors.New("name is required")
	}

	if strings.TrimSpace(input.ImageURL) == "" {
		return nil, nil, errors.New("image_url is required")
	}

	if input.CostPrice < 0 {
		return nil, nil, errors.New("cost_price must be greater than or equal to 0")
	}

	if input.MarginPercentage < 0 {
		return nil, nil, errors.New("margin_percentage must be greater than or equal to 0")
	}

	if strings.TrimSpace(input.ValueType) == "" {
		return nil, nil, errors.New("value_type is required")
	}

	if strings.TrimSpace(input.Value) == "" {
		return nil, nil, errors.New("value is required")
	}

	productID := uuid.New().String()
	minimumSellPrice := input.CostPrice * (1 + input.MarginPercentage/100)

	product := &models.Product{
		ID:          productID,
		Name:        input.Name,
		Description: input.Description,
		Type:        "COUPON",
		ImageURL:    input.ImageURL,
	}

	coupon := &models.Coupon{
		ProductID:        productID,
		CostPrice:        input.CostPrice,
		MarginPercentage: input.MarginPercentage,
		MinimumSellPrice: minimumSellPrice,
		IsSold:           false,
		ValueType:        input.ValueType,
		Value:            input.Value,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, nil, err
	}

	if err := s.couponRepo.Create(coupon); err != nil {
		return nil, nil, err
	}

	return product, coupon, nil
}

func (s *CouponService) GetAvailableProducts() ([]repository.AvailableProductResponse, error) {
	return s.productRepo.GetAvailableForAPI()
}

func (s *CouponService) GetProductByID(id string) (*repository.AvailableProductResponse, error) {
	return s.productRepo.GetAvailableByIDForAPI(id)
}
