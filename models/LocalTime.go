package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type LocalTime struct { // 内嵌方式（推荐）
	time.Time
}

func (t LocalTime) MarshalJSON() ([]byte, error) {
	// tune := fmt.Sprintf(`"%s"`, t.Format("2006-01-02 15:04:05"))
	tune := t.Local().Format(`"2006-01-02 15:04:05"`)
	return []byte(tune), nil
}

// Value insert timestamp into mysql need this function.
func (t LocalTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	if t.Time.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return t.Time, nil
}

// Scan valueof time.Time
func (t *LocalTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = LocalTime{Time: value}
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}

func (t LocalTime) After(u LocalTime) bool {
	return t.Time.After(u.Time)
}
func (t LocalTime) Before(u LocalTime) bool {
	return t.Time.Before(u.Time)
}
func (t LocalTime) Equal(u LocalTime) bool {
	return t.Time.Equal(u.Time)
}
