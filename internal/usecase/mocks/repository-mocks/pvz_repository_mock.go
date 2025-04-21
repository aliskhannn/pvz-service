package repository_mocks

import (
	"context"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockPVZRepository struct {
	mock.Mock
}

func (m *MockPVZRepository) CreatePVZ(ctx context.Context, pvz *domain.PVZ) error {
	args := m.Called(ctx, pvz)
	return args.Error(0)
}

func (m *MockPVZRepository) GetAllPVZs(ctx context.Context, offset, limit int) ([]*domain.PVZ, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]*domain.PVZ), args.Error(1)
}

func (m *MockPVZRepository) GetReceptionsByPVZId(ctx context.Context, pvzId uuid.UUID, startDate, endDate time.Time) ([]*domain.Reception, error) {
	args := m.Called(ctx, pvzId, startDate, endDate)
	return args.Get(0).([]*domain.Reception), args.Error(1)
}

func (m *MockPVZRepository) GetAllProductsFromReception(ctx context.Context, receptionId uuid.UUID) ([]*domain.Product, error) {
	args := m.Called(ctx, receptionId)
	return args.Get(0).([]*domain.Product), args.Error(1)
}
