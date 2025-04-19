package repository

import (
	"context"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/google/uuid"
	"time"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type PVZRepository interface {
	CreatePVZ(ctx context.Context, pvz *domain.PVZ) error
	GetAllPVZs(ctx context.Context, limit, offset int) ([]*domain.PVZ, error)
	GetReceptionsByPVZId(ctx context.Context, pvzId uuid.UUID, from, to time.Time) ([]*domain.Reception, error)
	GetAllProductsFromReception(ctx context.Context, receptionId uuid.UUID) ([]*domain.Product, error)
}

type ReceptionRepository interface {
	CreateReception(ctx context.Context, reception *domain.Reception) error
	CloseLastReception(ctx context.Context, pvzId uuid.UUID) error
	HasOpenReception(ctx context.Context, pvzId uuid.UUID) (bool, error)
}

type ProductRepository interface {
	AddProductToReception(ctx context.Context, pvzId uuid.UUID, product *domain.Product) error
	DeleteLatProductFromReception(ctx context.Context, pvzId uuid.UUID) error
}
