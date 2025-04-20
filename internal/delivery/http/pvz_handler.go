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

	startDateStr := query.Get("startDate")
	endDateStr := query.Get("endDate")

	var startDate, endDate time.Time
	var err error

	if startDateStr != "" {
		startDate, err = time.Parse(time.DateOnly, startDateStr)
		if err != nil {
			http.Error(w, "invalid 'startDate' date, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	} else {
		startDate = time.Time{}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.DateOnly, endDateStr)
		if err != nil {
			http.Error(w, "invalid 'endDate' date, use YYYY-MM-DD", http.StatusBadRequest)
			return
		}
	} else {
		endDate = time.Time{}
	}

	pageStr := query.Get("page")
	limitStr := query.Get("limit")

	page := 0
	limit := 10

	if pageStr != "" {
		page, err = strconv.Atoi(pageStr)
		if err != nil || page < 0 {
			http.Error(w, "invalid offset", http.StatusBadRequest)
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			http.Error(w, "invalid limit", http.StatusBadRequest)
			return
		}
	}

	offset := (page - 1) * limit

	pvzs, err := h.pvzUseCase.GetAllPVZsWithReceptions(r.Context(), user, startDate, endDate, offset, limit)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(pvzs)
}
