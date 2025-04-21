package repository_mocks

import (
	"context"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockReceptionRepository struct {
	mock.Mock
}

func (m *MockReceptionRepository) CreateReception(ctx context.Context, reception *domain.Reception) error {
	args := m.Called(ctx, reception)
	return args.Error(0)
}

func (m *MockReceptionRepository) CloseLastReception(ctx context.Context, pvzId uuid.UUID) error {
	args := m.Called(ctx, pvzId)
	return args.Error(0)
}

func (m *MockReceptionRepository) HasOpenReception(ctx context.Context, pvzId uuid.UUID) (bool, error) {
	args := m.Called(ctx, pvzId)
	return args.Bool(0), args.Error(1)
}
