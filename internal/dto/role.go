package dto

import (
	"open-fermentations/internal/model"

	"github.com/google/uuid"
)

type RoleDTO struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	Permissions []PermissionDTO `json:"permissions"`
}

func (r *RoleDTO) FromModel(m *model.Role) *RoleDTO {
	r.ID = m.ID
	r.Name = m.Name

	r.Permissions = make([]PermissionDTO, len(m.Permissions))
	for i, p := range m.Permissions {
		r.Permissions[i] = *new(PermissionDTO).FromModel(&p)
	}

	return r
}
