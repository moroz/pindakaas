package dbtypes

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"time"
)

type UnixTimestamp struct {
	*time.Time
}

var (
	_ driver.Valuer = &UnixTimestamp{}
	_ sql.Scanner   = &UnixTimestamp{}
)

func (u UnixTimestamp) Value() (driver.Value, error) {
	return u.Unix(), nil
}

func (u *UnixTimestamp) Scan(src any) error {
	if src == nil {
		return nil
	}

	val, ok := src.(int64)
	if !ok {
		return errors.New("failed to decode value as int64")
	}

	value := time.Unix(val, 0)
	u.Time = &value
	return nil
}
