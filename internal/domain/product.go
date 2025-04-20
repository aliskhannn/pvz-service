package domain

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	Id          uuid.UUID `json:"id" validate:"uuid"`
	DateTime    time.Time `json:"date_time"`
	Type        string    `json:"type" validate:"required,oneof=Электроника Одежда Обувь"`
	PVZId       uuid.UUID `json:"pvz_id" validate:"required,uuid"`
	ReceptionId uuid.UUID `json:"reception_id"`
}
