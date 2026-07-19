package route

import (
	"log/slog"
	"net/http"
	"time"
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

type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.statusCode = code
	s.ResponseWriter.WriteHeader(code)
}

func (s *statusRecorder) Write(b []byte) (int, error) {
	if s.statusCode == 0 {
		s.WriteHeader(http.StatusOK)
	}

	return s.ResponseWriter.Write(b)
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecorder{ResponseWriter: w, statusCode: 0}
		next.ServeHTTP(rec, r)

		slog.Info("Request", slog.Group("request",
			slog.String("method", r.Method),
			slog.String("url", r.RequestURI),
			slog.Int("status_code", rec.statusCode),
			slog.Duration("duration", time.Since(start)),
		))
	})
}
