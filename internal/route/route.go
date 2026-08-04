package route

import (
	"fmt"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.Handler

type Route struct {
	Handler        http.Handler
	Route          string
	Method         string
	RolesKey       string
	PermissionsKey string
}

func (r *Route) WithMiddleware(m Middleware) *Route {
	r.Handler = m(r.Handler)
	return r
}

func (r *Route) WithPrefix(prefix string) *Route {
	r.Route = fmt.Sprintf("/%v/%v", strings.Trim(prefix, "/"), strings.Trim(r.Route, "/"))

	return r
}

func (r *Route) WithJsonBody() *Route {
	return r.WithMiddleware(jsonBodyMiddleware)
}

func (r *Route) WithRolesKey(key string) *Route {
	r.RolesKey = key
	return r
}

func (r *Route) WithPermissionsKey(key string) *Route {
	r.PermissionsKey = key
	return r
}

func (r *Route) WithPermissions(permissions ...string) *Route {
	return r.WithMiddleware(r.PermissionsMiddleware(permissions))
}

func New(method, route string, handler http.Handler) *Route {
	return &Route{
		Route:   route,
		Method:  method,
		Handler: handler,
	}
}
