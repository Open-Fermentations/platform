package model

import (
	"open-fermentations/internal/database/sqlc"
	"time"

	"github.com/google/uuid"
)

type Device struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	MacAddress []byte    `json:"macAddress"`
	UserID     uuid.UUID `json:"userId"`
	Created    time.Time `json:"created"`
	Modified   time.Time `json:"modified"`
}

func (d *Device) FromModel(m sqlc.Device) *Device {
	d.ID = m.ID
	d.Name = m.Name
	d.MacAddress = m.MacAddress
	d.UserID = m.UserID
	d.Created = m.Created
	d.Modified = m.Modified

	return d
}

func (d *Device) FromSearchDevicesRow(m sqlc.SearchDevicesRow) *Device {
	d.ID = m.ID
	d.Name = m.Name
	d.MacAddress = m.MacAddress
	d.UserID = m.UserID
	d.Created = m.Created
	d.Modified = m.Modified

	return d
}
