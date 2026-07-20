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
func (s service) SearchDevices(name string, limit int, offset int) ([]model.Device, int, error) {
	devices, err := s.db.Querier().SearchDevices(s.ctx, sqlc.SearchDevicesParams{
		Name:      fmt.Sprintf("%%%v%%", name),
		Limitval:  int32(limit),
		Offsetval: int32(offset),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []model.Device{}, 0, nil
		}
		return nil, 0, err
	}

	total := 0
	if len(devices) > 0 {
		total = int(devices[0].Total)
	}

	deviceModels := make([]model.Device, len(devices))
	for i, d := range devices {
		deviceModels[i] = *new(model.Device).FromSearchDevicesRow(d)
	}

	return deviceModels, total, nil
}
