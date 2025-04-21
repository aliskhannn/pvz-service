package mocks

import (
	"github.com/aliskhannn/pvz-service/internal/infrastructure/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type MockJWTGenerator struct {
	mock.Mock
}

func (m *MockJWTGenerator) CreateToken(userId uuid.UUID, role string) (string, error) {
	args := m.Called(userId, role)
	return args.String(0), args.Error(1)
}

func (m *MockJWTGenerator) ValidateToken(tokenString string) (*jwt.Claims, error) {
	args := m.Called(tokenString)
	return args.Get(0).(*jwt.Claims), args.Error(1)
}
