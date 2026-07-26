package route

import (
	"context"
	"log/slog"
	"net/http"
	"open-fermentations/internal/logging"

	"github.com/google/uuid"
)

func pathUuidMiddleware(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rawId := r.PathValue(key)
			id, err := uuid.Parse(rawId)
			if err != nil {
				slog.Error("unable to parse path uuid", []any{slog.String("key", key), logging.Request(r)}...)
			}

			ctx := context.WithValue(r.Context(), key, id)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (r *Route) WithPathUuid(key string) *Route {
	return r.WithMiddleware(pathUuidMiddleware(key))
}
