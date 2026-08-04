package dto

import (
	"open-fermentations/internal/model"

	"github.com/google/uuid"
)

type RoleDTO struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (r *RoleDTO) FromModel(m *model.Role) *RoleDTO {
	r.ID = m.ID
	r.Name = m.Name

	return r
}
