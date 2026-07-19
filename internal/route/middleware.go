package route

import (
	"log/slog"
	"net/http"
)

func jsonBodyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != ContentTypeJSON {
			slog.Error("content type was not set to application/json", slog.String("Content-Type", contentType))
			http.Error(w, "content type is not application/json", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
