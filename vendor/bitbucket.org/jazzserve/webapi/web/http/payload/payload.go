package payload

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"strings"
)

const (
	ContentType   = "Content-Type"
	Origin        = "Origin"
	Accept        = "Accept"
	Authorization = "Authorization"

	JsonContentType = "application/json"
)

var (
	validate      = validator.New()
	schemaDecoder = schema.NewDecoder()
)

func init() {
	schemaDecoder.IgnoreUnknownKeys(true)
}

func RegisterValidation(tagKey string, f validator.Func) {
	if err := validate.RegisterValidation(tagKey, f); err != nil {
		panic("can't register validation: " + err.Error())
	}
}

func SetContentType(w http.ResponseWriter, value string) {
	if w == nil {
		panic("cannot SetContentType with nil ResponseWriter")
	}
	w.Header().Set(ContentType, value)
}

func SetDownloadableHeaders(w http.ResponseWriter, contentType, filename string) {
	SetContentType(w, contentType)
	w.Header().Set("Content-disposition", "attachment; filename="+filename)
}

func JsonEncode(w http.ResponseWriter, v interface{}) {
	if w == nil {
		panic("cannot JsonEncode with nil ResponseWriter")
	}
	if v == nil {
		return
	}

	SetContentType(w, JsonContentType)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Error().Msgf("can't encode json: %s", err.Error())
	}
}

func WriteBytes(w http.ResponseWriter, bytes []byte, contentType string) {
	SetContentType(w, contentType)
	if _, err := w.Write(bytes); err != nil {
		log.Error().Msgf("can't write bytes: %s", err.Error())
	}
}

func SetAttachmentHeader(w http.ResponseWriter, filename string) {
	w.Header().Set("Content-disposition", fmt.Sprintf("attachment; filename=%s", filename))
}

func getStrD(value, defaultValue string) string {
	v := strings.TrimSpace(value)
	if v == "" {
		return defaultValue
	}
	return v
}

func getStr(value string) string {
	return getStrD(value, "")
}

func getStrWhiteListed(value string, whiteList ...string) (string, error) {
	if len(whiteList) == 0 {
		return "", NewInternalServerErr("empty white list")
	}

	str := getStr(value)
	if str == "" {
		return whiteList[0], nil
	}

	for i := range whiteList {
		if str == whiteList[i] {
			return str, nil
		}
	}
	return "", NewBadRequestErr("%s does not validate as %s", str, whiteList)
}

func GetQueryStrD(r *http.Request, key, defaultValue string) string {
	return getStrD(r.FormValue(key), defaultValue)
}

func GetQueryStr(r *http.Request, key string) string {
	return GetQueryStrD(r, key, "")
}

func GetQueryStrListed(r *http.Request, key string, whiteList ...string) (string, error) {
	return getStrWhiteListed(r.FormValue(key), whiteList...)
}

func GetPathStr(r *http.Request, key string) string {
	return getStrD(mux.Vars(r)[key], "")
}

func GetPathStrListed(r *http.Request, key string, whiteList ...string) (string, error) {
	return getStrWhiteListed(mux.Vars(r)[key], whiteList...)
}

func getInt64D(key, value string, defaultValue int64) (int64, error) {
	strV := getStr(value)
	if strV == "" {
		return defaultValue, nil
	}

	v, err := strconv.ParseInt(strV, 10, 64)
	if err != nil {
		return 0, NewBadRequestErr("%s query param (%s) is invalid: %s", key, strV, err.Error())
	}

	return v, nil
}

func getInt(key, value string) (int64, error) {
	return getInt64D(key, value, 0)
}

func GetQueryInt64D(r *http.Request, key string, defaultValue int64) (int64, error) {
	return getInt64D(key, r.FormValue(key), defaultValue)
}

func GetQueryInt64(r *http.Request, key string) (int64, error) {
	return GetQueryInt64D(r, key, 0)
}

func GetPathInt64(r *http.Request, key string) (int64, error) {
	return getInt64D(key, mux.Vars(r)[key], 0)
}
