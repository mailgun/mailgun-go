package mtypes

import (
	"strconv"
	"time"

	"github.com/mailgun/errors"
)

// ISO8601Time Mailgun uses ISO8601 format for timestamps in the API key API endpoints ('2011/10/13T18:02:00'), but
// by default Go's JSON package uses another format when decoding/encoding timestamps. This goes against the documentation
// at https://documentation.mailgun.com/docs/mailgun/api-reference#date-format but here we are.
type ISO8601Time time.Time

const (
	ISO8601Format = "2006-01-02T15:04:05"
)

func NewISO8601Time(str string) (ISO8601Time, error) {
	t, err := time.Parse(ISO8601Format, str)
	if err != nil {
		return ISO8601Time{}, err
	}
	return ISO8601Time(t), nil
}

func (t ISO8601Time) Unix() int64 {
	return time.Time(t).Unix()
}

func (t ISO8601Time) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t ISO8601Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Time(t).Format(ISO8601Format))), nil
}

func (t *ISO8601Time) UnmarshalJSON(s []byte) error {
	q, err := strconv.Unquote(string(s))
	if err != nil {
		return err
	}

	if q == "" {
		return nil
	}

	var err1 error
	*(*time.Time)(t), err1 = time.Parse(ISO8601Format, q)
	if err1 != nil {
		var err2 error
		*(*time.Time)(t), err2 = time.Parse(ISO8601Format, q)
		if err2 != nil {
			// TODO(go1.20): use errors.Join:
			return errors.Errorf("%s; %s", err1, err2)
		}
	}

	return nil
}

func (t ISO8601Time) String() string {
	return time.Time(t).Format(ISO8601Format)
}
