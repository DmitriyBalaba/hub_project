package payload

import (
	"bitbucket.org/jazzserve/webapi/utils"
	"context"
	"encoding/json"
	"errors"
	"github.com/jinzhu/copier"
	"github.com/justinas/alice"
	"io"
	"net/http"
	"reflect"
)

const (
	DtoCtxKey         = "dto-validator-ctx-key"
	ModelCtxKey       = "model-validator-ctx-key"
	QueryParamsCtxKey = "query-params-ctx-key"
)

func getOrPanic(ctx context.Context, key string) (v interface{}) {
	if v = ctx.Value(key); v == nil {
		panic("nothing found in context by " + key)
	}
	return
}

func GetDTO(ctx context.Context) interface{} {
	return getOrPanic(ctx, DtoCtxKey)
}

func GetDTOoptional(ctx context.Context) (v interface{}, ok bool) {
	v = ctx.Value(DtoCtxKey)
	if v != nil {
		ok = true
	}
	return
}

func GetModel(ctx context.Context) interface{} {
	return getOrPanic(ctx, ModelCtxKey)
}

func GetModelOptional(ctx context.Context) (v interface{}, ok bool) {
	v = ctx.Value(ModelCtxKey)
	if v != nil {
		ok = true
	}
	return
}

func getCopyByType(t interface{}) interface{} {
	return reflect.New(reflect.TypeOf(t)).Interface()
}

func withDtoAndModel(r *http.Request, dto, model interface{}) *http.Request {
	ctx := r.Context()
	ctx = context.WithValue(ctx, DtoCtxKey, dto)
	ctx = context.WithValue(ctx, ModelCtxKey, model)
	return r.WithContext(ctx)
}

func JsonBody(dto, model interface{}) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			dto, model := getCopyByType(dto), getCopyByType(model)

			if err := getValid(r.Body, dto); err != nil {
				ErrorResponse(w, WrapBadRequestErr(err))
				return
			}

			if err := copyStructFields(dto, model); err != nil {
				ErrorResponse(w, err)
				return
			}

			next.ServeHTTP(w, withDtoAndModel(r, dto, model))
		})
	}
}

func JsonBodyOptional(dto, model interface{}) alice.Constructor {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			dto, model := getCopyByType(dto), getCopyByType(model)

			err := getValid(r.Body, dto)

			if errors.Is(err, io.EOF) {
				next.ServeHTTP(w, withDtoAndModel(r, dto, model))
				return
			}

			if err != nil {
				ErrorResponse(w, WrapBadRequestErr(err))
				return
			}

			if err := copyStructFields(dto, model); err != nil {
				ErrorResponse(w, err)
				return
			}

			next.ServeHTTP(w, withDtoAndModel(r, dto, model))
		})
	}
}

func getValid(body io.ReadCloser, obj interface{}) (err error) {
	decoder := json.NewDecoder(body)
	defer utils.CheckCloseError(body, &err)

	if err = decoder.Decode(obj); err != nil {
		return
	}

	if err = validate.Struct(obj); err != nil {
		return
	}

	return
}

func getReflectValue(ptr interface{}) (v reflect.Value) {
	v = reflect.ValueOf(ptr)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	return v
}

func copyStructFields(src interface{}, dst interface{}) error {
	return copier.Copy(dst, src)
}
