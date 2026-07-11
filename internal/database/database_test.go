package database

import (
	"context"
	"log"
	"open-fermentations/internal/env"
	"strconv" // Import strconv for the suggested changes
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	database string
	password string
	username string
	host     string
	port     string
)

func mustStartPostgresContainer() (func(context.Context, ...testcontainers.TerminateOption) error, error) {
	var (
		dbName = "database"
		dbPwd  = "password"
		dbUser = "user"
	)

	dbContainer, err := postgres.Run(
		context.Background(),
		"postgres:latest",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPwd),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, err
	}

	database = dbName
	password = dbPwd
	username = dbUser

	dbHost, err := dbContainer.Host(context.Background())
	if err != nil {
		return dbContainer.Terminate, err
	}

	dbPort, err := dbContainer.MappedPort(context.Background(), "5432/tcp")
	if err != nil {
		return dbContainer.Terminate, err
	}

	host = dbHost
	port = dbPort.Port()

	return dbContainer.Terminate, err
}

func setupEnv() *env.Env {
	env := env.Env{
		LogLevel: env.LogLevelEnum.None,
		Database: env.DatabaseEnv{
			DbName:   database,
			Host:     host,
			Port:     port,
			User:     username,
			Password: password,
		},
	}

	return &env
}

func TestMain(m *testing.M) {
	teardown, err := mustStartPostgresContainer()
	if err != nil {
		log.Fatalf("could not start postgres container: %v", err)
	}

	m.Run()

	if teardown != nil && teardown(context.Background()) != nil {
		log.Fatalf("could not teardown postgres container: %v", err)
	}
}

func TestNew(t *testing.T) {
	env := setupEnv()
	srv, err := New(env)
	if err != nil {
		t.Fatalf("Error creating database service: %v", err.Error())
	}
	if srv == nil {
		t.Fatal("New() returned nil")
	}
}

func TestHealth(t *testing.T) {
	env := setupEnv()
	srv, err := New(env)
	if err != nil {
		t.Fatalf("Error creating database service: %v", err.Error())
	}

	stats := srv.Health()

	// The Health function now returns connection statistics, not fixed status strings.
	expectedKeys := []string{"connections", "idle_connections", "max_connections", "total_connections", "new_connections"}
	if len(stats) != len(expectedKeys) {
		t.Fatalf("Expected %d stats keys, got %d.", len(expectedKeys), len(stats))
	}

	for _, key := range expectedKeys {
		if _, ok := stats[key]; !ok {
			t.Fatalf("Missing expected health stat key: %s", key)
		}
		// Basic check that the value is a non-negative integer string
		if _, err := strconv.Atoi(stats[key]); err != nil {
			t.Errorf("Stat key %s has invalid integer value: %s", key, stats[key])
		}
	}
}

func TestClose(t *testing.T) {
	env := setupEnv()
	srv, err := New(env)
	if err != nil {
		t.Fatalf("Error creating database service: %v", err.Error())
	}

	srv.Close()
}
