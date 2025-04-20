package usecase

import (
	"context"
	"errors"
	"github.com/aliskhannn/pvz-service/internal/auth"
	"github.com/aliskhannn/pvz-service/internal/auth/jwt"
	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	appErr "github.com/aliskhannn/pvz-service/internal/errors"
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
		return "", appErr.ErrInvalidRole
	}

	userId := uuid.New()
	token, err := jwt.CreateToken(userId, role)
	if err != nil {
		return "", appErr.ErrCreatingToken
	}

	return token, nil
}

func (uc *authUseCase) Login(ctx context.Context, email string, password string) (string, error) {
	if email == "" || password == "" {
		return "", appErr.ErrMissingAuthFields
	}

	user, err := uc.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", appErr.ErrGettingUser
	}

	err = auth.CheckPassword(password, user.Password)
	if err != nil {
		return "", appErr.ErrInvalidAuthFields
	}

	token, err := jwt.CreateToken(user.Id, user.Role)
	if err != nil {
		return "", appErr.ErrCreatingToken
	}

	return token, nil
}

func (uc *authUseCase) Register(ctx context.Context, user *domain.User) error {
	if user == nil {
		return appErr.ErrUserRequired
	}

	if user.Email == "" || user.Password == "" || user.Role == "" {
		return appErr.ErrMissingAuthFields
	}

	if user.Role != constants.UserRoleEmployee && user.Role != constants.UserRoleModerator {
		return appErr.ErrInvalidRole
	}

	existingUser, err := uc.repo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return appErr.ErrUserEmailExists
	}

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return appErr.ErrCheckingExistingUser
	}

	err = uc.repo.CreateUser(ctx, user)
	if err != nil {
		return appErr.ErrCreatingUser
	}

	return nil
}
