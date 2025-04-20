package postgres

import (
	"context"
	"errors"
	"fmt"
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

func (r *productRepository) AddProductToReception(ctx context.Context, pvzId uuid.UUID, productType string) error {
	var receptionId uuid.UUID
	query := `
		SELECT id FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY date_time DESC
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

	cmdTag, err := r.db.Exec(ctx, insert, productType, receptionId)
	if err != nil {
		return fmt.Errorf("error inserting product: %w", err)
	}

	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("no active reception found for pvz %s", pvzId)
	}

	return nil
}

func (r *productRepository) DeleteLatProductFromReception(ctx context.Context, pvzId uuid.UUID) error {
	query := `
		DELETE FROM products
		WHERE id = (
		      SELECT id FROM products
		      WHERE reception_id = (
		    	SELECT id FROM receptions
				WHERE pvz_id = $1 AND status = 'in_progress'
				ORDER BY date_time DESC
			  	LIMIT 1
		      )
		      ORDER BY date_time DESC
			  LIMIT 1
		)
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
