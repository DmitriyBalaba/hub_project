package routing

import (
	"context"
	"github.com/gorilla/mux"
)

// RouteInitializer is a function which inits routes for specified router
type RouteInitializer func(r *mux.Router, ctx context.Context)

var routeInitializers []RouteInitializer

// Adds RouteInitializer to global list for NewRouter
func AddRouteInitializer(f RouteInitializer) {
	if f == nil {
		panic("cannot add nil RouteInitializer")
	}
	routeInitializers = append(routeInitializers, f)
}

func NewRouter(pathPrefix string) *mux.Router {
	r := mux.NewRouter()
	if pathPrefix != "" {
		r = r.PathPrefix(pathPrefix).Subrouter()
	}
	return r
}

// AddRoutesTo consumes all the routes from routeInitializers
func AddRoutesTo(r *mux.Router, ctx context.Context) {
	if r == nil || ctx == nil {
		panic("incomplete input to router.AddRoutesTo")
	}

	if len(routeInitializers) == 0 {
		panic("no routes found")
	}

	for i := range routeInitializers {
		routeInitializers[i](r, ctx)
	}

	// reset routeInitializers
	routeInitializers = nil
}
