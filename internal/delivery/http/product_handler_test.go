package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/usecase/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const UserContextKey = "user"

func TestProductHandler_AddProductToReception(t *testing.T) {
	mockUseCase := new(mocks.MockProductUseCase)
	handler := NewProductHandler(mockUseCase)

	tests := []struct {
		name           string
		method         string
		body           map[string]interface{}
		user           *domain.User
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "Valid request",
			method: http.MethodPost,
			body: map[string]interface{}{
				"pvz_id":       uuid.New().String(),
				"product_type": "electronics",
			},
			user: &domain.User{
				Role: "employee",
			},
			mockSetup: func() {
				mockUseCase.On("AddProductToReception", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:   "Invalid method",
			method: http.MethodGet,
			body:   map[string]interface{}{},
			user:   &domain.User{},
			mockSetup: func() {
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "Invalid product type",
			method: http.MethodPost,
			body: map[string]interface{}{
				"pvz_id":       uuid.New().String(),
				"product_type": "invalid",
			},
			user: &domain.User{
				Role: "employee",
			},
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			body, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(tt.method, "/api/v1/products", bytes.NewReader(body))
			ctx := context.WithValue(req.Context(), UserContextKey, tt.user)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.AddProductToReception(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}

func TestProductHandler_DeleteLatProductFromReception(t *testing.T) {
	mockUseCase := new(mocks.MockProductUseCase)
	handler := NewProductHandler(mockUseCase)

	tests := []struct {
		name           string
		method         string
		pvzId          string
		user           *domain.User
		mockSetup      func()
		expectedStatus int
	}{
		{
			name:   "Valid request",
			method: http.MethodDelete,
			pvzId:  uuid.New().String(),
			user: &domain.User{
				Role: "employee",
			},
			mockSetup: func() {
				mockUseCase.On("DeleteLatProductFromReception", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:   "Invalid method",
			method: http.MethodPost,
			pvzId:  uuid.New().String(),
			user:   &domain.User{},
			mockSetup: func() {
			},
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:   "Invalid PVZ ID",
			method: http.MethodDelete,
			pvzId:  "invalid-uuid",
			user: &domain.User{
				Role: "employee",
			},
			mockSetup: func() {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()

			req := httptest.NewRequest(tt.method, "/api/v1/products/"+tt.pvzId, nil)
			ctx := context.WithValue(req.Context(), UserContextKey, tt.user)
			req = req.WithContext(ctx)

			rec := httptest.NewRecorder()
			handler.DeleteLatProductFromReception(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			mockUseCase.AssertExpectations(t)
		})
	}
}
