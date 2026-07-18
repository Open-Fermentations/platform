package service

import (
	"open-fermentations/internal/dto"
	"open-fermentations/internal/model"

	"github.com/google/uuid"
)

// CreateBatch implements [Service].
func (s service) CreateBatch(id uuid.UUID, d []dto.CreateBatchDTO) ([]model.Batch, error) {
	panic("unimplemented")
}
