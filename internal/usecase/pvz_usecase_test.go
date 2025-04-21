package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	appErr "github.com/aliskhannn/pvz-service/internal/errors"
	repository_mocks "github.com/aliskhannn/pvz-service/internal/usecase/mocks/repository-mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPvzUseCase_CreatePVZ(t *testing.T) {
	repo := &repository_mocks.MockPVZRepository{}
	pvzUC := NewPvzUseCase(repo)

	validUser := &domain.User{
		Id:   uuid.New(),
		Role: constants.UserRoleModerator,
	}

	validPVZ := &domain.PVZ{
		Id:   uuid.New(),
		City: constants.PVZCityMoscow,
	}

	tests := []struct {
		name      string
		pvz       *domain.PVZ
		user      *domain.User
		createErr error
		expectErr error
	}{
		{
			name:      "Valid PVZ creation",
			pvz:       validPVZ,
			user:      validUser,
			createErr: nil,
			expectErr: nil,
		},
		{
			name:      "Nil user",
			pvz:       validPVZ,
			user:      nil,
			expectErr: appErr.ErrUserRequired,
		},
		{
			name:      "Non-moderator user",
			pvz:       validPVZ,
			user:      &domain.User{Role: constants.UserRoleEmployee},
			expectErr: appErr.ErrOnlyModeratorAllowed,
		},
		{
			name:      "Nil PVZ",
			pvz:       nil,
			user:      validUser,
			expectErr: appErr.ErrPVZIdRequired,
		},
		{
			name:      "Invalid city",
			pvz:       &domain.PVZ{City: "invalid"},
			user:      validUser,
			expectErr: appErr.ErrInvalidCity,
		},
		{
			name:      "Repository error",
			pvz:       validPVZ,
			user:      validUser,
			createErr: errors.New("db error"),
			expectErr: appErr.ErrCreatingPVZ,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createErr != nil || (tt.expectErr == nil && tt.user != nil && tt.user.Role == constants.UserRoleModerator && tt.pvz != nil && tt.pvz.City != "invalid") {
				repo.On("CreatePVZ", mock.Anything, tt.pvz).
					Return(tt.createErr).
					Once()
			}

			err := pvzUC.CreatePVZ(context.Background(), tt.pvz, tt.user)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestPvzUseCase_GetAllPVZsWithReceptions(t *testing.T) {
	repo := &repository_mocks.MockPVZRepository{}
	pvzUC := NewPvzUseCase(repo)

	validModerator := &domain.User{
		Id:   uuid.New(),
		Role: constants.UserRoleModerator,
	}

	validEmployee := &domain.User{
		Id:   uuid.New(),
		Role: constants.UserRoleEmployee,
	}

	validPVZID := uuid.New()
	validReceptionID := uuid.New()
	validProductID := uuid.New()

	validPVZ := &domain.PVZ{
		Id: validPVZID,
	}

	validReception := &domain.Reception{
		Id:     validReceptionID,
		PVZId:  validPVZID,
		Status: constants.ReceptionStatusInProgress,
	}

	validProduct := &domain.Product{
		Id:          validProductID,
		ReceptionId: validReceptionID,
		Type:        constants.ProductTypeElectronics,
	}

	startDate := time.Now().Add(-24 * time.Hour)
	endDate := time.Now()

	tests := []struct {
		name          string
		user          *domain.User
		startDate     time.Time
		endDate       time.Time
		offset        int
		limit         int
		pvzs          []*domain.PVZ
		pvzsErr       error
		receptions    []*domain.Reception
		receptionsErr error
		products      []*domain.Product
		productsErr   error
		expected      []*domain.PVZ
		expectErr     error
	}{
		{
			name:          "Valid moderator request",
			user:          validModerator,
			startDate:     startDate,
			endDate:       endDate,
			offset:        0,
			limit:         10,
			pvzs:          []*domain.PVZ{validPVZ},
			pvzsErr:       nil,
			receptions:    []*domain.Reception{validReception},
			receptionsErr: nil,
			products:      []*domain.Product{validProduct},
			productsErr:   nil,
			expected:      []*domain.PVZ{validPVZ},
			expectErr:     nil,
		},
		{
			name:          "Valid employee request",
			user:          validEmployee,
			startDate:     startDate,
			endDate:       endDate,
			offset:        0,
			limit:         10,
			pvzs:          []*domain.PVZ{validPVZ},
			pvzsErr:       nil,
			receptions:    []*domain.Reception{validReception},
			receptionsErr: nil,
			products:      []*domain.Product{validProduct},
			productsErr:   nil,
			expected:      []*domain.PVZ{validPVZ},
			expectErr:     nil,
		},
		{
			name:      "Nil user",
			user:      nil,
			expectErr: appErr.ErrUserRequired,
		},
		{
			name:      "Invalid role",
			user:      &domain.User{Role: "invalid"},
			expectErr: appErr.ErrInvalidRole,
		},
		{
			name:      "Error getting PVZs",
			user:      validModerator,
			startDate: startDate,
			endDate:   endDate,
			offset:    0,
			limit:     10,
			pvzsErr:   errors.New("db error"),
			expectErr: appErr.ErrGettingPVZs,
		},
		{
			name:          "Error getting receptions",
			user:          validModerator,
			startDate:     startDate,
			endDate:       endDate,
			offset:        0,
			limit:         10,
			pvzs:          []*domain.PVZ{validPVZ},
			pvzsErr:       nil,
			receptionsErr: errors.New("db error"),
			expectErr:     appErr.ErrGettingReceptions,
		},
		{
			name:          "Error getting products",
			user:          validModerator,
			startDate:     startDate,
			endDate:       endDate,
			offset:        0,
			limit:         10,
			pvzs:          []*domain.PVZ{validPVZ},
			pvzsErr:       nil,
			receptions:    []*domain.Reception{validReception},
			receptionsErr: nil,
			productsErr:   errors.New("db error"),
			expectErr:     appErr.ErrGettingProducts,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.pvzsErr != nil || (tt.expectErr == nil && tt.user != nil && (tt.user.Role == constants.UserRoleModerator || tt.user.Role == constants.UserRoleEmployee)) {
				repo.On("GetAllPVZs", mock.Anything, tt.offset, tt.limit).
					Return(tt.pvzs, tt.pvzsErr).
					Once()
			}

			if tt.receptionsErr != nil || (tt.expectErr == nil && tt.user != nil && tt.pvzsErr == nil) {
				for _, pvz := range tt.pvzs {
					repo.On("GetReceptionsByPVZId", mock.Anything, pvz.Id, tt.startDate, tt.endDate).
						Return(tt.receptions, tt.receptionsErr).
						Once()
				}
			}

			if tt.productsErr != nil || (tt.expectErr == nil && tt.user != nil && tt.pvzsErr == nil && tt.receptionsErr == nil) {
				for _, reception := range tt.receptions {
					repo.On("GetProductsByReceptionId", mock.Anything, reception.Id).
						Return(tt.products, tt.productsErr).
						Once()
				}
			}

			result, err := pvzUC.GetAllPVZsWithReceptions(context.Background(), tt.user, tt.startDate, tt.endDate, tt.offset, tt.limit)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			repo.AssertExpectations(t)
		})
	}
}
