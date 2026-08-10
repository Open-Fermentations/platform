package service

import (
	"context"
	"log/slog"
	"open-fermentations/internal/database"
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/env"
	"open-fermentations/internal/model"

	"github.com/google/uuid"
)

type Service interface {
	Login(ctx context.Context, username string, password string) (*model.AuthenticatedUser, error)

	CreateBatches(ctx context.Context, d []dto.CreateBatchDTO) ([]model.Batch, error)
	DeleteBatch(ctx context.Context, id uuid.UUID) error
	SearchBatches(ctx context.Context, name string, limit, offset int) ([]model.Batch, int, error)
	GetBatchById(ctx context.Context, id uuid.UUID) (*model.Batch, error)
	UpdateBatch(ctx context.Context, id uuid.UUID, name string) (*model.Batch, error)
	AddDeviceToBatch(ctx context.Context, id, deviceId uuid.UUID) (*model.BatchDevice, error)
	RemoveDeviceFromBatch(ctx context.Context, id, deviceId uuid.UUID) error

	CreateDevices(ctx context.Context, d []dto.CreateDeviceDTO) ([]model.Device, error)
	SearchDevices(ctx context.Context, search dto.SearchDTO) (*dto.PageDTO[model.Device], error)
	GetDeviceById(ctx context.Context, id uuid.UUID) (*model.Device, error)
	DeleteDevice(ctx context.Context, id uuid.UUID) error
	UpdateDevice(ctx context.Context, id uuid.UUID, update sqlc.UpdateDeviceParams) (*model.Device, error)
}

type service struct {
	env *env.Env
	ctx context.Context
	db  database.Service
}

var _ Service = service{}

var serviceInstance *service

func New(ctx context.Context, env *env.Env, dbService database.Service) *service {
	if serviceInstance == nil {
		serviceInstance = &service{
			env: env,
			db:  dbService,
			ctx: ctx,
		}
		slog.Info("Service instance instantiated")
	}

	return serviceInstance
}
