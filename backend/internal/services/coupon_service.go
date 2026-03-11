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

type PurchaseInput struct {
	ResellerPrice float64 `json:"reseller_price"`
}

type UpdateCouponInput struct {
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	ImageURL         string  `json:"image_url"`
	CostPrice        float64 `json:"cost_price"`
	MarginPercentage float64 `json:"margin_percentage"`
	ValueType        string  `json:"value_type"`
	Value            string  `json:"value"`
}

func (s *CouponService) validateCouponFields(
	name string,
	imageURL string,
	costPrice float64,
	marginPercentage float64,
	valueType string,
	value string,
) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}

	if strings.TrimSpace(imageURL) == "" {
		return errors.New("image_url is required")
	}

	if costPrice < 0 {
		return errors.New("cost_price must be greater than or equal to 0")
	}

	if marginPercentage < 0 {
		return errors.New("margin_percentage must be greater than or equal to 0")
	}

	if strings.TrimSpace(valueType) == "" {
		return errors.New("value_type is required")
	}

	if strings.TrimSpace(value) == "" {
		return errors.New("value is required")
	}

	return nil
}

func (s *CouponService) CreateCoupon(input CreateCouponInput) (*models.Product, *models.Coupon, error) {
	if err := s.validateCouponFields(
		input.Name,
		input.ImageURL,
		input.CostPrice,
		input.MarginPercentage,
		input.ValueType,
		input.Value,
	); err != nil {
		return nil, nil, err
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

func (s *CouponService) Purchase(productID string, price float64) (*models.Coupon, error) {
	return s.couponRepo.Purchase(productID, price)
}

func (s *CouponService) GetAllAdminProducts() ([]repository.AdminProductDetails, error) {
	return s.productRepo.GetAllAdminProducts()
}

func (s *CouponService) GetAdminProductByID(id string) (*repository.AdminProductDetails, error) {
	return s.productRepo.GetAdminProductByID(id)
}

func (s *CouponService) DeleteProduct(id string) error {
	return s.productRepo.DeleteByID(id)
}

func (s *CouponService) UpdateCoupon(productID string, input UpdateCouponInput) error {
	if err := s.validateCouponFields(
		input.Name,
		input.ImageURL,
		input.CostPrice,
		input.MarginPercentage,
		input.ValueType,
		input.Value,
	); err != nil {
		return err
	}

	return s.couponRepo.UpdateCouponAndProduct(
		productID,
		input.Name,
		input.Description,
		input.ImageURL,
		input.CostPrice,
		input.MarginPercentage,
		input.ValueType,
		input.Value,
	)
}
