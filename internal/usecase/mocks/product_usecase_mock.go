package mocks

import (
	"context"

	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockProductUseCase struct {
	mock.Mock
}

func (m *MockProductUseCase) AddProductToReception(ctx context.Context, pvzId uuid.UUID, productType string, user *domain.User) error {
	args := m.Called(ctx, pvzId, productType, user)
	return args.Error(0)
}

func (m *MockProductUseCase) DeleteLatProductFromReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error {
	args := m.Called(ctx, pvzId, user)
	return args.Error(0)
}
