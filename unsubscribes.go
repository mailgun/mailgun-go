package mailgun

import (
	"context"
	"strconv"
)

type Unsubscribe struct {
	CreatedAt RFC2822Time `json:"created_at"`
	Tags      []string    `json:"tags"`
	ID        string      `json:"id"`
	Address   string      `json:"address"`
}

type unsubscribesResponse struct {
	Paging Paging        `json:"paging"`
	Items  []Unsubscribe `json:"items"`
}

// Fetches the list of unsubscribes
func (mg *MailgunImpl) ListUnsubscribes(opts *ListOptions) *UnsubscribesIterator {
	r := newHTTPRequest(generateApiUrl(mg, unsubscribesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	if opts != nil {
		if opts.Limit != 0 {
			r.addParameter("limit", strconv.Itoa(opts.Limit))
		}
	}
	url, err := r.generateUrlWithParameters()
	return &UnsubscribesIterator{
		mg:                   mg,
		unsubscribesResponse: unsubscribesResponse{Paging: Paging{Next: url, First: url}},
		err:                  err,
	}
}

type UnsubscribesIterator struct {
	unsubscribesResponse
	mg  Mailgun
	err error
}

// If an error occurred during iteration `Err()` will return non nil
func (ci *UnsubscribesIterator) Err() error {
	return ci.err
}

// Next retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ci *UnsubscribesIterator) Next(ctx context.Context, items *[]Unsubscribe) bool {
	if ci.err != nil {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.Next)
	if ci.err != nil {
		return false
	}
	cpy := make([]Unsubscribe, len(ci.Items))
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
func (ci *UnsubscribesIterator) First(ctx context.Context, items *[]Unsubscribe) bool {
	if ci.err != nil {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.First)
	if ci.err != nil {
		return false
	}
	cpy := make([]Unsubscribe, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ci *UnsubscribesIterator) Last(ctx context.Context, items *[]Unsubscribe) bool {
	if ci.err != nil {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.Last)
	if ci.err != nil {
		return false
	}
	cpy := make([]Unsubscribe, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ci *UnsubscribesIterator) Previous(ctx context.Context, items *[]Unsubscribe) bool {
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
	cpy := make([]Unsubscribe, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	if len(ci.Items) == 0 {
		return false
	}
	return true
}

func (ci *UnsubscribesIterator) fetch(ctx context.Context, url string) error {
	r := newHTTPRequest(url)
	r.setClient(ci.mg.Client())
	r.setBasicAuth(basicAuthUser, ci.mg.APIKey())

	return getResponseFromJSON(ctx, r, &ci.unsubscribesResponse)
}

// Retreives a single unsubscribe record. Can be used to check if a given address is present in the list of unsubscribed users.
func (mg *MailgunImpl) GetUnsubscribe(ctx context.Context, address string) (Unsubscribe, error) {
	r := newHTTPRequest(generateApiUrlWithTarget(mg, unsubscribesEndpoint, address))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	envelope := Unsubscribe{}
	err := getResponseFromJSON(ctx, r, &envelope)

	return envelope, err
}

// Unsubscribe adds an e-mail address to the domain's unsubscription table.
func (mg *MailgunImpl) CreateUnsubscribe(ctx context.Context, address, tag string) error {
	r := newHTTPRequest(generateApiUrl(mg, unsubscribesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("address", address)
	p.addValue("tag", tag)
	_, err := makePostRequest(ctx, r, p)
	return err
}

// DeleteUnsubscribe removes the e-mail address given from the domain's unsubscription table.
// If passing in an ID (discoverable from, e.g., ListUnsubscribes()), the e-mail address associated
// with the given ID will be removed.
func (mg *MailgunImpl) DeleteUnsubscribe(ctx context.Context, address string) error {
	r := newHTTPRequest(generateApiUrlWithTarget(mg, unsubscribesEndpoint, address))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// DeleteUnsubscribeWithTag removes the e-mail address given from the domain's unsubscription table with a matching tag.
// If passing in an ID (discoverable from, e.g., ListUnsubscribes()), the e-mail address associated
// with the given ID will be removed.
func (mg *MailgunImpl) DeleteUnsubscribeWithTag(ctx context.Context, a, t string) error {
	r := newHTTPRequest(generateApiUrlWithTarget(mg, unsubscribesEndpoint, a))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.addParameter("tag", t)
	_, err := makeDeleteRequest(ctx, r)
	return err
}
