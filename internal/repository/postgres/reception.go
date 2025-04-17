package postgres

import (
	"context"
	"fmt"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type receptionRepository struct {
	db *pgxpool.Pool
}

func NewReceptionRepository(db *pgxpool.Pool) repository.ReceptionRepository {
	return &receptionRepository{db: db}
}

func (r *receptionRepository) Create(ctx context.Context, reception *domain.Reception) error {
	if reception.PVZId == uuid.Nil || reception.DateTime.IsZero() || reception.Status == "" {
		return fmt.Errorf("pvz id, status and date time are required")
	}

	query := `INSERT INTO receptions (date_time, pvz_id, status) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, reception.DateTime, reception.PVZId, reception.Status)
	if err != nil {
		return fmt.Errorf("reception could not be created: %w", err)
	}

	return nil
}

func (r *receptionRepository) CloseLastReception(ctx context.Context, pvzId uuid.UUID) error {
	if pvzId == uuid.Nil {
		return fmt.Errorf("pvz id is required")
	}

	query := `
		UPDATE receptions
		SET status = $1 
		WHERE id = (
		    SELECT id FROM receptions
		    WHERE pvz_id = $2 AND status = 'in_progress'
		    ORDER BY date_time DESC
		    LIMIT 1
		)`

	_, err := r.db.Exec(ctx, query, "close", pvzId)
	if err != nil {
		return fmt.Errorf("reception could not be closed: %w", err)
	}

	return nil
}

func (r *receptionRepository) HasOpenReception(ctx context.Context, pvzId uuid.UUID) (bool, error) {
	if pvzId == uuid.Nil {
		return false, fmt.Errorf("pvz id is required")
	}

	var exists bool

	query := `
		SELECT EXISTS (
			SELECT 1 FROM receptions
         	WHERE pvz_id = $1 AND status = 'in_progress'
         )
	`

	err := r.db.QueryRow(ctx, query, pvzId).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check open reception: %w", err)
	}

	return exists, nil
}
