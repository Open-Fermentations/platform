package server

import (
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"open-fermentations/internal/database"
	"open-fermentations/internal/env"
	"open-fermentations/internal/logger"
)

type Server struct {
	env    *env.Env
	logger logger.Logger
	db     database.Service
}

func NewServer(env *env.Env) *http.Server {
	NewServer := &Server{
		env:    env,
		logger: logger.New(env),
		db:     database.New(env),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", env.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
