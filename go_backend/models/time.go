package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Time struct {
	time.Time
}

func Now() Time {
	return Time{Time: time.Now()}
}

func (t *Time) Scan(src interface{}) error {
	if src == nil {
		t.Time = time.Time{}
		return nil
	}
	switch v := src.(type) {
	case string:
		return t.parse(v)
	case []byte:
		return t.parse(string(v))
	case time.Time:
		t.Time = v
		return nil
	default:
		return fmt.Errorf("models.Time: tipo no soportado %T", src)
	}
}

func (t Time) Value() (driver.Value, error) {
	if t.IsZero() {
		return nil, nil
	}
	return t.Format(time.RFC3339), nil
}

func (t *Time) parse(s string) error {
	formats := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05Z",
	}
	for _, f := range formats {
		if parsed, err := time.Parse(f, s); err == nil {
			t.Time = parsed
			return nil
		}
	}
	return fmt.Errorf("models.Time: no se pudo parsear %q", s)
}
