package http

import (
	"encoding/json"
	"github.com/aliskhannn/pvz-service/internal/delivery/http/response"
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
		response.WriteJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		response.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized User")
		return
	}

	var req CreateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, "Invalid Request")
		return
	}

	reception, err := h.receptionUseCase.CreateReception(r.Context(), req.PVZId, user)
	if err != nil {
		status := response.MapErrorToStatusCode(err)
		response.WriteJSONError(w, status, err.Error())
		return
	}

	response.WriteJSONResponse(w, http.StatusCreated, reception)
}

func (h *ReceptionHandler) CloseLastReception(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		response.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized User")
		return
	}

	pvzIdParam := chi.URLParam(r, "pvzId")
	id, err := uuid.Parse(pvzIdParam)
	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, "Invalid ID")
		return
	}

	err = h.receptionUseCase.CloseLastReception(r.Context(), id, user)
	if err != nil {
		status := response.MapErrorToStatusCode(err)
		response.WriteJSONError(w, status, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
}
