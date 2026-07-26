package logging

import (
	"log/slog"
	"net/http"
)

type SlogErr interface {
	SlogErr(err error) []any
}

type Slog interface {
	Slog() []any
}

func Err(err error) slog.Attr {
	return slog.Group("error", slog.String("message", err.Error()))
}

func Request(r *http.Request) []any {
	return []any{slog.Group("request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)}
}
