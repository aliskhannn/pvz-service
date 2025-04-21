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

func TestReceptionUseCase_CreateReception(t *testing.T) {
	repo := &repository_mocks.MockReceptionRepository{}
	receptionUC := NewReceptionUseCase(repo)

	validUser := &domain.User{
		Id:   uuid.New(),
		Role: constants.UserRoleEmployee,
	}

	validPVZID := uuid.New()
	validReception := &domain.Reception{
		PVZId:    validPVZID,
		Status:   constants.ReceptionStatusInProgress,
		DateTime: time.Now(),
	}

	tests := []struct {
		name       string
		pvzId      uuid.UUID
		user       *domain.User
		hasOpen    bool
		hasOpenErr error
		createErr  error
		expected   *domain.Reception
		expectErr  error
	}{
		{
			name:       "Valid reception creation",
			pvzId:      validPVZID,
			user:       validUser,
			hasOpen:    false,
			hasOpenErr: nil,
			createErr:  nil,
			expected:   validReception,
			expectErr:  nil,
		},
		{
			name:      "Nil user",
			pvzId:     validPVZID,
			user:      nil,
			expectErr: appErr.ErrUserRequired,
		},
		{
			name:      "Non-employee user",
			pvzId:     validPVZID,
			user:      &domain.User{Role: constants.UserRoleModerator},
			expectErr: appErr.ErrOnlyEmployeeAllowed,
		},
		{
			name:      "Nil PVZ ID",
			pvzId:     uuid.Nil,
			user:      validUser,
			expectErr: appErr.ErrPVZIdRequired,
		},
		{
			name:       "Has open reception",
			pvzId:      validPVZID,
			user:       validUser,
			hasOpen:    true,
			hasOpenErr: nil,
			expectErr:  appErr.ErrPVZHasOpenReception,
		},
		{
			name:       "Error checking open reception",
			pvzId:      validPVZID,
			user:       validUser,
			hasOpenErr: errors.New("db error"),
			expectErr:  errors.New("db error"),
		},
		{
			name:       "Error creating reception",
			pvzId:      validPVZID,
			user:       validUser,
			hasOpen:    false,
			hasOpenErr: nil,
			createErr:  errors.New("db error"),
			expectErr:  appErr.ErrCreatingReception,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.hasOpenErr != nil || (tt.expectErr == nil && tt.user != nil && tt.user.Role == constants.UserRoleEmployee && tt.pvzId != uuid.Nil) {
				repo.On("HasOpenReception", mock.Anything, tt.pvzId).
					Return(tt.hasOpen, tt.hasOpenErr).
					Once()
			}

			if tt.createErr != nil || (tt.expectErr == nil && tt.user != nil && tt.user.Role == constants.UserRoleEmployee && tt.pvzId != uuid.Nil && !tt.hasOpen) {
				repo.On("CreateReception", mock.Anything, mock.MatchedBy(func(r *domain.Reception) bool {
					return r.PVZId == tt.pvzId && r.Status == constants.ReceptionStatusInProgress
				})).
					Return(tt.createErr).
					Once()
			}

			result, err := receptionUC.CreateReception(context.Background(), tt.pvzId, tt.user)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.pvzId, result.PVZId)
				assert.Equal(t, constants.ReceptionStatusInProgress, result.Status)
			}

			repo.AssertExpectations(t)
		})
	}
}

func TestReceptionUseCase_CloseLastReception(t *testing.T) {
	repo := &repository_mocks.MockReceptionRepository{}
	receptionUC := NewReceptionUseCase(repo)

	validUser := &domain.User{
		Id:   uuid.New(),
		Role: constants.UserRoleEmployee,
	}

	validPVZID := uuid.New()

	tests := []struct {
		name      string
		pvzId     uuid.UUID
		user      *domain.User
		closeErr  error
		expectErr error
	}{
		{
			name:      "Valid reception closure",
			pvzId:     validPVZID,
			user:      validUser,
			closeErr:  nil,
			expectErr: nil,
		},
		{
			name:      "Nil user",
			pvzId:     validPVZID,
			user:      nil,
			expectErr: appErr.ErrUserRequired,
		},
		{
			name:      "Non-employee user",
			pvzId:     validPVZID,
			user:      &domain.User{Role: constants.UserRoleModerator},
			expectErr: appErr.ErrOnlyEmployeeAllowed,
		},
		{
			name:      "Nil PVZ ID",
			pvzId:     uuid.Nil,
			user:      validUser,
			expectErr: appErr.ErrPVZIdRequired,
		},
		{
			name:      "Error closing reception",
			pvzId:     validPVZID,
			user:      validUser,
			closeErr:  errors.New("db error"),
			expectErr: appErr.ErrClosingLastReception,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.closeErr != nil || (tt.expectErr == nil && tt.user != nil && tt.user.Role == constants.UserRoleEmployee && tt.pvzId != uuid.Nil) {
				repo.On("CloseLastReception", mock.Anything, tt.pvzId).
					Return(tt.closeErr).
					Once()
			}

			err := receptionUC.CloseLastReception(context.Background(), tt.pvzId, tt.user)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}

			repo.AssertExpectations(t)
		})
	}
}
