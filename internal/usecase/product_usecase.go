package usecase

import (
	"context"
	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	appErr "github.com/aliskhannn/pvz-service/internal/errors"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"github.com/google/uuid"
)

type ProductUseCase interface {
	AddProductToReception(ctx context.Context, pvzId uuid.UUID, productType string, user *domain.User) error
	DeleteLatProductFromReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error
}

type productUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) AddProductToReception(ctx context.Context, pvzId uuid.UUID, productType string, user *domain.User) error {
	if user == nil {
		return appErr.ErrUserRequired
	}

	if user.Role != constants.UserRoleEmployee {
		return appErr.ErrOnlyEmployeeAllowed
	}

	if pvzId == uuid.Nil || productType == "" {
		return appErr.ErrPVZIdAndProductTypeRequired
	}

	if productType != constants.ProductTypeElectronics && productType != constants.ProductsTypeCloth && productType != constants.ProductTypeShoes {
		return appErr.ErrInvalidProductType
	}

	err := uc.repo.AddProductToReception(ctx, pvzId, productType)
	if err != nil {
		return appErr.ErrCreatingProduct
	}

	return nil
}

func (uc *productUseCase) DeleteLatProductFromReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error {
	if user == nil {
		return appErr.ErrUserRequired
	}

	if user.Role != constants.UserRoleEmployee {
		return appErr.ErrOnlyEmployeeAllowed
	}

	if pvzId == uuid.Nil {
		return appErr.ErrPVZIdRequired
	}

	err := uc.repo.DeleteLatProductFromReception(ctx, pvzId)
	if err != nil {
		return appErr.ErrDeletingLastProduct
	}

	return nil
}
