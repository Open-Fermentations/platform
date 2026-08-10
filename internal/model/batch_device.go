package model

import "github.com/google/uuid"

type BatchDevice struct {
	ID     uuid.UUID
	Batch  Batch
	Device Device
}
