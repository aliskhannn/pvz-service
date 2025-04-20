package usecase

import (
	"context"
	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	appErr "github.com/aliskhannn/pvz-service/internal/errors"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"github.com/google/uuid"
	"time"
)

type ReceptionUseCase interface {
	CreateReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) (*domain.Reception, error)
	CloseLastReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error
}

type receptionUseCase struct {
	repo repository.ReceptionRepository
}

func NewReceptionUseCase(repo repository.ReceptionRepository) ReceptionUseCase {
	return &receptionUseCase{repo: repo}
}

func (uc *receptionUseCase) CreateReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) (*domain.Reception, error) {
	if user == nil {
		return nil, appErr.ErrUserRequired
	}

	if user.Role != constants.UserRoleEmployee {
		return nil, appErr.ErrOnlyEmployeeAllowed
	}

	if pvzId == uuid.Nil {
		return nil, appErr.ErrPVZIdRequired
	}

	hasOpen, err := uc.repo.HasOpenReception(ctx, pvzId)
	if err != nil {
		return nil, err
	}

	if hasOpen {
		return nil, appErr.ErrPVZHasOpenReception
	}

	reception := &domain.Reception{
		PVZId:    pvzId,
		Status:   constants.ReceptionStatusInProgress,
		DateTime: time.Now(),
	}

	err = uc.repo.CreateReception(ctx, reception)
	if err != nil {
		return nil, appErr.ErrCreatingReception
	}

	return reception, nil
}

func (uc *receptionUseCase) CloseLastReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error {
	if user == nil {
		return appErr.ErrUserRequired
	}

	if user.Role != constants.UserRoleEmployee {
		return appErr.ErrOnlyEmployeeAllowed
	}

	if pvzId == uuid.Nil {
		return appErr.ErrPVZIdRequired
	}

	err := uc.repo.CloseLastReception(ctx, pvzId)
	if err != nil {
		return appErr.ErrClosingLastReception
	}

	return nil
}
