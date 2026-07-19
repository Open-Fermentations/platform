package service

import (
	"errors"
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
