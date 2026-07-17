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
