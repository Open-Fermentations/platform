package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"open-fermentations/internal/database"
	"open-fermentations/internal/env"
	"open-fermentations/internal/service"
)

type Server struct {
	ctx context.Context
	env *env.Env
	db  database.Service
	svc service.Service
}

func NewServer(ctx context.Context, env *env.Env) (*http.Server, error) {
	db, err := database.New(env)
	if err != nil {
		slog.Error("failed setting up database service", slog.String("error", err.Error()))
		return nil, err
	}
	newServer := &Server{
		env: env,
		db:  db,
		svc: service.New(ctx, env, db),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server, nil
}
