package service

import (
	"fmt"
	"log/slog"
	"open-fermentations/internal/logging"

	"github.com/google/uuid"
)

type ErrInvalidCredentials struct {
	ID       uuid.UUID
	Username string
}

// Error implements [error].
func (e ErrInvalidCredentials) Error() string {
	return fmt.Sprintf("invalid credentials for user %v with id %v", e.Username, e.ID.String())
}

// Slog implements [logging.Slogs].
func (e ErrInvalidCredentials) SlogErr(err error) []any {
	return []any{slog.Group("error",
		slog.String("id", e.ID.String()),
		slog.String("username", e.Username),
		slog.String("message", err.Error()),
	)}
}

var _ error = ErrInvalidCredentials{}
var _ logging.SlogErr = ErrInvalidCredentials{}

type ErrUserDoesNotExist struct {
	ID uuid.UUID
}

// SlogErr implements [logging.SlogErr].
func (e ErrUserDoesNotExist) SlogErr(err error) []any {
	return []any{slog.Group("error",
		slog.String("id", e.ID.String()),
		slog.String("message", err.Error()),
	)}
}

// Error implements [error].
func (e ErrUserDoesNotExist) Error() string {
	return fmt.Sprintf("user does not exist: %v", e.ID.String())
}

var _ error = ErrUserDoesNotExist{}
var _ logging.SlogErr = ErrUserDoesNotExist{}

type ErrNotFound struct {
	ID   uuid.UUID
	Name string
}

// Error implements [error].
func (e ErrNotFound) Error() string {
	return fmt.Sprintf("could not find %s by id: %s", e.Name, e.ID.String())
}

var _ error = ErrNotFound{}
