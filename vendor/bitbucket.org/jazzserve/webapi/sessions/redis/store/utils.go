package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"strings"
)

const (
	// redis command list
	Select  = "SELECT"
	Del     = "DEL"
	HGetAll = "HGETALL"
	HMSet   = "HMSET"
	Exists  = "EXISTS"
	Expire  = "EXPIRE"
	Scan    = "SCAN"
	Match   = "MATCH"

	randSessIDSize = 32

	UserIDKey = "user-id-for-session"

	sessionIDSeparator = ":"
)

func randomize(size int) string {
	b := make([]byte, size)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func GetRandomSessionID(r *http.Request) (string, error) {
	return randomize(randSessIDSize), nil
}

func PutUserID(r *http.Request, userID string) (*http.Request, error) {
	if r == nil {
		return nil, errors.New("can't proceed with nil request")
	}
	if userID == "" {
		return nil, errors.New("user id is empty")
	}
	if strings.Contains(userID, sessionIDSeparator) {
		return nil, errors.Errorf("user id contains separator: %s", sessionIDSeparator)
	}
	r = r.WithContext(context.WithValue(r.Context(), UserIDKey, userID))
	return r, nil
}

func GetSessionIDWithUserID(r *http.Request) (sessionID string, err error) {
	if r == nil {
		err = errors.New("can't proceed with nil request")
		return
	}
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		err = errors.New("user id is empty")
		return
	}
	sessionID = fmt.Sprintf("%s%s%s", userID, sessionIDSeparator, randomize(randSessIDSize))
	return
}

func mergeIntoStrMap(main, secondary map[string]string) {
	for k := range secondary {
		if _, ok := main[k]; !ok {
			main[k] = secondary[k]
		}
	}
}

func getTypeAndValue(obj interface{}) (t reflect.Type, v reflect.Value) {
	t, v = reflect.TypeOf(obj), reflect.ValueOf(obj)
	for t.Kind() == reflect.Ptr {
		t, v = t.Elem(), v.Elem()
	}
	return
}
