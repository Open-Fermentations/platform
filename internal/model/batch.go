package model

import (
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
