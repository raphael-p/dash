package sqlutil

import (
	"database/sql"
	"database/sql/driver"
	"time"
)

func NowNullTime() sql.NullTime {
	return sql.NullTime{Time: time.Now(), Valid: true}
}

type nullableType interface {
	driver.Valuer
}

func OrNil[T nullableType](v T) any {
	val, err := v.Value()
	if err != nil || val == nil {
		return nil
	}
	return val
}
