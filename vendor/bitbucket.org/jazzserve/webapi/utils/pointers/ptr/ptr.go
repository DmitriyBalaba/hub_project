package ptr

import (
	"time"
)

func Float64(p float64) *float64 {
	return &p
}

func Float32(p float32) *float32 {
	return &p
}

func Bool(v bool) *bool {
	return &v
}

func NilBool(v bool) *bool {
	if v == false {
		return nil
	}
	return &v
}

func Str(v string) *string {
	return &v
}

func NilStr(v string) *string {
	if v == "" {
		return nil
	}
	return &v
}

func Int(v int) *int {
	return &v
}

func NilInt(v int) *int {
	if v == 0 {
		return nil
	}
	return &v
}

func Int64(v int64) *int64 {
	return &v
}

func NilInt64(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func Uint64(v uint64) *uint64 {
	return &v
}

func NilUint64(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	return &v
}

func Time(v time.Time) *time.Time {
	return &v
}

func NilTime(v time.Time) *time.Time {
	if v.IsZero() {
		return nil
	}
	return &v
}

func ByteArr(v []byte) *[]byte {
	return &v
}

func NilByteArr(v []byte) *[]byte {
	if len(v) == 0 {
		return nil
	}
	return &v
}

func CopyBool(v *bool) *bool {
	if v == nil {
		return nil
	}
	return Bool(*v)
}

func CopyStr(v *string) *string {
	if v == nil {
		return nil
	}
	return Str(*v)
}

func CopyInt(v *int) *int {
	if v == nil {
		return nil
	}
	return Int(*v)
}

func CopyInt64(v *int64) *int64 {
	if v == nil {
		return nil
	}
	return Int64(*v)
}

func CopyUint64(v *uint64) *uint64 {
	if v == nil {
		return nil
	}
	return Uint64(*v)
}

func CopyTime(v *time.Time) *time.Time {
	if v == nil {
		return nil
	}
	return Time(*v)
}

func CopyByteArr(v *[]byte) *[]byte {
	if v == nil {
		return nil
	}
	return ByteArr(*v)
}
