package model

import (
	"open-fermentations/internal/database/sqlc"

	"github.com/google/uuid"
)

type Role struct {
	ID   uuid.UUID
	Name string
}

func (r *Role) FromModel(m *sqlc.Role) *Role {
	r.ID = m.ID
	r.Name = m.Name
	return r
}
