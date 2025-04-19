package domain

import (
	"github.com/google/uuid"
	"time"
)

type PVZ struct {
	Id               uuid.UUID    `json:"id" validate:"uuid"`
	RegistrationDate time.Time    `json:"registration_date"`
	City             string       `json:"city" validate:"required,oneof=Москва Санкт-Петербург Казань"`
	Receptions       []*Reception `json:"receptions" validate:"required"`
}
