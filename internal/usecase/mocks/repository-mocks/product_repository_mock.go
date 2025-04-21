package repository_mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) AddProductToReception(ctx context.Context, pvzId uuid.UUID, productType string) error {
	args := m.Called(ctx, pvzId, productType)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteLatProductFromReception(ctx context.Context, pvzId uuid.UUID) error {
	args := m.Called(ctx, pvzId)
	return args.Error(0)
}
