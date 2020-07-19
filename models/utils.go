package models

import "strings"

func unifyStr(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}
