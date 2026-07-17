package server

import (
	"fmt"
	"log/slog"
	"open-fermentations/internal/logging"
	"time"
)

type ErrJWTExpired struct {
	Exp time.Time
	Now time.Time
}

// Slog implements [logging.Slog].
func (e ErrJWTExpired) Slog() []any {
	return []any{slog.Group("jwt", slog.Time("exp", e.Exp), slog.Time("now", e.Now))}
}

// Error implements [error].
func (e ErrJWTExpired) Error() string {
	return fmt.Sprintf("JWT expired. Exp: %v, Now: %v", e.Exp, e.Now)
}

var _ error = ErrJWTExpired{}

var _ logging.Slog = ErrJWTExpired{}
