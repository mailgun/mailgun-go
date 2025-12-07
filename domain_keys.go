package mailgun

import (
	"context"
	"strconv"

	"github.com/mailgun/mailgun-go/v5/mtypes"
)

type ListDomainKeysOptions struct {
	Limit int
}

// ListDomainKeys retrieves a set of domain keys from Mailgun.
func (mg *Client) ListDomainKeys(opts *ListDomainKeysOptions) *DomainKeysIterator {
	var limit int
	if opts != nil {
		limit = opts.Limit
	}

	if limit == 0 {
		limit = 100
	}
	return &DomainKeysIterator{
		mg:                     mg,
		url:                    generateApiUrl(mg, 1, dkimEndpoint+"/keys"),
		ListDomainKeysResponse: mtypes.ListDomainKeysResponse{TotalCount: -1},
		limit:                  limit,
	}
}

type DomainKeysIterator struct {
	mtypes.ListDomainKeysResponse

	limit           int
	mg              Mailgun
	offset          int
	firstPageUrl    string
	previousPageUrl string
	nextPageUrl     string
	lastPageUrl     string
	url             string
	err             error
}

// Err if an error occurred during iteration `Err()` will return non nil
func (ri *DomainKeysIterator) Err() error {
	return ri.err
}

// Offset returns the current offset of the iterator
func (ri *DomainKeysIterator) Offset() int {
	return ri.offset
}

// Next retrieves the next page of items from the api. Returns false when there
// are no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ri *DomainKeysIterator) Next(ctx context.Context, items *[]mtypes.DomainKey) bool {
	if ri.err != nil {
		return false
	}

	ri.err = ri.fetch(ctx, ri.nextPageUrl, ri.limit)
	if ri.err != nil {
		return false
	}

	cpy := make([]mtypes.DomainKey, len(ri.Items))
	*items = cpy
	if len(ri.Items) == 0 {
		return false
	}
	ri.offset += len(ri.Items)
	return true
}

// First retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (ri *DomainKeysIterator) First(ctx context.Context, items *[]mtypes.DomainKey) bool {
	if ri.err != nil {
		return false
	}
	ri.err = ri.fetch(ctx, ri.firstPageUrl, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]mtypes.DomainKey, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	ri.offset = len(ri.Items)
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ri *DomainKeysIterator) Last(ctx context.Context, items *[]mtypes.DomainKey) bool {
	if ri.err != nil {
		return false
	}

	if ri.TotalCount == -1 {
		return false
	}

	ri.offset = ri.TotalCount - ri.limit
	if ri.offset < 0 {
		ri.offset = 0
	}

	ri.err = ri.fetch(ctx, ri.lastPageUrl, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]mtypes.DomainKey, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// are no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ri *DomainKeysIterator) Previous(ctx context.Context, items *[]mtypes.DomainKey) bool {
	if ri.err != nil {
		return false
	}

	if ri.TotalCount == -1 {
		return false
	}

	ri.offset -= ri.limit * 2
	if ri.offset < 0 {
		ri.offset = 0
	}

	ri.err = ri.fetch(ctx, ri.previousPageUrl, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]mtypes.DomainKey, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy

	return len(ri.Items) != 0
}

func (ri *DomainKeysIterator) fetch(ctx context.Context, pageUrl string, limit int) error {
	ri.Items = nil
	url := ri.url
	if pageUrl != "" {
		url = pageUrl
	}
	r := newHTTPRequest(url)

	r.setBasicAuth(basicAuthUser, ri.mg.APIKey())
	r.setClient(ri.mg.HTTPClient())

	if limit != 0 {
		r.addParameter("limit", strconv.Itoa(limit))
	}

	return getResponseFromJSON(ctx, r, &ri.ListDomainKeysResponse)
}

// UpdateDomainDkimSelector updates the DKIM selector for a domain
func (mg *Client) UpdateDomainDkimSelector(ctx context.Context, domain, dkimSelector string) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/dkim_selector")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("dkim_selector", dkimSelector)
	_, err := makePutRequest(ctx, r, payload)
	return err
}
