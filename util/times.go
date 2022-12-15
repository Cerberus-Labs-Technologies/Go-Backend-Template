package util

import (
	"database/sql/driver"
	"strings"
	"time"
)

type DateTime struct {
	time.Time
}

func (t DateTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.Format("2006-01-02") + "\""), nil
}

func (t *DateTime) UnmarshalJSON(data []byte) error {
	a, err := time.Parse("2006-01-02", strings.Trim(string(data), "\""))
	t.Time = a
	return err
}

func (dateT *DateTime) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		dateT.Time = t
	}
	return nil
}

func (dateT DateTime) Value() (driver.Value, error) {
	return dateT.Time, nil
}

type TimeStamp struct {
	time.Time
}

func (t TimeStamp) MarshalJSON() ([]byte, error) {
	return []byte("\"" + t.Format("2006-01-02 15:04:05") + "\""), nil
}

func (t *TimeStamp) UnmarshalJSON(data []byte) error {
	a, err := time.Parse("2006-01-02 15:04:05", strings.Trim(string(data), "\""))
	t.Time = a
	return err
}

func (dateT *TimeStamp) Scan(src interface{}) error {
	if t, ok := src.(time.Time); ok {
		dateT.Time = t
	}
	return nil
}

func (dateT TimeStamp) Value() (driver.Value, error) {
	return dateT.Time, nil
}
