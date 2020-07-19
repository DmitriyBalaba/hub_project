package env

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

const (
	CtxKey = "environment"
)

var (
	NilConfigErr        = errors.New("no config")
	NotBuiltProperlyErr = errors.New("operation cannot be performed due to env build")
)

func AllowedHeaders() []string {
	return []string{"origin", "content-type", "accept", "authorization"}
}

func AllowedMethods() []string {
	return []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions}
}

func Add(ctx context.Context, e Environment) context.Context {
	return context.WithValue(ctx, CtxKey, e)
}

func Get(ctx context.Context) Environment {
	e, ok := ctx.Value(CtxKey).(Environment)
	if !ok {
		panic("ctx environment is invalid or not found")
	}
	return e
}
