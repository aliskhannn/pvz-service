package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	appErr "github.com/aliskhannn/pvz-service/internal/errors"
	repository_mocks "github.com/aliskhannn/pvz-service/internal/usecase/mocks/repository-mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductUseCase_AddProductToReception(t *testing.T) {
	productRepo := &repository_mocks.MockProductRepository{}
	productUC := NewProductUseCase(productRepo)

	tests := []struct {
		name        string
		user        *domain.User
		pvzId       uuid.UUID
		productType string
		repoErr     error
		expectErr   error
	}{
		{
			name:        "Valid product addition",
			user:        &domain.User{Role: constants.UserRoleEmployee},
			pvzId:       uuid.New(),
			productType: constants.ProductTypeElectronics,
			repoErr:     nil,
			expectErr:   nil,
		},
		{
			name:        "Nil user",
			user:        nil,
			pvzId:       uuid.New(),
			productType: constants.ProductTypeElectronics,
			expectErr:   appErr.ErrUserRequired,
		},
		{
			name:        "Non-employee user",
			user:        &domain.User{Role: constants.UserRoleModerator},
			pvzId:       uuid.New(),
			productType: constants.ProductTypeElectronics,
			expectErr:   appErr.ErrOnlyEmployeeAllowed,
		},
		{
			name:        "Invalid product type",
			user:        &domain.User{Role: constants.UserRoleEmployee},
			pvzId:       uuid.New(),
			productType: "",
			expectErr:   appErr.ErrInvalidProductType,
		},
		{
			name:        "Repository error",
			user:        &domain.User{Role: constants.UserRoleEmployee},
			pvzId:       uuid.New(),
			productType: constants.ProductTypeElectronics,
			repoErr:     errors.New("repository error"),
			expectErr:   appErr.ErrCreatingProduct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user != nil && tt.user.Role == constants.UserRoleEmployee && tt.productType != "" {
				productRepo.On("AddProductToReception", mock.Anything, tt.pvzId, tt.productType).
					Return(tt.repoErr).
					Once()
			}

			err := productUC.AddProductToReception(context.Background(), tt.pvzId, tt.productType, tt.user)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}

			productRepo.AssertExpectations(t)
		})
	}
}

func TestProductUseCase_DeleteLatProductFromReception(t *testing.T) {
	productRepo := &repository_mocks.MockProductRepository{}
	productUC := NewProductUseCase(productRepo)

	tests := []struct {
		name      string
		user      *domain.User
		pvzId     uuid.UUID
		repoErr   error
		expectErr error
	}{
		{
			name:      "Valid product deletion",
			user:      &domain.User{Role: constants.UserRoleEmployee},
			pvzId:     uuid.New(),
			repoErr:   nil,
			expectErr: nil,
		},
		{
			name:      "Nil user",
			user:      nil,
			pvzId:     uuid.New(),
			expectErr: appErr.ErrUserRequired,
		},
		{
			name:      "Non-employee user",
			user:      &domain.User{Role: constants.UserRoleModerator},
			pvzId:     uuid.New(),
			expectErr: appErr.ErrOnlyEmployeeAllowed,
		},
		{
			name:      "Repository error",
			user:      &domain.User{Role: constants.UserRoleEmployee},
			pvzId:     uuid.New(),
			repoErr:   errors.New("repository error"),
			expectErr: appErr.ErrDeletingLastProduct,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user != nil && tt.user.Role == constants.UserRoleEmployee {
				productRepo.On("DeleteLatProductFromReception", mock.Anything, tt.pvzId).
					Return(tt.repoErr).
					Once()
			}

			err := productUC.DeleteLatProductFromReception(context.Background(), tt.pvzId, tt.user)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}

			productRepo.AssertExpectations(t)
		})
	}
}
