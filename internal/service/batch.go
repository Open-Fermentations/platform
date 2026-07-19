package service

import (
	"errors"
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateBatch implements [Service].
func (s service) CreateBatch(id uuid.UUID, d []dto.CreateBatchDTO) ([]model.Batch, error) {
	u, err := s.db.Querier().GetUserById(s.ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
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

// DeleteBatch implements [Service].
func (s service) DeleteBatch(id uuid.UUID) error {
	return s.db.Querier().DeleteBatch(s.ctx, id)
}

// GetBatches implements [Service].
func (s service) GetBatches(name string, limit int, offset int) ([]model.Batch, int, error) {
	batches, err := s.db.Querier().GetBatches(s.ctx, sqlc.GetBatchesParams{
		Column1: name,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.Batch{}, 0, nil
		}
		return nil, 0, err
	}

	total := 0
	if len(batches) > 0 {
		total = int(batches[0].Total)
	}

	batchModels := make([]model.Batch, len(batches))
	for i, b := range batches {
		batchModels[i] = *new(model.Batch).FromGetBatchesRow(b)
	}

	return batchModels, total, nil
}

// GetBatchById implements [Service].
func (s service) GetBatchById(id uuid.UUID) (*model.Batch, error) {
	batch, err := s.db.Querier().GetBatchById(s.ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return new(model.Batch).FromModel(batch), nil
}
