package http

import (
	"encoding/json"
	"github.com/aliskhannn/pvz-service/internal/delivery/http/response"
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
		response.WriteJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		response.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized User")
		return
	}

	if user.Role != "moderator" {
		response.WriteJSONError(w, http.StatusForbidden, "access denied")
		return
	}

	var pvz domain.PVZ
	err := json.NewDecoder(r.Body).Decode(&pvz)
	if err != nil {
		response.WriteJSONError(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = h.pvzUseCase.CreatePVZ(r.Context(), &pvz, user)
	if err != nil {
		status := response.MapErrorToStatusCode(err)
		response.WriteJSONError(w, status, err.Error())
		return
	}

	response.WriteJSONResponse(w, http.StatusCreated, pvz)
}

func (h *PVZHandler) GetAllPVZsWithReceptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		response.WriteJSONError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	user, ok := middleware.GetUserFromContext(r.Context())
	if !ok || user == nil {
		response.WriteJSONError(w, http.StatusUnauthorized, "Unauthorized User")
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
			response.WriteJSONError(w, http.StatusBadRequest, "invalid 'startDate' date, use YYYY-MM-DD")
			return
		}
	} else {
		startDate = time.Time{}
	}

	if endDateStr != "" {
		endDate, err = time.Parse(time.DateOnly, endDateStr)
		if err != nil {
			response.WriteJSONError(w, http.StatusBadRequest, "invalid 'endDate' date, use YYYY-MM-DD")
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
			response.WriteJSONError(w, http.StatusBadRequest, "invalid offset")
			return
		}
	}

	if limitStr != "" {
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			response.WriteJSONError(w, http.StatusBadRequest, "invalid limit")
			return
		}
	}

	offset := (page - 1) * limit

	pvzs, err := h.pvzUseCase.GetAllPVZsWithReceptions(r.Context(), user, startDate, endDate, offset, limit)
	if err != nil {
		status := response.MapErrorToStatusCode(err)
		response.WriteJSONError(w, status, err.Error())
		return
	}

	response.WriteJSONResponse(w, http.StatusOK, pvzs)
}
