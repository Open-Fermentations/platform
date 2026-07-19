package dto

import (
	"open-fermentations/internal/database/sqlc"
	"time"

	"github.com/google/uuid"
)

type CreateDeviceDTO struct {
	Name       string `json:"name"`
	MacAddress string `json:"macAddress"`
}

func (d *CreateDeviceDTO) ToCreateDeviceParams(userId uuid.UUID) sqlc.CreateDeviceParams {
	return sqlc.CreateDeviceParams{
		Name:       d.Name,
		MacAddress: []byte(d.MacAddress),
		UserID:     userId,
	}
}

type DeviceDTO struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	MacAddress string    `json:"macAddress"`
	UserID     uuid.UUID `json:"userId"`
	Created    time.Time `json:"created"`
	Modified   time.Time `json:"modified"`
}
