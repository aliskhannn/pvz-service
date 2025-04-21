package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/aliskhannn/pvz-service/internal/domain"
	appErr "github.com/aliskhannn/pvz-service/internal/errors"
	"github.com/aliskhannn/pvz-service/internal/usecase"
	"github.com/aliskhannn/pvz-service/internal/usecase/mocks"
	repository_mocks "github.com/aliskhannn/pvz-service/internal/usecase/mocks/repository-mocks"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthHandler_DummyLogin(t *testing.T) {
	userRepo := &repository_mocks.MockUserRepository{}
	tokens := &mocks.MockJWTGenerator{}
	hasher := &mocks.MockPasswordHasher{}
	authUC := usecase.NewAuthUseCase(userRepo, tokens, hasher)
	handler := NewAuthHandler(authUC)

	tests := []struct {
		name           string
		body           interface{}
		token          string
		tokenErr       error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Valid role",
			body:           DummyLoginRequest{Role: "employee"},
			token:          "valid-token",
			tokenErr:       nil,
			expectedStatus: http.StatusOK,
			expectedBody:   TokenResponse{Token: "valid-token"},
		},
		{
			name:           "Empty role",
			body:           DummyLoginRequest{Role: ""},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "role is required"},
		},
		{
			name:           "Invalid role",
			body:           DummyLoginRequest{Role: "invalid"},
			tokenErr:       appErr.ErrInvalidRole,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": appErr.ErrInvalidRole.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/dummyLogin", bytes.NewReader(body))
			w := httptest.NewRecorder()

			if tt.token != "" || tt.tokenErr != nil {
				tokens.On("CreateToken", mock.Anything, "employee").
					Return(tt.token, tt.tokenErr).
					Once()
			}

			handler.DummyLogin(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var resp interface{}
			json.NewDecoder(w.Body).Decode(&resp)
			assert.Equal(t, tt.expectedBody, resp)
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	userRepo := &repository_mocks.MockUserRepository{}
	tokens := &mocks.MockJWTGenerator{}
	hasher := &mocks.MockPasswordHasher{}
	authUC := usecase.NewAuthUseCase(userRepo, tokens, hasher)
	handler := NewAuthHandler(authUC)

	tests := []struct {
		name           string
		body           interface{}
		user           *domain.User
		userErr        error
		hashErr        error
		token          string
		tokenErr       error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Valid login",
			body:           LoginRequest{Email: "test@example.com", Password: "password"},
			user:           &domain.User{Id: uuid.New(), Email: "test@example.com", Password: "hashed", Role: "employee"},
			userErr:        nil,
			hashErr:        nil,
			token:          "valid-token",
			tokenErr:       nil,
			expectedStatus: http.StatusOK,
			expectedBody:   TokenResponse{Token: "valid-token"},
		},
		{
			name:           "Missing fields",
			body:           LoginRequest{Email: "", Password: ""},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "email and password are required"},
		},
		{
			name:           "User not found",
			body:           LoginRequest{Email: "test@example.com", Password: "password"},
			userErr:        pgx.ErrNoRows,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": appErr.ErrGettingUser.Error()},
		},
		{
			name:           "Invalid password",
			body:           LoginRequest{Email: "test@example.com", Password: "password"},
			user:           &domain.User{Id: uuid.New(), Email: "test@example.com", Password: "hashed", Role: "employee"},
			hashErr:        errors.New("invalid password"),
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   map[string]string{"error": appErr.ErrInvalidAuthFields.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(body))
			w := httptest.NewRecorder()

			if tt.user != nil || tt.userErr != nil {
				userRepo.On("GetUserByEmail", mock.Anything, mock.Anything).Return(tt.user, tt.userErr).Once()
			}
			if tt.hashErr != nil || tt.user != nil {
				hasher.On("CheckPassword", mock.Anything, mock.Anything).Return(tt.hashErr).Once()
			}
			if tt.token != "" || tt.tokenErr != nil {
				tokens.On("CreateToken", mock.Anything, mock.Anything).Return(tt.token, tt.tokenErr).Once()
			}

			handler.Login(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var resp interface{}
			json.NewDecoder(w.Body).Decode(&resp)
			assert.Equal(t, tt.expectedBody, resp)
		})
	}
}

func TestAuthHandler_Register(t *testing.T) {
	userRepo := &repository_mocks.MockUserRepository{}
	tokens := &mocks.MockJWTGenerator{}
	hasher := &mocks.MockPasswordHasher{}
	authUC := usecase.NewAuthUseCase(userRepo, tokens, hasher)
	handler := NewAuthHandler(authUC)

	tests := []struct {
		name           string
		body           interface{}
		existingUser   *domain.User
		existingErr    error
		createErr      error
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			name:           "Valid registration",
			body:           domain.User{Email: "test@example.com", Password: "password", Role: "employee"},
			existingErr:    pgx.ErrNoRows,
			createErr:      nil,
			expectedStatus: http.StatusCreated,
			expectedBody:   "test@example.com",
		},
		{
			name:           "Missing fields",
			body:           domain.User{Email: "", Password: "", Role: ""},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": "email, password and role are required"},
		},
		{
			name:           "User exists",
			body:           domain.User{Email: "test@example.com", Password: "password", Role: "employee"},
			existingUser:   &domain.User{Email: "test@example.com"},
			existingErr:    nil,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": appErr.ErrUserEmailExists.Error()},
		},
		{
			name:           "Invalid role",
			body:           domain.User{Email: "test@example.com", Password: "password", Role: "invalid"},
			existingErr:    pgx.ErrNoRows,
			createErr:      appErr.ErrInvalidRole,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   map[string]string{"error": appErr.ErrInvalidRole.Error()},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader(body))
			w := httptest.NewRecorder()

			if tt.existingUser != nil || tt.existingErr != nil {
				userRepo.On("GetUserByEmail", mock.Anything, mock.Anything).Return(tt.existingUser, tt.existingErr).Once()
			}
			if tt.createErr != nil || tt.existingErr == pgx.ErrNoRows {
				userRepo.On("CreateUser", mock.Anything, mock.Anything).Return(tt.createErr).Once()
			}

			handler.Register(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)
			var resp interface{}
			json.NewDecoder(w.Body).Decode(&resp)
			assert.Equal(t, tt.expectedBody, resp)
		})
	}
}
