package response

import (
	appErr "github.com/aliskhannn/pvz-service/internal/errors"
	"net/http"
)

// MapErrorToStatusCode маппит бизнес-ошибки на HTTP-статусы.
func MapErrorToStatusCode(err error) int {
	switch err {
	// 401 Unauthorized
	case appErr.ErrUnauthorized,
		appErr.ErrUserRequired:
		return http.StatusUnauthorized

	// 403 Forbidden
	case appErr.ErrForbidden,
		appErr.ErrOnlyEmployeeAllowed,
		appErr.ErrOnlyModeratorAllowed:
		return http.StatusForbidden

	// 404 Not Found
	case appErr.ErrNotFound:
		return http.StatusNotFound

	// 400 Bad Request — валидация, дубликаты, отсутствие полей, бизнес-ошибки клиента
	case appErr.ErrBadRequest,
		appErr.ErrValidation,
		appErr.ErrAlreadyExists,
		appErr.ErrUserAlreadyExists,
		appErr.ErrUserEmailExists,
		appErr.ErrInvalidRole,
		appErr.ErrMissingAuthFields,
		appErr.ErrInvalidAuthFields,
		appErr.ErrPVZIdRequired,
		appErr.ErrPVZRequired,
		appErr.ErrInvalidCity,
		appErr.ErrPVZHasOpenReception,
		appErr.ErrPVZIdAndProductTypeRequired,
		appErr.ErrInvalidProductType:
		return http.StatusBadRequest

	// 500 Internal Server Error — технические ошибки
	case appErr.ErrInternal,
		appErr.ErrCheckingExistingUser,
		appErr.ErrCreatingUser,
		appErr.ErrGettingUser,
		appErr.ErrCreatingToken,
		appErr.ErrCreatingPVZ,
		appErr.ErrGettingPVZs,
		appErr.ErrGettingReceptions,
		appErr.ErrCreatingReception,
		appErr.ErrClosingLastReception,
		appErr.ErrGettingProducts,
		appErr.ErrCreatingProduct,
		appErr.ErrDeletingLastProduct:
		return http.StatusInternalServerError

	// fallback на 500
	default:
		return http.StatusInternalServerError
	}
}
