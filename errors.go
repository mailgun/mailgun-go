package mailgun

import (
	"errors"
	"fmt"
)

// UnexpectedResponseError this error will be returned whenever a Mailgun API returns an error response.
// Your application can check the Actual field to see the actual HTTP response code returned.
// URL contains the base URL accessed, sans any query parameters.
type UnexpectedResponseError struct {
	Expected []int
	Actual   int
	Method   string
	URL      string
	Data     []byte
}

// String() converts the error into a human-readable, logfmt-compliant string.
// See http://godoc.org/github.com/kr/logfmt for details on logfmt formatting.
func (e *UnexpectedResponseError) String() string {
	return fmt.Sprintf("UnexpectedResponseError Method=%s URL=%s ExpectedOneOf=%#v Got=%d Error: %s",
		e.Method, e.URL, e.Expected, e.Actual, string(e.Data))
}

// Error() performs as String().
func (e *UnexpectedResponseError) Error() string {
	return e.String()
}

// newError creates a new error condition to be returned.
func newError(method, url string, expected []int, got *httpResponse) error {
	return &UnexpectedResponseError{
		Expected: expected,
		Actual:   got.Code,
		Method:   method,
		URL:      url,
		Data:     got.Data,
	}
}

// GetStatusFromErr extracts the http status code from error object
func GetStatusFromErr(err error) int {
	var obj *UnexpectedResponseError
	if errors.As(err, &obj) {
		return obj.Actual
	}

	return -1
}
