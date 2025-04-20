package http

import (
	"encoding/json"
	"github.com/aliskhannn/pvz-service/internal/middleware"
	"github.com/aliskhannn/pvz-service/internal/usecase"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
)

type ReceptionHandler struct {
	receptionUseCase usecase.ReceptionUseCase
}

func NewReceptionHandler(receptionUseCase usecase.ReceptionUseCase) *ReceptionHandler {
	return &ReceptionHandler{
		receptionUseCase: receptionUseCase,
	}
}

type CreateRequest struct {
	PVZId uuid.UUID `json:"pvz_id"`
}

func (h *ReceptionHandler) CreateReception(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var req CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	err = h.receptionUseCase.CreateReception(r.Context(), req.PVZId, user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *ReceptionHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	pvzIdParam := chi.URLParam(r, "pvzId")
	id, err := uuid.Parse(pvzIdParam)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	err = h.receptionUseCase.CloseLastReception(r.Context(), id, user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
