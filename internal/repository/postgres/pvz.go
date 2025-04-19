package postgres

import (
	"context"
	"fmt"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type pvzRepository struct {
	db *pgxpool.Pool
}

func NewPVZRepository(db *pgxpool.Pool) repository.PVZRepository {
	return &pvzRepository{db: db}
}

func (r *pvzRepository) CreatePVZ(ctx context.Context, pvz *domain.PVZ) error {
	query := `INSERT INTO pvz (city) VALUES ($1)`
	_, err := r.db.Exec(ctx, query, pvz.City)
	if err != nil {
		return fmt.Errorf("pvz could not be created: %w", err)
	}

	return nil
}

func (r *pvzRepository) GetAllPVZs(ctx context.Context, limit, offset int) ([]*domain.PVZ, error) {
	query := `
		SELECT id, registration_date, city 
		FROM pvz
		ORDER BY registration_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("pvz could not be retrieved: %w", err)
	}
	defer rows.Close()

	var pvzs []*domain.PVZ
	for rows.Next() {
		var pvz domain.PVZ
		err = rows.Scan(&pvz.Id, &pvz.RegistrationDate, &pvz.City)
		if err != nil {
			return nil, fmt.Errorf("pvz could not be retrieved: %w", err)
		}

		pvzs = append(pvzs, &pvz)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return pvzs, nil
}

func (r *pvzRepository) GetReceptionsByPVZId(ctx context.Context, pvzId uuid.UUID, from, to time.Time) ([]*domain.Reception, error) {
	query := `
		SELECT id, pvz_id, date_time, status
		FROM receptions
		WHERE pvz_id = $1 AND date_time BETWEEN $2 AND $3
		ORDER BY date_time DESC
	`

	rows, err := r.db.Query(ctx, query, pvzId, from, to)
	if err != nil {
		return nil, fmt.Errorf("reception could not be found: %w", err)
	}
	defer rows.Close()

	var receptions []*domain.Reception
	for rows.Next() {
		var reception domain.Reception
		if err = rows.Scan(&reception.Id, &reception.PVZId, &reception.DateTime, &reception.Status); err != nil {
			return nil, fmt.Errorf("could not scan reception: %w", err)
		}

		receptions = append(receptions, &reception)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return receptions, nil
}

func (r *pvzRepository) GetAllProductsFromReception(ctx context.Context, receptionId uuid.UUID) ([]*domain.Product, error) {
	query := `
		SELECT id, type, reception_id, date_time
		FROM products 
		WHERE reception_id = $1
		ORDER BY date_time DESC
	`

	rows, err := r.db.Query(ctx, query, receptionId)
	if err != nil {
		return nil, fmt.Errorf("error fetching products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		var product domain.Product
		if err = rows.Scan(&product.Id, &product.Type, &product.ReceptionId, &product.DateTime); err != nil {
			return nil, fmt.Errorf("products could not be retrieved: %w", err)
		}

		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %w", err)
	}

	return products, nil
}
