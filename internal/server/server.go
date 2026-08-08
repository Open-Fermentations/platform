package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"

	"open-fermentations/internal/database"
	"open-fermentations/internal/env"
	"open-fermentations/internal/logging"
	"open-fermentations/internal/model"
	"open-fermentations/internal/service"
)

type Server struct {
	ctx         context.Context
	env         *env.Env
	db          database.Service
	svc         service.Service
	validate    *validator.Validate
	roles       []model.Role
	permissions []model.Permission
}

func NewServer(ctx context.Context, env *env.Env) (*http.Server, error) {
	db, err := database.New(env)
	if err != nil {
		slog.Error("failed setting up database service", logging.Err(err))
		return nil, err
	}
	newServer := &Server{
		env:      env,
		db:       db,
		svc:      service.New(ctx, env, db),
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	slog.Info("Server instantiated")

	return server, nil
}

func (s *Server) withRoles() *Server {
	rls, err := s.db.Querier().GetRolesWithPermissions(s.ctx)
	if err != nil {
		slog.Error("fetching roles for server", logging.Err(err))
	} else {
		s.roles = model.FromGetRolesWithPermissionsToRoles(rls)
	}

	return s
}

func (s *Server) withPermissions() *Server {
	ps, err := s.db.Querier().GetPermissions(s.ctx)
	if err != nil {
		slog.Error("fetching permissions for server", logging.Err(err))
	} else {
		s.permissions = make([]model.Permission, len(ps))
		for i, p := range ps {
			s.permissions[i] = model.Permission{ID: p.ID, Name: p.Name}
		}
	}

	return s
}
