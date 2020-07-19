package routing

import (
	"bitbucket.org/jazzserve/webapi/web/http/payload"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/pkg/errors"
	"net/http"
)

// FallibleHandlerFunc is a wrapper of http.HandlerFunc for detecting handler errors
type FallibleHandlerFunc func(w http.ResponseWriter, r *http.Request) (err error)

// Exec is a func which responds in correspondence with error type
func Exec(f FallibleHandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				payload.ErrorResponse(w, errors.Errorf("panic: %v", r))
			}
		}()
		if err := f(w, r); err != nil {
			payload.ErrorResponse(w, err)
		}
	})
}

func New(r *mux.Router, path string, h FallibleHandlerFunc, constructors ...alice.Constructor) *mux.Route {
	if r == nil || h == nil {
		panic("cannot create new route with nil input")
	}
	return r.Handle(path, alice.New(constructors...).Then(Exec(h)))
}
