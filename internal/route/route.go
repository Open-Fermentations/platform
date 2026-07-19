package route

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

type Route struct {
	Handler http.Handler
	Route   string
	Method  string
}

func (r *Route) WithMiddleware(m Middleware) *Route {
	r.Handler = m(r.Handler)
	return r
}

func (r *Route) WithPrefix(prefix string) *Route {
	r.Route = fmt.Sprintf("/%v/%v", strings.Trim(prefix, "/"), strings.Trim(r.Route, "/"))

	return r
}

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

func (r *Route) WithJsonBody() *Route {
	r.Handler = jsonBodyMiddleware(r.Handler)
	return r
}

func New(method, route string, handler http.Handler) *Route {
	return &Route{
		Route:   route,
		Method:  method,
		Handler: handler,
	}
}
