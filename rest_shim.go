package mailgun

import (
	"context"
	"fmt"
	"net/http"
)

// The UserAgent identifies the client to the server, for logging purposes.
// In the event of problems requiring a human administrator's assistance,
// this user agent allows them to identify the client from human-generated activity.
const UserAgent = "mailgun-go/" + Version

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

// notGood searches a list of response codes (the haystack) for a matching entry (the needle).
// If found, the response code is considered good, and thus false is returned.
// Otherwise true is returned.
func notGood(needle int, haystack []int) bool {
	for _, i := range haystack {
		if needle == i {
			return false
		}
	}
	return true
}

// expected denotes the expected list of known-good HTTP response codes possible from the Mailgun API.
var expected = []int{200, 202, 204}

// makeRequest shim performs a generic request, checking for a positive outcome.
// See simplehttp.MakeRequest for more details.
func makeRequest(ctx context.Context, r *httpRequest, method string, p payload) (*httpResponse, error) {
	r.addHeader("User-Agent", UserAgent)
	rsp, err := r.makeRequest(ctx, method, p)
	if (err == nil) && notGood(rsp.Code, expected) {
		return rsp, newError(method, r.URL, expected, rsp)
	}
	return rsp, err
}

// getResponseFromJSON shim performs a GET request, checking for a positive outcome.
// See simplehttp.GetResponseFromJSON for more details.
func getResponseFromJSON(ctx context.Context, r *httpRequest, v any) error {
	r.addHeader("User-Agent", UserAgent)
	response, err := r.makeGetRequest(ctx)
	if err != nil {
		return err
	}
	if notGood(response.Code, expected) {
		return newError(http.MethodGet, r.URL, expected, response)
	}
	return response.parseFromJSON(v)
}

// postResponseFromJSON shim performs a POST request, checking for a positive outcome.
// See simplehttp.PostResponseFromJSON for more details.
func postResponseFromJSON(ctx context.Context, r *httpRequest, p payload, v any) error {
	r.addHeader("User-Agent", UserAgent)
	response, err := r.makePostRequest(ctx, p)
	if err != nil {
		return err
	}
	if notGood(response.Code, expected) {
		return newError(http.MethodPost, r.URL, expected, response)
	}
	return response.parseFromJSON(v)
}

// putResponseFromJSON shim performs a PUT request, checking for a positive outcome.
// See simplehttp.PutResponseFromJSON for more details.
func putResponseFromJSON(ctx context.Context, r *httpRequest, p payload, v any) error {
	r.addHeader("User-Agent", UserAgent)
	response, err := r.makePutRequest(ctx, p)
	if err != nil {
		return err
	}
	if notGood(response.Code, expected) {
		return newError(http.MethodPut, r.URL, expected, response)
	}
	return response.parseFromJSON(v)
}

// makeGetRequest shim performs a GET request, checking for a positive outcome.
// See simplehttp.MakeGetRequest for more details.
func makeGetRequest(ctx context.Context, r *httpRequest) (*httpResponse, error) {
	r.addHeader("User-Agent", UserAgent)
	rsp, err := r.makeGetRequest(ctx)
	if (err == nil) && notGood(rsp.Code, expected) {
		return rsp, newError(http.MethodGet, r.URL, expected, rsp)
	}
	return rsp, err
}

// makePostRequest shim performs a POST request, checking for a positive outcome.
// See simplehttp.MakePostRequest for more details.
func makePostRequest(ctx context.Context, r *httpRequest, p payload) (*httpResponse, error) {
	r.addHeader("User-Agent", UserAgent)
	rsp, err := r.makePostRequest(ctx, p)
	if (err == nil) && notGood(rsp.Code, expected) {
		return rsp, newError(http.MethodPost, r.URL, expected, rsp)
	}
	return rsp, err
}

// makePutRequest shim performs a PUT request, checking for a positive outcome.
// See simplehttp.MakePutRequest for more details.
func makePutRequest(ctx context.Context, r *httpRequest, p payload) (*httpResponse, error) {
	r.addHeader("User-Agent", UserAgent)
	rsp, err := r.makePutRequest(ctx, p)
	if (err == nil) && notGood(rsp.Code, expected) {
		return rsp, newError(http.MethodPut, r.URL, expected, rsp)
	}
	return rsp, err
}

// makeDeleteRequest shim performs a DELETE request, checking for a positive outcome.
// See simplehttp.MakeDeleteRequest for more details.
func makeDeleteRequest(ctx context.Context, r *httpRequest) (*httpResponse, error) {
	r.addHeader("User-Agent", UserAgent)
	rsp, err := r.makeDeleteRequest(ctx)
	if (err == nil) && notGood(rsp.Code, expected) {
		return rsp, newError(http.MethodDelete, r.URL, expected, rsp)
	}
	return rsp, err
}

// Extract the http status code from error object
func GetStatusFromErr(err error) int {
	obj, ok := err.(*UnexpectedResponseError)
	if !ok {
		return -1
	}
	return obj.Actual
}
