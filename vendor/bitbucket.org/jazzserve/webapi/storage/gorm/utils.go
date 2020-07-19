package gorm

import (
	"strings"
)

func unifyStr(s string) string {
	return strings.TrimSpace(strings.ToLower(s))
}
