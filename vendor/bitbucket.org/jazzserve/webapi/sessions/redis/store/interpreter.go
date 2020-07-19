package store

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"strings"
)

type Interpreter interface {
	// interpreting values from / to redis
	ConvertObjFromRedis(redisValues map[string]string, objType interface{}) (obj interface{}, err error)
	ConvertObjToRedis(obj interface{}) (redisValues map[string]string, err error)
	// adding / getting a type in which obj should be returned from redis
	PutObjType(*http.Request, interface{}) (*http.Request, error)
	GetObjType(*http.Request) interface{}
	// adding / getting an object from session
	PutObj(*sessions.Session, interface{}) error
	GetObj(*sessions.Session) interface{}
}

var _ Interpreter = &jsonInterpreter{}

type jsonInterpreter struct{}

// these operations (Convert From / To Redis) can be optimised by manually mapping every key with marshalled value
func (ji *jsonInterpreter) ConvertObjFromRedis(redisValues map[string]string, objType interface{}) (obj interface{}, err error) {
	t := reflect.TypeOf(objType)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	obj = reflect.New(t).Interface()

	fieldsRaw := make(map[string]json.RawMessage)
	for key := range redisValues {
		fieldsRaw[key] = json.RawMessage(redisValues[key])
	}
	fullRaw, err := json.Marshal(fieldsRaw)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	err = json.Unmarshal(fullRaw, obj)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	return
}

func (ji *jsonInterpreter) structToRedis(sValue reflect.Value, sType reflect.Type) (map[string]string, error) {
	if k := sValue.Kind(); k != reflect.Struct {
		panic(fmt.Sprintf("struct to redis works with structures only; have %v", k))
	}

	redisValues := make(map[string]string)

	for i := 0; i < sValue.NumField(); i++ {
		fieldVal, structField := sValue.Field(i), sType.Field(i)

		// gets a name from json tag
		// (it's safe to get the first element if sep is not empty)
		name := strings.Split(structField.Tag.Get("json"), ",")[0]
		hasJsonTag := true
		if name == "" {
			name = structField.Name
			hasJsonTag = false
		}

		if !fieldVal.IsValid() || !fieldVal.CanInterface() ||
			fieldVal.IsZero() || (hasJsonTag && name == "-") {
			continue
		}

		// anonymous fields should be marshalled at the same level
		if !hasJsonTag && structField.Anonymous {
			partOfRedisValues, err := ji.convertObjAsStructF(fieldVal.Interface())
			if err != nil {
				return nil, err
			}
			mergeIntoStrMap(redisValues, partOfRedisValues)
			continue
		}

		marshalledFieldVal, err := json.Marshal(fieldVal.Interface())
		if err != nil {
			return nil, errors.WithStack(err)
		}

		redisValues[name] = string(marshalledFieldVal)
	}

	return redisValues, nil
}

func (ji *jsonInterpreter) mapToRedis(mValue reflect.Value, mType reflect.Type) (map[string]string, error) {
	if k := mValue.Kind(); k != reflect.Map {
		panic(fmt.Sprintf("struct to redis works with structures only; have %v", k))
	}

	jsonVal, err := json.Marshal(mValue.Interface())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rawMap := make(map[string]json.RawMessage)
	if err = json.Unmarshal(jsonVal, &rawMap); err != nil {
		return nil, errors.WithStack(err)
	}

	redisValues := make(map[string]string)
	for k := range rawMap {
		redisValues[k] = string(rawMap[k])
	}

	return redisValues, nil
}

func (ji *jsonInterpreter) anyTypeToRedis(sValue reflect.Value, sType reflect.Type) (map[string]string, error) {
	if !sValue.IsValid() || !sValue.CanInterface() || sValue.IsZero() {
		return map[string]string{}, nil
	}

	marshalledFieldVal, err := json.Marshal(sValue.Interface())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	redisValues := map[string]string{
		sType.Name(): string(marshalledFieldVal),
	}

	return redisValues, nil
}

func (ji *jsonInterpreter) convertObjAsStructF(obj interface{}) (map[string]string, error) {
	objType, objVal := getTypeAndValue(obj)

	switch objVal.Kind() {
	case reflect.Struct:
		return ji.structToRedis(objVal, objType)
	default:
		return ji.anyTypeToRedis(objVal, objType)
	}
}

func (ji *jsonInterpreter) ConvertObjToRedis(obj interface{}) (map[string]string, error) {
	objType, objVal := getTypeAndValue(obj)

	switch objVal.Kind() {
	case reflect.Struct:
		return ji.structToRedis(objVal, objType)
	case reflect.Map:
		return ji.mapToRedis(objVal, objType)
	default:
		return nil, errors.Errorf("can't convert %v to redis values", objType)
	}
}

const (
	objCtxKey     = "redis-store-obj"
	objTypeCtxKey = "redis-store-obj-type"
)

func (ji *jsonInterpreter) PutObjType(r *http.Request, objType interface{}) (*http.Request, error) {
	if r == nil {
		return nil, errors.New("can't proceed: nil request")
	}
	if objType == nil {
		return nil, errors.New("can't proceed: nil value type")
	}
	// maybe, we should validate objType here
	return r.WithContext(context.WithValue(r.Context(), objTypeCtxKey, objType)), nil
}

func (ji *jsonInterpreter) GetObjType(r *http.Request) interface{} {
	t := r.Context().Value(objTypeCtxKey)
	if t == nil {
		return &map[string]json.RawMessage{}
	}
	return t
}

func (ji *jsonInterpreter) PutObj(session *sessions.Session, obj interface{}) error {
	if session == nil || session.Values == nil {
		return errors.New("can't proceed: session is incomplete")
	}
	if obj == nil {
		return nil
	}
	session.Values[objCtxKey] = obj
	return nil
}

func (ji *jsonInterpreter) GetObj(session *sessions.Session) interface{} {
	if session == nil {
		return nil
	}
	return session.Values[objCtxKey]
}
