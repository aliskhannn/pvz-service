package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"github.com/google/uuid"
)

type PvzUseCase interface {
	CreatePVZ(ctx context.Context, pvz *domain.PVZ, user *domain.User) error
	GetAllPVZs(ctx context.Context, user *domain.User) ([]*domain.PVZ, error)
	GetPVZById(ctx context.Context, pvzId uuid.UUID) (*domain.PVZ, error)
}

type pvzUseCase struct {
	repo repository.PVZRepository
}

func NewPvzUseCase(repo repository.PVZRepository) PvzUseCase {
	return &pvzUseCase{repo: repo}
}

func (uc *pvzUseCase) CreatePVZ(ctx context.Context, pvz *domain.PVZ, user *domain.User) error {
	if user == nil {
		return errors.New("user is required")
	}

	if user.Role != constants.UserRoleModerator {
		return errors.New("only moderator can create PVZ")
	}

	if pvz == nil {
		return errors.New("pvz is nil")
	}

	if pvz.City != constants.PVZCityMoscow && pvz.City != constants.PVZCitySaintPetersburg && pvz.City != constants.PVZCityKazan {
		return fmt.Errorf("city must be one of %s, %s or %s", constants.PVZCityMoscow, constants.PVZCitySaintPetersburg, constants.PVZCityKazan)
	}

	err := uc.repo.CreatePVZ(ctx, pvz)
	if err != nil {
		return fmt.Errorf("failed to create PVZ: %w", err)
	}

	return nil
}

func (uc *pvzUseCase) GetAllPVZs(ctx context.Context, user *domain.User) ([]*domain.PVZ, error) {
	if user == nil {
		return nil, errors.New("user is required")
	}

	if user.Role != constants.UserRoleModerator && user.Role != constants.UserRoleEmployee {
		return nil, errors.New("only moderator or employee can get all PVZs")
	}

	pvzs, err := uc.repo.GetAllPVZs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pvzs: %w", err)
	}

	return pvzs, nil
}

func (uc *pvzUseCase) GetPVZById(ctx context.Context, pvzId uuid.UUID) (*domain.PVZ, error) {
	if pvzId == uuid.Nil {
		return nil, errors.New("pvz id is required")
	}

	pvz, err := uc.repo.GetPVZById(ctx, pvzId)
	if err != nil {
		return nil, fmt.Errorf("failed to get pvz: %w", err)
	}

	return pvz, nil
}
