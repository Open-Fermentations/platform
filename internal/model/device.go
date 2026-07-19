package model

import (
	"open-fermentations/internal/database/sqlc"
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	MacAddress string    `json:"macAddress"`
	UserID     uuid.UUID `json:"userId"`
	Created    time.Time `json:"created"`
	Modified   time.Time `json:"modified"`
}

func (d *Device) FromModel(m sqlc.Device) *Device {
	d.ID = m.ID
	d.Name = m.Name
	d.MacAddress = string(m.MacAddress)
	d.UserID = m.UserID
	d.Created = m.Created
	d.Modified = m.Modified

	return d
}
