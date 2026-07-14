package logging

import "log/slog"

type SlogErr interface {
	SlogErr(err error) []any
}

type Slog interface {
	Slog() []any
}

func Err(err error) slog.Attr {
	return slog.Group("error", slog.String("message", err.Error()))
}
