package database

import (
	"context"
	"fmt"
	"log/slog"
	"open-fermentations/internal/database/sqlc"
	"open-fermentations/internal/env"
	"strconv"

	// Added import for pgx.Row usage in suggested edit logic
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

// Service represents a service that interacts with a database.
type Service interface {
	Queries() *sqlc.Queries
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close()
}

type service struct {
	env     *env.Env
	queries *sqlc.Queries
	dbpool  *pgxpool.Pool
}

// Queries implements [Service].
func (s service) Queries() *sqlc.Queries {
	return s.queries
}

var dbInstance *service

func New(env *env.Env) (Service, error) {
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

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	dbInstance = &service{
		env:     env,
		dbpool:  pool,
		queries: sqlc.New(pool),
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
	slog.Info("disconnected from database", slog.String("database", s.env.Database.DbName))
	s.dbpool.Close()
}
