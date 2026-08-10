package dto

import (
	"open-fermentations/internal/model"

	"github.com/google/uuid"
)

type BatchDeviceDTO struct {
	ID     uuid.UUID `json:"id"`
	Batch  BatchDTO  `json:"batch"`
	Device DeviceDTO `json:"device"`
}

func (d *BatchDeviceDTO) FromModel(m *model.BatchDevice) *BatchDeviceDTO {
	d.ID = m.ID
	d.Batch.FromModel(&m.Batch)
	d.Device.FromModel(&m.Device)

	return d
}
