package dto

import (
	"open-fermentations/internal/model"

	"github.com/google/uuid"
)

type PermissionDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (p *PermissionDTO) FromModel(m *model.Permission) *PermissionDTO {
	p.ID = m.ID
	p.Name = m.Name

	return p
}
