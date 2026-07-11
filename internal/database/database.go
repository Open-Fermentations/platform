package database

import (
	"context"
	"fmt"
	"open-fermentations/internal/env"
	"open-fermentations/internal/logger"
	"open-fermentations/internal/repository"
	"strconv"

	// Added import for pgx.Row usage in suggested edit logic
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	Queries() *repository.Queries
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close()
}

type service struct {
	env     *env.Env
	logger  logger.Logger
	queries *repository.Queries
	dbpool  *pgxpool.Pool
}

// Queries implements [Service].
func (s service) Queries() *repository.Queries {
	return s.queries
}

var dbInstance *service

func New(env *env.Env) (Service, error) {
	logger := logger.New(env)
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance, nil
	}
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s",
		env.Database.User,
		env.Database.Password,
		env.Database.Host,
		"5555",
		env.Database.DbName,
		env.Database.Schema,
	)
	fmt.Println(connStr)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	dbInstance = &service{

		logger:  logger,
		env:     env,
		dbpool:  pool,
		queries: repository.New(pool),
	}
	return dbInstance, nil
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	stats := s.dbpool.Stat()

	return map[string]string{
		"connections":       strconv.Itoa(int(stats.AcquiredConns())),
		"idle_connections":  strconv.Itoa(int(stats.IdleConns())),
		"max_connections":   strconv.Itoa(int(stats.MaxConns())),
		"total_connections": strconv.Itoa(int(stats.TotalConns())),
		"new_connections":   strconv.Itoa(int(stats.NewConnsCount())),
	}
}

func (s *service) Close() {
	s.logger.Infof("Disconnected from database: %s", s.env.Database.DbName)
	s.dbpool.Close()
}
