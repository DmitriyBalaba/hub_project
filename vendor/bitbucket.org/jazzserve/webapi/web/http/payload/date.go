package payload

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/pkg/errors"
	"time"
)

type Date struct {
	Time time.Time
}

const (
	DateFormat = "2006-01-02"
	PSQLFormat = "2006-01-02"
)

func NewDate(t *time.Time) *Date {
	if t == nil {
		return nil
	}
	return &Date{
		Time: time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location()),
	}
}

func (d *Date) GetTime() *time.Time {
	if d == nil {
		return nil
	}
	return &d.Time
}

func (d *Date) MarshalJSON() (data []byte, err error) {
	return json.Marshal(d.Time.Format(DateFormat))
}

func (d *Date) UnmarshalJSON(data []byte) (err error) {
	if string(data) == "null" {
		return nil
	}

	t, err := time.Parse("\""+DateFormat+"\"", string(data))
	if err != nil {
		return errors.WithStack(err)
	}

	*d = Date{t}
	return
}

func (d *Date) Scan(value interface{}) error {
	var ok bool
	d.Time, ok = value.(time.Time)
	if !ok {
		return errors.New("can't scan date: not time.Time")
	}
	return nil
}

func (d *Date) Value() (driver.Value, error) {
	if d == nil {
		return nil, nil
	}
	return d.Time, nil
}

func (d *Date) PSQLString() string {
	return d.GetTime().Format(PSQLFormat)
}
