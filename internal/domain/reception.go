package domain

import (
	"github.com/google/uuid"
	"time"
)

type Reception struct {
	Id       uuid.UUID  `json:"id" validate:"uuid"`
	DateTime time.Time  `json:"date_time" validate:"required"`
	PVZId    uuid.UUID  `json:"pvz_id" validate:"required, uuid"`
	Products []*Product `json:"products" validate:"required"`
	Status   string     `json:"status" validate:"required,oneof=in_progress close"`
}
