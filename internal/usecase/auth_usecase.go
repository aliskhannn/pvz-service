package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliskhannn/pvz-service/internal/auth"
	"github.com/aliskhannn/pvz-service/internal/auth/jwt"
	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AuthUseCase interface {
	DummyLogin(ctx context.Context, role string) (string, error)
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, user *domain.User) error
}

type authUseCase struct {
	repo repository.UserRepository
}

func NewAuthUseCase(repo repository.UserRepository) AuthUseCase {
	return &authUseCase{repo: repo}
}

func (uc *authUseCase) DummyLogin(ctx context.Context, role string) (string, error) {
	if role != constants.UserRoleEmployee && role != constants.UserRoleModerator {
		return "", errors.New("invalid role")
	}

	userId := uuid.New()
	token, err := jwt.CreateToken(userId, role)
	if err != nil {
		return "", fmt.Errorf("token creation failed: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) Login(ctx context.Context, email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", errors.New("email and password are required")
	}

	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	err = auth.CheckPassword(password, user.Password)
	if err != nil {
		return "", errors.New("invalid email or password")
	}

	token, err := jwt.CreateToken(user.Id, user.Role)
	if err != nil {
		return "", fmt.Errorf("token creation failed: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) Register(ctx context.Context, user *domain.User) error {
	if user == nil {
		return errors.New("user is nil")
	}

	if user.Email == "" || user.Password == "" || user.Role == "" {
		return errors.New("email, password and role are required")
	}

	if user.Role != constants.UserRoleEmployee && user.Role != constants.UserRoleModerator {
		return errors.New("role must be 'employee' or 'moderator'")
	}

	existingUser, err := uc.repo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.New("user with this email already exists")
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("failed to check existing user: %w", err)
	}

	err = uc.repo.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("user creation failed: %w", err)
	}

	return nil
}
