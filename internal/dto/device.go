package dto

import (
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/model"
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

func (d *DeviceDTO) FromModel(m *model.Device) *DeviceDTO {
	d.ID = m.ID
	d.Name = m.Name
	d.MacAddress = m.MacAddress
	d.UserID = m.UserID
	d.Created = m.Created
	d.Modified = m.Modified

	return d
}
