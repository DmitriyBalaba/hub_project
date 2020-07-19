package val

import (
	"time"
)

func Float64(p *float64) float64 {
	if p == nil {
		return 0
	}
	return *p
}

func Float32(p *float32) float32 {
	if p == nil {
		return 0
	}
	return *p
}

func Bool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

func Str(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func Int(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func Int64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

func Uint64(p *uint64) uint64 {
	if p == nil {
		return 0
	}
	return *p
}

func Time(p *time.Time) time.Time {
	if p == nil {
		return time.Time{}
	}
	return *p
}

func ByteArr(p *[]byte) []byte {
	if p == nil {
		return nil
	}
	return *p
}
