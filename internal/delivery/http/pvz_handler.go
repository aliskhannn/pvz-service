package http

import (
	"encoding/json"
	"github.com/aliskhannn/pvz-service/internal/domain"
	"github.com/aliskhannn/pvz-service/internal/middleware"
	"github.com/aliskhannn/pvz-service/internal/usecase"
	"net/http"
	"strconv"
	"time"
)

type PVZHandler struct {
	pvzUseCase usecase.PvzUseCase
}

func NewPVZHandler(pvzUseCase usecase.PvzUseCase) *PVZHandler {
	return &PVZHandler{
		pvzUseCase: pvzUseCase,
	}
}

func (h *PVZHandler) CreatePVZ(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	var pvz domain.PVZ
	err := json.NewDecoder(r.Body).Decode(&pvz)
	if err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	err = h.pvzUseCase.CreatePVZ(r.Context(), &pvz, user)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *PVZHandler) GetAllPVZsWithReceptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	query := r.URL.Query()

	fromStr := query.Get("from")
	toStr := query.Get("to")

	var from, to time.Time
	var err error

	if fromStr != "" {
		from, err = time.Parse(time.DateOnly, fromStr)
		if err != nil {
			http.Error(w, "invalid 'from' date, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	} else {
		from = time.Time{}
	}

	if toStr != "" {
		to, err = time.Parse(time.DateOnly, toStr)
		if err != nil {
			http.Error(w, "invalid 'to' date, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	} else {
		to = time.Time{}
	}

	limitStr := query.Get("limit")
	offsetStr := query.Get("offset")

	limit := 10
	offset := 0

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
	}

	if offsetStr != "" {
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
	}

	pvzs, err := h.pvzUseCase.GetAllPVZsWithReceptions(r.Context(), user, from, to, limit, offset)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pvzs)
}
