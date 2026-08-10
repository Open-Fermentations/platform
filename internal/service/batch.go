package service

import (
	"context"
	"errors"
	"open-fermentations/internal/constants"
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateBatches implements [Service].
func (s service) CreateBatches(ctx context.Context, d []dto.CreateBatchDTO) ([]model.Batch, error) {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	u, err := s.db.Querier().GetUserById(s.ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserDoesNotExist{ID: userId}
		}
		return nil, err
	}

	batches := []model.Batch{}
	errs := []error{}
	for _, dt := range d {
		b, err := s.db.Querier().CreateBatch(ctx, *dt.ToCreateBatchParams(u.ID))
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
func (s service) DeleteBatch(ctx context.Context, id uuid.UUID) error {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	return s.db.Querier().DeleteBatch(s.ctx, sqlc.DeleteBatchParams{
		ID:     id,
		UserID: userId,
	})
}

// SearchBatches implements [Service].
func (s service) SearchBatches(ctx context.Context, name string, limit int, offset int) ([]model.Batch, int, error) {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	batches, err := s.db.Querier().SearchBatches(ctx, sqlc.SearchBatchesParams{
		Name:      name,
		Limitval:  int32(limit),
		Offsetval: int32(offset),
		UserID:    userId,
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
		batchModels[i] = *new(model.Batch).FromSearchBatchesRow(b)
	}

	return batchModels, total, nil
}

// GetBatchById implements [Service].
func (s service) GetBatchById(ctx context.Context, id uuid.UUID) (*model.Batch, error) {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	batch, err := s.db.Querier().GetBatchById(s.ctx, sqlc.GetBatchByIdParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return new(model.Batch).FromModel(batch), nil
}

// UpdateBatch implements [Service].
func (s service) UpdateBatch(ctx context.Context, id uuid.UUID, name string) (*model.Batch, error) {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	batch, err := s.db.Querier().UpdateBatch(s.ctx, sqlc.UpdateBatchParams{
		ID:     id,
		Name:   name,
		UserID: userId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return new(model.Batch).FromModel(batch), nil
}

// AddDeviceToBatch implements [Service].
func (s service) AddDeviceToBatch(ctx context.Context, id uuid.UUID, deviceId uuid.UUID) (*model.BatchDevice, error) {
	batch, err := s.GetBatchById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound{ID: id, Name: "batch id"}
		}
		return nil, err
	}

	device, err := s.GetDeviceById(ctx, deviceId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound{ID: id, Name: "device id"}
		}
		return nil, err
	}

	batchDevice, err := s.db.Querier().AddDeviceToBatch(ctx, sqlc.AddDeviceToBatchParams{
		BatchID:  id,
		DeviceID: deviceId,
	})
	if err != nil {
		return nil, err
	}

	batchDeviceModel := model.BatchDevice{
		ID:     batchDevice.ID,
		Batch:  *batch,
		Device: *device,
	}

	return &batchDeviceModel, nil
}

// RemoveDeviceFromBatch implements [Service].
func (s service) RemoveDeviceFromBatch(ctx context.Context, id uuid.UUID, deviceId uuid.UUID) error {
	var err error
	_, err = s.GetBatchById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound{ID: id, Name: "batch id"}
		}
		return err
	}

	_, err = s.GetDeviceById(ctx, deviceId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound{ID: id, Name: "device id"}
		}
		return err
	}

	return s.db.Querier().RemoveDeviceFromBatch(ctx, sqlc.RemoveDeviceFromBatchParams{
		BatchID:  id,
		DeviceID: deviceId,
	})
}
