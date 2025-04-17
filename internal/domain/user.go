package domain

import "github.com/google/uuid"

type User struct {
	Id       uuid.UUID `json:"id" validate:"uuid"`
	Email    string    `json:"email" validate:"required,email"`
	Password string    `json:"password,omitempty" validate:"required,min=6"`
	Role     string    `json:"role" validate:"required,oneof=employee moderator"`
}
