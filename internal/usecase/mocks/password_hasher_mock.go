package mocks

import "github.com/stretchr/testify/mock"

type MockPasswordHasher struct {
	mock.Mock
}

func (m *MockPasswordHasher) CheckPassword(password, hash string) error {
	args := m.Called(password, hash)
	return args.Error(0)
}
