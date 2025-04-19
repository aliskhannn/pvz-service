package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"github.com/google/uuid"
	"time"
)

type ReceptionUseCase interface {
	CreateReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error
	CloseLastReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error
}

type receptionUseCase struct {
	repo repository.ReceptionRepository
}

func NewReceptionUseCase(repo repository.ReceptionRepository) ReceptionUseCase {
	return &receptionUseCase{repo: repo}
}

func (uc *receptionUseCase) CreateReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error {
	if user == nil {
		return errors.New("user is required")
	}

	if user.Role != constants.UserRoleEmployee {
		return errors.New("only employee can create a reception")
	}

	if pvzId == uuid.Nil {
		return fmt.Errorf("pvz id is required")
	}

	hasOpen, err := uc.repo.HasOpenReception(ctx, pvzId)
	if err != nil {
		return err
	}

	if hasOpen {
		return errors.New("pvz already has an open reception")
	}

	reception := &domain.Reception{
		PVZId:    pvzId,
		Status:   constants.ReceptionStatusInProgress,
		DateTime: time.Now(),
	}

	err = uc.repo.CreateReception(ctx, reception)
	if err != nil {
		return fmt.Errorf("failed to create reception: %w", err)
	}

	return nil
}

func (uc *receptionUseCase) CloseLastReception(ctx context.Context, pvzId uuid.UUID, user *domain.User) error {
	if user == nil {
		return errors.New("user is required")
	}

	if user.Role != constants.UserRoleEmployee {
		return errors.New("only employee can delete product from reception")
	}

	if pvzId == uuid.Nil {
		return fmt.Errorf("pvz id is required")
	}

	err := uc.repo.CloseLastReception(ctx, pvzId)
	if err != nil {
		return fmt.Errorf("failed to close last reception: %w", err)
	}

	return nil
}
