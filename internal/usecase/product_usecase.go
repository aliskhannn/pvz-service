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

type ProductUseCase interface {
	AddProductToReception(ctx context.Context, pvzId uuid.UUID, product *domain.Product, user *domain.User) error
	DeleteLatProductFromReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error
}

type productUseCase struct {
	repo repository.ProductRepository
}

func NewProductUseCase(repo repository.ProductRepository) ProductUseCase {
	return &productUseCase{repo: repo}
}

func (uc *productUseCase) AddProductToReception(ctx context.Context, pvzId uuid.UUID, product *domain.Product, user *domain.User) error {
	if user == nil {
		return errors.New("user is required")
	}

	if user.Role != constants.UserRoleEmployee {
		return errors.New("only employee can add product to reception")
	}

	if product == nil {
		return fmt.Errorf("product is nil")
	}

	if pvzId == uuid.Nil || product.Type == "" || product.ReceptionId == uuid.Nil {
		return fmt.Errorf("pvz id, product type and reception id are required")
	}

	if product.Type != constants.ProductTypeElectronics && product.Type != constants.ProductsTypeCloth && product.Type != constants.ProductTypeShoes {
		return fmt.Errorf("city must be one of %s, %s or %s", constants.ProductTypeElectronics, constants.ProductsTypeCloth, constants.ProductTypeShoes)
	}

	err := uc.repo.AddProductToReception(ctx, pvzId, product)
	if err != nil {
		return fmt.Errorf("failed to add product to reception: %w", err)
	}

	return nil
}

func (uc *productUseCase) DeleteLatProductFromReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error {
	if user == nil {
		return errors.New("user is required")
	}

	if user.Role != constants.UserRoleEmployee {
		return errors.New("only employee can delete product from reception")
	}

	if pvzId == uuid.Nil {
		return fmt.Errorf("pvz id is required")
	}

	err := uc.repo.DeleteLatProductFromReception(ctx, pvzId)
	if err != nil {
		return fmt.Errorf("failed to delete latitude from reception: %w", err)
	}

	return nil
}
