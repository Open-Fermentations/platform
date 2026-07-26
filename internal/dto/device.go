package dto

import (
	"log/slog"
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/logging"
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

// Slog implements [logging.Slog].
func (d DeviceDTO) Slog() []any {
	return []any{slog.Group("device",
		slog.String("id", d.ID.String()),
		slog.String("name", d.Name),
		slog.String("macAddress", d.MacAddress),
		slog.String("userId", d.UserID.String()),
		slog.Time("created", d.Created),
		slog.Time("modified", d.Modified),
	)}
}

func (d *DeviceDTO) FromModel(m *model.Device) *DeviceDTO {
	d.ID = m.ID
	d.Name = m.Name
	d.MacAddress = string(m.MacAddress)
	d.UserID = m.UserID
	d.Created = m.Created
	d.Modified = m.Modified

	return d
}

var _ logging.Slog = DeviceDTO{}
