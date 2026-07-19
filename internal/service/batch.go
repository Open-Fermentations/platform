package service

import (
	"errors"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/model"

	"github.com/google/uuid"
)

// CreateBatch implements [Service].
func (s service) CreateBatch(id uuid.UUID, d []dto.CreateBatchDTO) ([]model.Batch, error) {
	u, err := s.db.Querier().GetUserById(s.ctx, id)
	if err != nil {
		return nil, err
	}

	batches := []model.Batch{}
	errs := []error{}
	for _, dto := range d {
		b, err := s.db.Querier().CreateBatch(s.ctx, *dto.ToCreateBatchParams(u.ID))
		if err != nil {
			errs = append(errs, err)
			continue
		}

		batches = append(batches, *new(model.Batch).FromModel(b))
	}

	if len(errs) != 0 {
		return batches, errors.Join(errs...)
	}

	return batches, nil
}
