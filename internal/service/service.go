package service

import (
	"context"
	"log/slog"
	"open-fermentations/internal/database"
	"open-fermentations/internal/dto"
	"open-fermentations/internal/env"
	"open-fermentations/internal/model"

	"github.com/google/uuid"
)

type Service interface {
	Login(username string, password string) (*model.User, error)

	CreateBatch(id uuid.UUID, d []dto.CreateBatchDTO) ([]model.Batch, error)
	DeleteBatch(id uuid.UUID) error
	GetBatches(name string, limit, offset int) ([]model.Batch, int, error)
	GetBatchById(id uuid.UUID) (*model.Batch, error)
	UpdateBatch(id uuid.UUID, name string) (*model.Batch, error)
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
