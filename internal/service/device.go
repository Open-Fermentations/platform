package service

import (
	"context"
	"errors"
	"fmt"
	"open-fermentations/internal/constants"
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateDevices implements [Service].
func (s service) CreateDevices(ctx context.Context, d []dto.CreateDeviceDTO) ([]model.Device, error) {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	user, err := s.db.Querier().GetUserById(ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserDoesNotExist{ID: userId}
		}
		return nil, err
	}

	devices := []model.Device{}
	errs := []error{}
	for _, dt := range d {
		m, err := s.db.Querier().CreateDevice(ctx, dt.ToCreateDeviceParams(user.ID))
		if err != nil {
			errs = append(errs, err)
			continue
		}

		devices = append(devices, *new(model.Device).FromModel(m))
	}

	if len(errs) != 0 {
		return devices, errors.Join(errs...)
	}

	return devices, nil
}

// SearchDevices implements [Service].
func (s service) SearchDevices(ctx context.Context, search dto.SearchDTO) (*dto.PageDTO[model.Device], error) {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	devices, err := s.db.Querier().SearchDevices(ctx, sqlc.SearchDevicesParams{
		Search:    fmt.Sprintf("%%%v%%", search.Search),
		Limitval:  int32(search.Limit),
		Offsetval: int32(search.Offset),
		OrderCol:  search.OrderBy,
		Asc:       search.Asc,
		UserID:    userId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	total := 0
	if len(devices) > 0 {
		total = int(devices[0].Total)
	}

	devicesModel := make([]model.Device, len(devices))
	for i, m := range devices {
		devicesModel[i] = *new(model.Device).FromSearchDevicesRow(m)
	}

	p := dto.PageDTO[model.Device]{
		Total:  total,
		Limit:  search.Limit,
		Offset: search.Offset,
		Data:   devicesModel,
	}

	return &p, nil
}

// GetDeviceById implements [Service].
func (s service) GetDeviceById(ctx context.Context, id uuid.UUID) (*model.Device, error) {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	device, err := s.db.Querier().GetDeviceById(ctx, sqlc.GetDeviceByIdParams{
		ID:     id,
		UserID: userId,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return new(model.Device).FromModel(device), nil
}

// DeleteDevice implements [Service].
func (s service) DeleteDevice(ctx context.Context, id uuid.UUID) error {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	return s.db.Querier().RemoveDeviceById(ctx, sqlc.RemoveDeviceByIdParams{
		ID:     id,
		UserID: userId,
	})
}

// UpdateDevice implements [Service].
func (s service) UpdateDevice(ctx context.Context, id uuid.UUID, update sqlc.UpdateDeviceParams) (*model.Device, error) {
	userId := ctx.Value(constants.ContextUserIdKey).(uuid.UUID)
	update.UserID = userId
	device, err := s.db.Querier().UpdateDevice(ctx, update)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return new(model.Device).FromModel(device), nil
}
