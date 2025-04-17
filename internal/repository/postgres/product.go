package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) AddProductToReception(ctx context.Context, pvzId uuid.UUID, product *domain.Product) error {
	if pvzId == uuid.Nil || product.Type == "" || product.ReceptionId == uuid.Nil {
		return fmt.Errorf("pvz id, product type and reception id are required")
	}

	if product.Type != "электроника" && product.Type != "одежда" && product.Type != "обувь" {
		return fmt.Errorf("city must be one of электроника, одежда or обувь")
	}

	var receptionId uuid.UUID
	query := `
		SELECT id FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDDER BY data_time DESC
		LIMIT 1
    `

	err := r.db.QueryRow(ctx, query, pvzId).Scan(&receptionId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("no active reception found for pvz %s", pvzId)
		}
		return fmt.Errorf("error fetching reception: %w", err)
	}

	insert := `INSERT INTO products (type, reception_id) VALUES ($1, $2)`

	cmdTag, err := r.db.Exec(ctx, insert, product.Type, receptionId)
	if err != nil {
		return fmt.Errorf("error inserting product: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no active reception found for pvz %s", pvzId)
	}

	return nil
}

func (r *productRepository) DeleteLatProductFromReception(ctx context.Context, pvzId uuid.UUID) error {
	if pvzId == uuid.Nil {
		return fmt.Errorf("pvz id is required")
	}

	query := `
		DELETE FROM products
		WHERE id (
		      SELECT id FROM products
		      WHERE reception_id = (
		    	SELECT id FROM receptions
				WHERE pvz_id = $1 AND status = 'in_progress'
				ORDDER BY data_time DESC
			  	LIMIT 1
		      )
		)
		ORDDER BY data_time DESC
		LIMIT 1
	`

	cmdTag, err := r.db.Exec(ctx, query, pvzId)
	if err != nil {
		return fmt.Errorf("error deleting reception: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no active reception found for pvz %s", pvzId)
	}

	return nil
}

func (r *productRepository) GetAllProductsFromReception(ctx context.Context, receptionId uuid.UUID) ([]*domain.Product, error) {
	if receptionId == uuid.Nil {
		return nil, fmt.Errorf("reception id is required")
	}

	products := make([]*domain.Product, 0)

	query := `SELECT id, type, reception_id, date_time FROM products WHERE reception_id = $1`

	rows, err := r.db.Query(ctx, query, receptionId)
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var product domain.Product
		err = rows.Scan(&product.Id, &product.Type, &product.ReceptionId, &product.DateTime)
		if err != nil {
			return nil, fmt.Errorf("products could not be retrieved: %w", err)
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
