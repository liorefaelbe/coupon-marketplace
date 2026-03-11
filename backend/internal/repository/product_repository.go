package repository

import (
	"context"

	"coupon-marketplace/internal/database"
	"coupon-marketplace/internal/models"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

type AvailableProductResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageURL    string  `json:"image_url"`
	Price       float64 `json:"price"`
}

type AdminProductDetails struct {
	ID               string  `json:"id"`
	Name             string  `json:"name"`
	Description      string  `json:"description"`
	Type             string  `json:"type"`
	ImageURL         string  `json:"image_url"`
	CostPrice        float64 `json:"cost_price"`
	MarginPercentage float64 `json:"margin_percentage"`
	MinimumSellPrice float64 `json:"minimum_sell_price"`
	IsSold           bool    `json:"is_sold"`
	ValueType        string  `json:"value_type"`
	Value            string  `json:"value"`
}

func (r *ProductRepository) Create(product *models.Product) error {
	query := `
		INSERT INTO products (id, name, description, type, image_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
	`

	_, err := database.DB.Exec(
		context.Background(),
		query,
		product.ID,
		product.Name,
		product.Description,
		product.Type,
		product.ImageURL,
	)

	return err
}

func (r *ProductRepository) GetByID(id string) (*models.Product, error) {
	query := `
		SELECT id, name, description, type, image_url, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	var product models.Product

	err := database.DB.QueryRow(context.Background(), query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.Type,
		&product.ImageURL,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetAvailable() ([]models.Product, error) {
	query := `
		SELECT p.id, p.name, p.description, p.type, p.image_url, p.created_at, p.updated_at
		FROM products p
		INNER JOIN coupons c ON c.product_id = p.id
		WHERE c.is_sold = FALSE
		ORDER BY p.created_at DESC
	`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Type,
			&product.ImageURL,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, rows.Err()
}

func (r *ProductRepository) GetAvailableForAPI() ([]AvailableProductResponse, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.description,
			p.image_url,
			c.minimum_sell_price
		FROM products p
		INNER JOIN coupons c ON c.product_id = p.id
		WHERE c.is_sold = FALSE
		ORDER BY p.created_at DESC
	`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := make([]AvailableProductResponse, 0)

	for rows.Next() {
		var product AvailableProductResponse

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.ImageURL,
			&product.Price,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	return products, rows.Err()
}

func (r *ProductRepository) GetAvailableByIDForAPI(id string) (*AvailableProductResponse, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.description,
			p.image_url,
			c.minimum_sell_price
		FROM products p
		INNER JOIN coupons c ON c.product_id = p.id
		WHERE p.id = $1
	`

	var product AvailableProductResponse

	err := database.DB.QueryRow(context.Background(), query, id).Scan(
		&product.ID,
		&product.Name,
		&product.Description,
		&product.ImageURL,
		&product.Price,
	)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *ProductRepository) GetAllAdminProducts() ([]AdminProductDetails, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.description,
			p.type,
			p.image_url,
			c.cost_price,
			c.margin_percentage,
			c.minimum_sell_price,
			c.is_sold,
			c.value_type,
			c.value
		FROM products p
		INNER JOIN coupons c ON c.product_id = p.id
		ORDER BY p.created_at DESC
	`

	rows, err := database.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []AdminProductDetails

	for rows.Next() {
		var item AdminProductDetails

		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Description,
			&item.Type,
			&item.ImageURL,
			&item.CostPrice,
			&item.MarginPercentage,
			&item.MinimumSellPrice,
			&item.IsSold,
			&item.ValueType,
			&item.Value,
		)
		if err != nil {
			return nil, err
		}

		products = append(products, item)
	}

	return products, rows.Err()
}

func (r *ProductRepository) GetAdminProductByID(id string) (*AdminProductDetails, error) {
	query := `
		SELECT
			p.id,
			p.name,
			p.description,
			p.type,
			p.image_url,
			c.cost_price,
			c.margin_percentage,
			c.minimum_sell_price,
			c.is_sold,
			c.value_type,
			c.value
		FROM products p
		INNER JOIN coupons c ON c.product_id = p.id
		WHERE p.id = $1
	`

	var item AdminProductDetails

	err := database.DB.QueryRow(context.Background(), query, id).Scan(
		&item.ID,
		&item.Name,
		&item.Description,
		&item.Type,
		&item.ImageURL,
		&item.CostPrice,
		&item.MarginPercentage,
		&item.MinimumSellPrice,
		&item.IsSold,
		&item.ValueType,
		&item.Value,
	)
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *ProductRepository) DeleteByID(id string) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := database.DB.Exec(context.Background(), query, id)
	return err
}
