package mailgun

import (
	"context"
	"strconv"
)

// Bounce aggregates data relating to undeliverable messages to a specific intended recipient,
// identified by Address.
type Bounce struct {
	// The time at which Mailgun detected the bounce.
	CreatedAt RFC2822Time `json:"created_at"`
	// Code provides the SMTP error code that caused the bounce
	Code string `json:"code"`
	// Address the bounce is for
	Address string `json:"address"`
	// human readable reason why
	Error string `json:"error"`
}

type Paging struct {
	First    string `json:"first,omitempty"`
	Next     string `json:"next,omitempty"`
	Previous string `json:"previous,omitempty"`
	Last     string `json:"last,omitempty"`
}

type bouncesListResponse struct {
	Items  []Bounce `json:"items"`
	Paging Paging   `json:"paging"`
}

// ListBounces returns a complete set of bounces logged against the sender's domain, if any.
// The results include the total number of bounces (regardless of skip or limit settings),
// and the slice of bounces specified, if successful.
// Note that the length of the slice may be smaller than the total number of bounces.
func (mg *MailgunImpl) ListBounces(opts *ListOptions) *BouncesIterator {
	r := newHTTPRequest(generateApiUrl(mg, bouncesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	if opts != nil {
		if opts.Limit != 0 {
			r.addParameter("limit", strconv.Itoa(opts.Limit))
		}
	}
	url, err := r.generateUrlWithParameters()
	return &BouncesIterator{
		mg:                  mg,
		bouncesListResponse: bouncesListResponse{Paging: Paging{Next: url, First: url}},
		err:                 err,
	}
}

type BouncesIterator struct {
	bouncesListResponse
	mg  Mailgun
	err error
}

// If an error occurred during iteration `Err()` will return non nil
func (ci *BouncesIterator) Err() error {
	return ci.err
}

// Next retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ci *BouncesIterator) Next(ctx context.Context, items *[]Bounce) bool {
	if ci.err != nil {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.Next)
	if ci.err != nil {
		return false
	}
	cpy := make([]Bounce, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	if len(ci.Items) == 0 {
		return false
	}
	return true
}

// First retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (ci *BouncesIterator) First(ctx context.Context, items *[]Bounce) bool {
	if ci.err != nil {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.First)
	if ci.err != nil {
		return false
	}
	cpy := make([]Bounce, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ci *BouncesIterator) Last(ctx context.Context, items *[]Bounce) bool {
	if ci.err != nil {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.Last)
	if ci.err != nil {
		return false
	}
	cpy := make([]Bounce, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ci *BouncesIterator) Previous(ctx context.Context, items *[]Bounce) bool {
	if ci.err != nil {
		return false
	}
	if ci.Paging.Previous == "" {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.Previous)
	if ci.err != nil {
		return false
	}
	cpy := make([]Bounce, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	if len(ci.Items) == 0 {
		return false
	}
	return true
}

func (ci *BouncesIterator) fetch(ctx context.Context, url string) error {
	ci.Items = nil
	r := newHTTPRequest(url)
	r.setClient(ci.mg.Client())
	r.setBasicAuth(basicAuthUser, ci.mg.APIKey())

	return getResponseFromJSON(ctx, r, &ci.bouncesListResponse)
}

// GetBounce retrieves a single bounce record, if any exist, for the given recipient address.
func (mg *MailgunImpl) GetBounce(ctx context.Context, address string) (Bounce, error) {
	r := newHTTPRequest(generateApiUrl(mg, bouncesEndpoint) + "/" + address)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var response Bounce
	err := getResponseFromJSON(ctx, r, &response)
	return response, err
}

// AddBounce files a bounce report.
// Address identifies the intended recipient of the message that bounced.
// Code corresponds to the numeric response given by the e-mail server which rejected the message.
// Error provides the corresponding human readable reason for the problem.
// For example,
// here's how the these two fields relate.
// Suppose the SMTP server responds with an error, as below.
// Then, . . .
//
//      550  Requested action not taken: mailbox unavailable
//     \___/\_______________________________________________/
//       |                         |
//       `-- Code                  `-- Error
//
// Note that both code and error exist as strings, even though
// code will report as a number.
func (mg *MailgunImpl) AddBounce(ctx context.Context, address, code, error string) error {
	r := newHTTPRequest(generateApiUrl(mg, bouncesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("address", address)
	if code != "" {
		payload.addValue("code", code)
	}
	if error != "" {
		payload.addValue("error", error)
	}
	_, err := makePostRequest(ctx, r, payload)
	return err
}

// DeleteBounce removes all bounces associted with the provided e-mail address.
func (mg *MailgunImpl) DeleteBounce(ctx context.Context, address string) error {
	r := newHTTPRequest(generateApiUrl(mg, bouncesEndpoint) + "/" + address)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// DeleteBounceList removes all bounces in the bounce list
func (mg *MailgunImpl) DeleteBounceList(ctx context.Context) error {
	r := newHTTPRequest(generateApiUrl(mg, bouncesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}
