package service

import (
	"errors"
	"fmt"
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// CreateDevices implements [Service].
func (s service) CreateDevices(userId uuid.UUID, d []dto.CreateDeviceDTO) ([]model.Device, error) {
	user, err := s.db.Querier().GetUserById(s.ctx, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserDoesNotExist{ID: userId}
		}
		return nil, err
	}

	devices := []model.Device{}
	errs := []error{}
	for _, dt := range d {
		m, err := s.db.Querier().CreateDevice(s.ctx, dt.ToCreateDeviceParams(user.ID))
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
func (s service) SearchDevices(search dto.SearchDTO) (*dto.PageDTO[model.Device], error) {
	devices, err := s.db.Querier().SearchDevices(s.ctx, sqlc.SearchDevicesParams{
		Search:    fmt.Sprintf("%%%v%%", search.Search),
		Limitval:  int32(search.Limit),
		Offsetval: int32(search.Offset),
		OrderCol:  search.OrderBy,
		Asc:       search.Asc,
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
func (s service) GetDeviceById(id uuid.UUID) (*model.Device, error) {
	device, err := s.db.Querier().GetDeviceById(s.ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return new(model.Device).FromModel(device), nil
}
