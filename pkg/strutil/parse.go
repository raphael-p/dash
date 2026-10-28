package strutil

import (
	"cmp"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type timeVal interface {
	time.Time | sql.NullTime
}

// FormatTime formats a timestamp into a standard datetime format
//
// supply args[0] to override the layout, args[1] to provide a fallback
func FormatTime[T timeVal](val T, args ...string) string {
	layout := time.DateTime
	fallback := ""

	if len(args) > 0 {
		layout = cmp.Or(args[0], time.DateTime)
	}
	if len(args) > 1 {
		fallback = args[1]
	}
	switch v := any(val).(type) {
	case time.Time:
		if v.IsZero() {
			return fallback
		}
		return v.Format(layout)
	case sql.NullTime:
		if !v.Valid {
			return fallback
		}
		return v.Time.Format(layout)
	default:
		return fallback
	}
}

// OrWarn returns the string if it is non-empty, warns otherwise
func OrWarn(s, label string, warn func(err string)) string {
	if value := strings.TrimSpace(s); value != "" {
		return value
	}
	warn(fmt.Sprintf("please provide a %s", label))
	return ""
}

func parseID(s string) (int64, error) {
	if s == "" {
		return 0, errors.New("ID is empty")
	}

	id, err := strconv.ParseInt(s, 10, 0)
	if err != nil || id < 1 {
		return 0, errors.New("ID is not an integer")
	}

	if id < 1 {
		return 0, errors.New("ID must be greater than zero")
	}
	return id, nil
}

// ParseIDOrWarn converts a string into a valid incremental ID
//
// warn() is invoked if this fails
func ParseIDOrWarn(s string, warn func(err string)) (int64, bool) {
	id, err := parseID(s)
	if err != nil {
		warn(fmt.Sprint("your input is invalid: ", err))
		return 0, false
	}
	return id, true
}

// RelativeTime returns how long ago the timestamp was in a
// concise format
//
// e.g. "1d", "2w", "3y"
func RelativeTime(t time.Time) string {
	d := time.Since(t)
	day := 24 * time.Hour
	switch {
	case d < time.Minute:
		return "now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < day:
		return fmt.Sprintf("%dh", int(d.Hours()))
	case d < 7*day:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	case d < 365*day:
		return fmt.Sprintf("%dw", int(d.Hours()/(24*7)))
	default:
		return fmt.Sprintf("%dy", int(d.Hours()/(24*365)))
	}
}
