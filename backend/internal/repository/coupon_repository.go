package repository

import (
	"context"
	"errors"

	"coupon-marketplace/internal/database"
	"coupon-marketplace/internal/models"
)

type CouponRepository struct{}

func NewCouponRepository() *CouponRepository {
	return &CouponRepository{}
}

func (r *CouponRepository) Create(coupon *models.Coupon) error {
	query := `
		INSERT INTO coupons (
			product_id,
			cost_price,
			margin_percentage,
			minimum_sell_price,
			is_sold,
			value_type,
			value
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := database.DB.Exec(
		context.Background(),
		query,
		coupon.ProductID,
		coupon.CostPrice,
		coupon.MarginPercentage,
		coupon.MinimumSellPrice,
		coupon.IsSold,
		coupon.ValueType,
		coupon.Value,
	)

	return err
}

func (r *CouponRepository) GetByProductID(productID string) (*models.Coupon, error) {
	query := `
		SELECT product_id, cost_price, margin_percentage, minimum_sell_price, is_sold, value_type, value
		FROM coupons
		WHERE product_id = $1
	`

	var coupon models.Coupon

	err := database.DB.QueryRow(context.Background(), query, productID).Scan(
		&coupon.ProductID,
		&coupon.CostPrice,
		&coupon.MarginPercentage,
		&coupon.MinimumSellPrice,
		&coupon.IsSold,
		&coupon.ValueType,
		&coupon.Value,
	)

	if err != nil {
		return nil, err
	}

	return &coupon, nil
}

func (r *CouponRepository) Purchase(productID string, resellerPrice float64) (*models.Coupon, error) {

	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(context.Background())

	query := `
		SELECT product_id, cost_price, margin_percentage, minimum_sell_price, is_sold, value_type, value
		FROM coupons
		WHERE product_id = $1
		FOR UPDATE
	`

	var coupon models.Coupon

	err = tx.QueryRow(context.Background(), query, productID).Scan(
		&coupon.ProductID,
		&coupon.CostPrice,
		&coupon.MarginPercentage,
		&coupon.MinimumSellPrice,
		&coupon.IsSold,
		&coupon.ValueType,
		&coupon.Value,
	)
	if err != nil {
		return nil, err
	}

	if coupon.IsSold {
		return nil, errors.New("PRODUCT_ALREADY_SOLD")
	}

	if resellerPrice < coupon.MinimumSellPrice {
		return nil, errors.New("RESELLER_PRICE_TOO_LOW")
	}

	updateQuery := `
		UPDATE coupons
		SET is_sold = TRUE
		WHERE product_id = $1
	`

	_, err = tx.Exec(context.Background(), updateQuery, productID)
	if err != nil {
		return nil, err
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}

	return &coupon, nil
}
