package model

import (
	"open-fermentations/internal/database/sqlc"
	"time"

	"github.com/google/uuid"
)

type Batch struct {
	ID       uuid.UUID
	Name     string
	UserID   uuid.UUID
	Created  time.Time
	Modified time.Time
}

func (b *Batch) FromModel(m sqlc.Batch) *Batch {
	b.ID = m.ID
	b.Name = m.Name
	b.UserID = m.UserID
	b.Created = m.Created
	b.Modified = m.Modified

	return b
}

func (b *Batch) FromGetBatchesRow(m sqlc.SearchBatchesRow) *Batch {
	b.ID = m.ID
	b.Name = m.Name
	b.UserID = m.UserID
	b.Created = m.Created
	b.Modified = m.Modified

	return b
}
