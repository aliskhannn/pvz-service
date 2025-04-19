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

func (r *pvzRepository) GetAllPVZs(ctx context.Context) ([]*domain.PVZ, error) {
	pvzs := make([]*domain.PVZ, 0)

	query := `SELECT id, registration_date, city FROM pvz`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("pvz could not be retrieved: %w", err)
	}
	defer rows.Close()

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

func (r *pvzRepository) GetPVZById(ctx context.Context, pvzId uuid.UUID) (*domain.PVZ, error) {
	var pvz domain.PVZ

	query := `SELECT id, registration_date, city FROM pvz WHERE id = $1`
	err := r.db.QueryRow(ctx, query, pvzId).Scan(&pvz.Id, &pvz.RegistrationDate, &pvz.City)
	if err != nil {
		return nil, fmt.Errorf("pvz could not be retrieved: %w", err)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("pvz not found")
	}

	return &pvz, nil
}
