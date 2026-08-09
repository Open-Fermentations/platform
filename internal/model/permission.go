package model

import (
	"open-fermentations/internal/database/sqlc"

	"github.com/google/uuid"
)

type Permission struct {
	ID   uuid.UUID
	Name string
}

func (p *Permission) FromModel(m *sqlc.Permission) *Permission {
	p.ID = m.ID
	p.Name = m.Name

	return p
}
