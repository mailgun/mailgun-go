package mtypes

import (
	"strconv"
	"time"
)

// RFC3339 Mailgun uses RFC3339 format for timestamps in ip warmups endpoints ('2025-01-01T00:00:00Z'), but
// by default Go's JSON package uses another format when decoding/encoding timestamps.
// TODO(v6): make a struct and embed time.Time to inherit all its methods.
type RFC3339Time time.Time

func NewRFC3339Time(str string) (RFC3339Time, error) {
	t, err := time.Parse(time.RFC3339, str)
	if err != nil {
		return RFC3339Time{}, err
	}
	return RFC3339Time(t), nil
}

func (t RFC3339Time) Unix() int64 {
	return time.Time(t).Unix()
}

func (t RFC3339Time) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t RFC3339Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Time(t).Format(time.RFC3339))), nil
}

func (t *RFC3339Time) UnmarshalJSON(s []byte) error {
	q, err := strconv.Unquote(string(s))
	if err != nil {
		return err
	}

	*(*time.Time)(t), err = time.Parse(time.RFC3339, q)
	if err != nil {
		return err
	}

	return nil
}

func (t RFC3339Time) String() string {
	return time.Time(t).Format(time.RFC3339)
}
