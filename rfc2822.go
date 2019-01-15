package mailgun

import (
	"strconv"
	"strings"
	"time"
)

// Mailgun uses RFC2822 format for timestamps everywhere ('Thu, 13 Oct 2011 18:02:00 GMT'), but
// by default Go's JSON package uses another format when decoding/encoding timestamps.
type RFC2822Time time.Time

func NewRFC2822Time(str string) (RFC2822Time, error) {
	t, err := time.Parse(time.RFC1123, str)
	if err != nil {
		return RFC2822Time{}, err
	}
	return RFC2822Time(t), nil
}

func (t RFC2822Time) Unix() int64 {
	return time.Time(t).Unix()
}

func (t RFC2822Time) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t RFC2822Time) MarshalJSON() ([]byte, error) {
	return []byte(strconv.Quote(time.Time(t).Format(time.RFC1123))), nil
}

func (t *RFC2822Time) UnmarshalJSON(s []byte) error {
	q, err := strconv.Unquote(string(s))
	if err != nil {
		return err
	}
	if *(*time.Time)(t), err = time.Parse(time.RFC1123, q); err != nil {
		if strings.Contains(err.Error(), "extra text") {
			if *(*time.Time)(t), err = time.Parse(time.RFC1123Z, q); err != nil {
				return err
			}
			return nil
		}
		return err
	}
	return nil
}

func (t RFC2822Time) String() string {
	return time.Time(t).Format(time.RFC1123)
}
