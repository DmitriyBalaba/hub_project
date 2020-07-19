package payload

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
)

type HttpError struct {
	error
	status int
}

func (e HttpError) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Message string `json:"message"`
	}{
		Message: e.Error(),
	})
}

func getHttpError(err error) (httpErr *HttpError) {
	httpErr, ok := err.(*HttpError)
	if !ok {
		return WrapInternalServerErr(err)
	}
	return httpErr
}

func ErrorResponse(w http.ResponseWriter, err error) {
	if w == nil {
		panic("ErrorResponse: nil ResponseWriter")
	}
	if err == nil {
		return
	}

	httpErr := getHttpError(err)
	w.WriteHeader(httpErr.status)
	JsonEncode(w, httpErr)
	if httpErr.status == http.StatusInternalServerError {
		log.Error().Msgf("Responded with 500: %+v", err)
	}
}

func newErrorf(status int, format string, args ...interface{}) *HttpError {
	return &HttpError{
		error:  errors.Errorf(format, args...),
		status: status,
	}
}

func wrapError(status int, err error) *HttpError {
	if err == nil {
		return nil
	}
	if httpErr, ok := err.(*HttpError); ok {
		httpErr.status = status
		return httpErr
	}
	return &HttpError{
		error:  err,
		status: status,
	}
}

func IsHttpErr(err error) bool {
	_, ok := err.(*HttpError)
	return ok
}

func NewInternalServerErr(format string, args ...interface{}) *HttpError {
	return newErrorf(http.StatusInternalServerError, format, args...)
}

func WrapInternalServerErr(err error) *HttpError {
	return wrapError(http.StatusInternalServerError, err)
}

func NewBadRequestErr(format string, args ...interface{}) *HttpError {
	return newErrorf(http.StatusBadRequest, format, args...)
}

func WrapBadRequestErr(err error) *HttpError {
	return wrapError(http.StatusBadRequest, err)
}

func NewUnauthorizedErr(format string, args ...interface{}) *HttpError {
	return newErrorf(http.StatusUnauthorized, format, args...)
}

func WrapUnauthorizedErr(err error) *HttpError {
	return wrapError(http.StatusUnauthorized, err)
}

func NewForbiddenErr(format string, args ...interface{}) *HttpError {
	return newErrorf(http.StatusForbidden, format, args...)
}

func WrapForbiddenErr(err error) *HttpError {
	return wrapError(http.StatusForbidden, err)
}

func NewNotFoundErr(format string, args ...interface{}) *HttpError {
	return newErrorf(http.StatusNotFound, format, args...)
}

func WrapNotFoundErr(err error) *HttpError {
	return wrapError(http.StatusNotFound, err)
}

func NewUnprocessableEntityErr(format string, args ...interface{}) *HttpError {
	return newErrorf(http.StatusUnprocessableEntity, format, args...)
}

func WrapUnprocessableEntityErr(err error) *HttpError {
	return wrapError(http.StatusUnprocessableEntity, err)
}
