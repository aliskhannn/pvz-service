package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/aliskhannn/pvz-service/internal/constants"
	"github.com/aliskhannn/pvz-service/internal/domain"
	appErr "github.com/aliskhannn/pvz-service/internal/errors"
	"github.com/aliskhannn/pvz-service/internal/usecase/mocks"
	repository_mocks "github.com/aliskhannn/pvz-service/internal/usecase/mocks/repository-mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthUseCase_DummyLogin(t *testing.T) {
	userRepo := &repository_mocks.MockUserRepository{}
	tokens := &mocks.MockJWTGenerator{}
	hasher := &mocks.MockPasswordHasher{}
	authUC := NewAuthUseCase(userRepo, tokens, hasher)

	tests := []struct {
		name      string
		role      string
		token     string
		tokenErr  error
		expected  string
		expectErr error
	}{
		{
			name:      "Valid employee role",
			role:      constants.UserRoleEmployee,
			token:     "valid-token",
			tokenErr:  nil,
			expected:  "valid-token",
			expectErr: nil,
		},
		{
			name:      "Valid moderator role",
			role:      constants.UserRoleModerator,
			token:     "valid-token",
			tokenErr:  nil,
			expected:  "valid-token",
			expectErr: nil,
		},
		{
			name:      "Invalid role",
			role:      "invalid",
			expectErr: appErr.ErrInvalidRole,
		},
		{
			name:      "Token creation error",
			role:      constants.UserRoleEmployee,
			tokenErr:  errors.New("token error"),
			expectErr: appErr.ErrCreatingToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.token != "" || tt.tokenErr != nil {
				tokens.On("CreateToken", mock.Anything, tt.role).
					Return(tt.token, tt.tokenErr).
					Once()
			}

			result, err := authUC.DummyLogin(context.Background(), tt.role)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			tokens.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	userRepo := &repository_mocks.MockUserRepository{}
	tokens := &mocks.MockJWTGenerator{}
	hasher := &mocks.MockPasswordHasher{}
	authUC := NewAuthUseCase(userRepo, tokens, hasher)

	userID := uuid.New()
	user := &domain.User{
		Id:       userID,
		Email:    "test@example.com",
		Password: "hashed-password",
		Role:     constants.UserRoleEmployee,
	}

	tests := []struct {
		name      string
		email     string
		password  string
		user      *domain.User
		userErr   error
		hashErr   error
		token     string
		tokenErr  error
		expected  string
		expectErr error
	}{
		{
			name:      "Valid login",
			email:     "test@example.com",
			password:  "password",
			user:      user,
			userErr:   nil,
			hashErr:   nil,
			token:     "valid-token",
			tokenErr:  nil,
			expected:  "valid-token",
			expectErr: nil,
		},
		{
			name:      "Missing email",
			email:     "",
			password:  "password",
			expectErr: appErr.ErrMissingAuthFields,
		},
		{
			name:      "Missing password",
			email:     "test@example.com",
			password:  "",
			expectErr: appErr.ErrMissingAuthFields,
		},
		{
			name:      "User not found",
			email:     "test@example.com",
			password:  "password",
			userErr:   pgx.ErrNoRows,
			expectErr: appErr.ErrGettingUser,
		},
		{
			name:      "Invalid password",
			email:     "test@example.com",
			password:  "wrong-password",
			user:      user,
			userErr:   nil,
			hashErr:   errors.New("invalid password"),
			expectErr: appErr.ErrInvalidAuthFields,
		},
		{
			name:      "Token creation error",
			email:     "test@example.com",
			password:  "password",
			user:      user,
			userErr:   nil,
			hashErr:   nil,
			tokenErr:  errors.New("token error"),
			expectErr: appErr.ErrCreatingToken,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.email != "" {
				userRepo.On("GetUserByEmail", mock.Anything, tt.email).
					Return(tt.user, tt.userErr).
					Once()
			}

			if tt.user != nil && tt.userErr == nil {
				hasher.On("CheckPassword", tt.password, tt.user.Password).
					Return(tt.hashErr).
					Once()
			}

			if tt.user != nil && tt.userErr == nil && tt.hashErr == nil {
				tokens.On("CreateToken", tt.user.Id, tt.user.Role).
					Return(tt.token, tt.tokenErr).
					Once()
			}

			result, err := authUC.Login(context.Background(), tt.email, tt.password)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}

			userRepo.AssertExpectations(t)
			hasher.AssertExpectations(t)
			tokens.AssertExpectations(t)
		})
	}
}

func TestAuthUseCase_Register(t *testing.T) {
	userRepo := &repository_mocks.MockUserRepository{}
	tokens := &mocks.MockJWTGenerator{}
	hasher := &mocks.MockPasswordHasher{}
	authUC := NewAuthUseCase(userRepo, tokens, hasher)

	userID := uuid.New()
	validUser := &domain.User{
		Id:       userID,
		Email:    "test@example.com",
		Password: "password",
		Role:     constants.UserRoleEmployee,
	}

	tests := []struct {
		name        string
		user        *domain.User
		existing    *domain.User
		existingErr error
		createErr   error
		expectErr   error
	}{
		{
			name:        "Valid registration",
			user:        validUser,
			existingErr: pgx.ErrNoRows,
			createErr:   nil,
			expectErr:   nil,
		},
		{
			name:      "Nil user",
			user:      nil,
			expectErr: appErr.ErrUserRequired,
		},
		{
			name:      "Missing email",
			user:      &domain.User{Password: "password", Role: constants.UserRoleEmployee},
			expectErr: appErr.ErrMissingAuthFields,
		},
		{
			name:      "Missing password",
			user:      &domain.User{Email: "test@example.com", Role: constants.UserRoleEmployee},
			expectErr: appErr.ErrMissingAuthFields,
		},
		{
			name:      "Missing role",
			user:      &domain.User{Email: "test@example.com", Password: "password"},
			expectErr: appErr.ErrMissingAuthFields,
		},
		{
			name:      "Invalid role",
			user:      &domain.User{Email: "test@example.com", Password: "password", Role: "invalid"},
			expectErr: appErr.ErrInvalidRole,
		},
		{
			name:        "User exists",
			user:        validUser,
			existing:    validUser,
			existingErr: nil,
			expectErr:   appErr.ErrUserEmailExists,
		},
		{
			name:        "Error checking existing user",
			user:        validUser,
			existingErr: errors.New("db error"),
			expectErr:   appErr.ErrCheckingExistingUser,
		},
		{
			name:        "Error creating user",
			user:        validUser,
			existingErr: pgx.ErrNoRows,
			createErr:   errors.New("create error"),
			expectErr:   appErr.ErrCreatingUser,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.user != nil && tt.user.Email != "" && tt.user.Password != "" && tt.user.Role != "" &&
				(tt.user.Role == constants.UserRoleEmployee || tt.user.Role == constants.UserRoleModerator) {
				userRepo.On("GetUserByEmail", mock.Anything, tt.user.Email).
					Return(tt.existing, tt.existingErr).
					Once()
			}

			if tt.createErr != nil || (tt.existingErr == pgx.ErrNoRows && tt.user != nil && tt.user.Email != "" &&
				tt.user.Password != "" && tt.user.Role != "" &&
				(tt.user.Role == constants.UserRoleEmployee || tt.user.Role == constants.UserRoleModerator)) {
				userRepo.On("CreateUser", mock.Anything, tt.user).
					Return(tt.createErr).
					Once()
			}

			err := authUC.Register(context.Background(), tt.user)

			if tt.expectErr != nil {
				assert.ErrorIs(t, err, tt.expectErr)
			} else {
				assert.NoError(t, err)
			}

			userRepo.AssertExpectations(t)
		})
	}
}
