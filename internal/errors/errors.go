package errors

import "errors"

var (
	ErrUnauthorized  = errors.New("unauthorized")
	ErrForbidden     = errors.New("forbidden")
	ErrBadRequest    = errors.New("bad request")
	ErrValidation    = errors.New("validation error")
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInternal      = errors.New("internal error")

	ErrUserRequired         = errors.New("user is required")
	ErrUserAlreadyExists    = errors.New("user already exists")
	ErrUserEmailExists      = errors.New("user with this email already exists")
	ErrCheckingExistingUser = errors.New("error checking existing user")
	ErrInvalidRole          = errors.New("invalid role")
	ErrOnlyEmployeeAllowed  = errors.New("only employee is allowed")
	ErrOnlyModeratorAllowed = errors.New("only moderator is allowed")
	ErrCreatingUser         = errors.New("error creating user")
	ErrGettingUser          = errors.New("error getting user")
	ErrCreatingToken        = errors.New("error creating token")
	ErrMissingAuthFields    = errors.New("email, password or role is required")
	ErrInvalidAuthFields    = errors.New("invalid email, password or type")

	ErrPVZIdRequired = errors.New("pvz id is required")
	ErrPVZRequired   = errors.New("pvz is required")
	ErrInvalidCity   = errors.New("invalid city")
	ErrCreatingPVZ   = errors.New("error creating pvz")
	ErrGettingPVZs   = errors.New("error getting pvzs")

	ErrGettingReceptions    = errors.New("error getting receptions")
	ErrPVZHasOpenReception  = errors.New("pvz already has an open reception")
	ErrCreatingReception    = errors.New("error creating reception")
	ErrClosingLastReception = errors.New("error closing last reception")

	ErrGettingProducts             = errors.New("error getting products")
	ErrPVZIdAndProductTypeRequired = errors.New("pvz id and product type is required")
	ErrInvalidProductType          = errors.New("invalid product type")
	ErrCreatingProduct             = errors.New("error creating product")
	ErrDeletingLastProduct         = errors.New("error deleting last product from reception")
)
