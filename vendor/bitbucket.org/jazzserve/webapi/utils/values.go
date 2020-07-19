package utils

import (
	"reflect"
	"strings"
)

func IsStrIn(s string, arr ...string) bool {
	for i := range arr {
		if s == arr[i] {
			return true
		}
	}
	return false
}

func IsZeroAny(v ...interface{}) bool {
	for i := range v {
		if reflect.ValueOf(v[i]).IsZero() {
			return true
		}
	}
	return false
}

func IsZeroAnyInt64(i ...int64) bool {
	for j := range i {
		if i[j] == 0 {
			return true
		}
	}
	return false
}

func IsEmptyStrAny(s ...string) bool {
	for i := range s {
		if strings.TrimSpace(s[i]) == "" {
			return true
		}
	}
	return false
}

func IsArray(v ...interface{}) bool {
	for i := range v {
		kind := reflect.ValueOf(v[i]).Kind()
		if kind != reflect.Slice && kind != reflect.Array {
			return false
		}
	}
	return true
}
