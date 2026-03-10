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
