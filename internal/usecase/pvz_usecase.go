package usecase

import (
	"context"
	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	appErr "github.com/aliskhannn/pvz-service/internal/errors"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"time"
)

type PvzUseCase interface {
	CreatePVZ(ctx context.Context, pvz *domain.PVZ, user *domain.User) error
	GetAllPVZsWithReceptions(ctx context.Context, user *domain.User, startDate, endDate time.Time, page, limit int) ([]*domain.PVZ, error)
}

type pvzUseCase struct {
	repo repository.PVZRepository
}

func NewPvzUseCase(repo repository.PVZRepository) PvzUseCase {
	return &pvzUseCase{repo: repo}
}

func (uc *pvzUseCase) CreatePVZ(ctx context.Context, pvz *domain.PVZ, user *domain.User) error {
	if user == nil {
		return appErr.ErrUserRequired
	}

	if user.Role != constants.UserRoleModerator {
		return appErr.ErrOnlyModeratorAllowed
	}

	if pvz == nil {
		return appErr.ErrPVZIdRequired
	}

	if pvz.City != constants.PVZCityMoscow && pvz.City != constants.PVZCitySaintPetersburg && pvz.City != constants.PVZCityKazan {
		return appErr.ErrInvalidCity
	}

	err := uc.repo.CreatePVZ(ctx, pvz)
	if err != nil {
		return appErr.ErrCreatingPVZ
	}

	return nil
}

func (uc *pvzUseCase) GetAllPVZsWithReceptions(ctx context.Context, user *domain.User, startDate, endDate time.Time, offset, limit int) ([]*domain.PVZ, error) {
	if user == nil {
		return nil, appErr.ErrUserRequired
	}

	if user.Role != constants.UserRoleModerator && user.Role != constants.UserRoleEmployee {
		return nil, appErr.ErrInvalidRole
	}

	pvzs, err := uc.repo.GetAllPVZs(ctx, offset, limit)
	if err != nil {
		return nil, appErr.ErrGettingPVZs
	}

	for _, pvz := range pvzs {
		receptions, err := uc.repo.GetReceptionsByPVZId(ctx, pvz.Id, startDate, endDate)
		if err != nil {
			return nil, appErr.ErrGettingReceptions
		}

		for _, reception := range receptions {
			products, err := uc.repo.GetAllProductsFromReception(ctx, reception.Id)
			if err != nil {
				return nil, appErr.ErrGettingProducts
			}

			reception.Products = products
		}

		pvz.Receptions = receptions
	}

	return pvzs, nil
}
