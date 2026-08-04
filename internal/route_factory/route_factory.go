package routefactory

import (
	"net/http"
	"open-fermentations/internal/route"
)

type routeFactory struct {
	PermissionsKey string
	RolesKey       string
}

// NewRoute implements [RouteFactory].
func (r routeFactory) New(method, path string, handler http.Handler) *route.Route {
	return route.New(method, path, handler).
		WithRolesKey(r.RolesKey).
		WithPermissionsKey(r.PermissionsKey)
}

var _ RouteFactory = routeFactory{}

type RouteFactory interface {
	New(method, route string, handler http.Handler) *route.Route
}

func New(permissionsKey, rolesKey string) *routeFactory {
	return &routeFactory{
		PermissionsKey: permissionsKey,
		RolesKey:       rolesKey,
	}
}
