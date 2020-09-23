package expensify

import (
	"bytes"
	"fmt"
	"time"
)

const (
	layoutISO = "2006-01-02"
)

// Time which marshals to and unmarshals from "yyyy-mm-dd" format.
type Time struct {
	time.Time
}

// NewTime returns a new Time from the given time.Time.
func NewTime(t time.Time) Time {
	return Time{t}
}

// MarshalJSON implements json.Marshaler.
func (t Time) MarshalJSON() ([]byte, error) {
	if t.IsZero() {
		return []byte("null"), nil
	}
	stamp := fmt.Sprintf("%q", t.Format(layoutISO))
	return []byte(stamp), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (t *Time) UnmarshalJSON(b []byte) (err error) {
	b = bytes.Trim(b, "\"")
	if bytes.Equal(b, []byte("null")) {
		return nil
	} else if t.Time, err = time.Parse(layoutISO, string(b)); err != nil {
		return err
	}
	return nil
}
