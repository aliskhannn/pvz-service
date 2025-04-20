package http

import (
	"encoding/json"
	"github.com/aliskhannn/pvz-service/internal/delivery/http/response"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/usecase"
	"net/http"
)

type AuthHandler struct {
	authUseCase usecase.AuthUseCase
}

func NewAuthHandler(authUseCase usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
	}
}

type DummyLoginRequest struct {
	Role string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	Token string `json:"token"`
}

func (h *AuthHandler) DummyLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var req DummyLoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if req.Role == "" {
		response.WriteJSONError(w, http.StatusBadRequest, "role is required")
		return
	}

	token, err := h.authUseCase.DummyLogin(r.Context(), req.Role)
	if err != nil {
		status := response.MapErrorToStatusCode(err)
		response.WriteJSONError(w, status, err.Error())
		return
	}

	response.WriteJSONResponse(w, http.StatusOK, TokenResponse{Token: token})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var loginReq LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if loginReq.Email == "" || loginReq.Password == "" {
		response.WriteJSONError(w, http.StatusBadRequest, "email and password are required")
		return
	}

	token, err := h.authUseCase.Login(r.Context(), loginReq.Email, loginReq.Password)
	if err != nil {
		status := response.MapErrorToStatusCode(err)
		response.WriteJSONError(w, status, err.Error())
		return
	}

	response.WriteJSONResponse(w, http.StatusOK, TokenResponse{Token: token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var user domain.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	if user.Email == "" || user.Password == "" || user.Role == "" {
		response.WriteJSONError(w, http.StatusBadRequest, "email, password and role are required")
		return
	}

	err = h.authUseCase.Register(r.Context(), &user)
	if err != nil {
		status := response.MapErrorToStatusCode(err)
		response.WriteJSONError(w, status, err.Error())
		return
	}

	response.WriteJSONResponse(w, http.StatusCreated, user.Email)
}
