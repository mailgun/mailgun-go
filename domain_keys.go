package mailgun

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/mailgun/mailgun-go/v5/mtypes"
)

type ListDomainKeysOptions struct {
	Limit int
}

type CreateDomainKeyOptions struct {
	Bits int
	PEM  string
}

type AllDomainsKeysIterator struct {
	mtypes.ListAllDomainsKeysResponse

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

type DomainKeysIterator struct {
	mtypes.ListDomainKeysResponse

	mg              Mailgun
	firstPageUrl    string
	previousPageUrl string
	nextPageUrl     string
	lastPageUrl     string
	url             string
	err             error
}

// ListAllDomainsKeys retrieves a set of domain keys from Mailgun.
func (mg *Client) ListAllDomainsKeys(opts *ListDomainKeysOptions) *AllDomainsKeysIterator {
	var limit int
	if opts != nil {
		limit = opts.Limit
	}

	if limit == 0 {
		limit = 100
	}
	return &AllDomainsKeysIterator{
		mg:                         mg,
		url:                        generateApiUrl(mg, 1, dkimEndpoint+"/keys"),
		ListAllDomainsKeysResponse: mtypes.ListAllDomainsKeysResponse{TotalCount: -1},
		limit:                      limit,
	}
}

// Err if an error occurred during iteration `Err()` will return non nil
func (ri *AllDomainsKeysIterator) Err() error {
	return ri.err
}

// Offset returns the current offset of the iterator
func (ri *AllDomainsKeysIterator) Offset() int {
	return ri.offset
}

// Next retrieves the next page of items from the api. Returns false when there
// are no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ri *AllDomainsKeysIterator) Next(ctx context.Context, items *[]mtypes.DomainKey) bool {
	if ri.err != nil {
		return false
	}

	ri.err = ri.fetch(ctx, ri.nextPageUrl, ri.limit)
	if ri.err != nil {
		return false
	}

	cpy := make([]mtypes.DomainKey, len(ri.Items))
	copy(cpy, ri.Items)
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
func (ri *AllDomainsKeysIterator) First(ctx context.Context, items *[]mtypes.DomainKey) bool {
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
func (ri *AllDomainsKeysIterator) Last(ctx context.Context, items *[]mtypes.DomainKey) bool {
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
func (ri *AllDomainsKeysIterator) Previous(ctx context.Context, items *[]mtypes.DomainKey) bool {
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

func (ri *AllDomainsKeysIterator) fetch(ctx context.Context, pageUrl string, limit int) error {
	ri.Items = nil
	uri := ri.url
	if pageUrl != "" {
		uri = pageUrl
	}
	r := newHTTPRequest(uri)

	r.setBasicAuth(basicAuthUser, ri.mg.APIKey())
	r.setClient(ri.mg.HTTPClient())

	if limit != 0 {
		r.addParameter("limit", strconv.Itoa(limit))
	}

	return getResponseFromJSON(ctx, r, &ri.ListAllDomainsKeysResponse)
}

// CreateDomainKey creates a domain key for the given domain
func (mg *Client) CreateDomainKey(ctx context.Context, domain, dkimSelector string, opts *CreateDomainKeyOptions) (mtypes.DomainKey, error) {
	r := newHTTPRequest(generateApiUrl(mg, 1, dkimEndpoint+"/keys"))
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("signing_domain", domain)
	payload.addValue("selector", dkimSelector)

	if opts.Bits != 0 {
		payload.addValue("bits", strconv.Itoa(opts.Bits))
	}

	if opts.PEM != "" {
		payload.addValue("pem", opts.PEM)
	}

	var resp mtypes.DomainKey
	err := postResponseFromJSON(ctx, r, payload, &resp)
	return resp, err
}

// DeleteDomainKey deletes a domain key from the given domain
func (mg *Client) DeleteDomainKey(ctx context.Context, domain, dkimSelector string) error {
	uri := generateDeleteDomainKeyApiUrl(dkimEndpoint, domain, dkimSelector)

	r := newHTTPRequest(generateApiUrl(mg, 1, uri))
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	_, err := makeDeleteRequest(ctx, r)
	return err
}

// ActivateDomainKey deactivates a domain key for the given domain
func (mg *Client) ActivateDomainKey(ctx context.Context, domain, dkimSelector string) error {
	uri := generateActivateDomainKeyApiUrl(domainsEndpoint, domain, dkimSelector)

	r := newHTTPRequest(generateApiUrl(mg, 4, uri))
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	_, err := makePutRequest(ctx, r, newUrlEncodedPayload())
	return err
}

// ListDomainKeys retrieves a set of domain keys from Mailgun.
func (mg *Client) ListDomainKeys(domain string) *DomainKeysIterator {
	uri := generateListDomainKeysApiUrl(domainsEndpoint, domain)

	return &DomainKeysIterator{
		mg:                     mg,
		url:                    generateApiUrl(mg, 4, uri),
		ListDomainKeysResponse: mtypes.ListDomainKeysResponse{},
	}
}

// Err if an error occurred during iteration `Err()` will return non nil
func (ri *DomainKeysIterator) Err() error {
	return ri.err
}

// Offset returns the current offset of the iterator
func (ri *DomainKeysIterator) Offset() int {
	return len(ri.Items)
}

// Next retrieves the next page of items from the api. Returns false when there
// are no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ri *DomainKeysIterator) Next(ctx context.Context, items *[]mtypes.DomainKey) bool {
	if ri.err != nil {
		return false
	}

	ri.err = ri.fetch(ctx, ri.nextPageUrl)
	if ri.err != nil {
		return false
	}

	cpy := make([]mtypes.DomainKey, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	return len(ri.Items) != 0
}

// First retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (ri *DomainKeysIterator) First(ctx context.Context, items *[]mtypes.DomainKey) bool {
	if ri.err != nil {
		return false
	}
	ri.err = ri.fetch(ctx, ri.firstPageUrl)
	if ri.err != nil {
		return false
	}
	cpy := make([]mtypes.DomainKey, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
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

	ri.err = ri.fetch(ctx, ri.lastPageUrl)
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

	ri.err = ri.fetch(ctx, ri.previousPageUrl)
	if ri.err != nil {
		return false
	}
	cpy := make([]mtypes.DomainKey, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy

	return len(ri.Items) != 0
}

func (ri *DomainKeysIterator) fetch(ctx context.Context, pageUrl string) error {
	ri.Items = nil
	uri := ri.url
	if pageUrl != "" {
		uri = pageUrl
	}
	r := newHTTPRequest(uri)

	r.setBasicAuth(basicAuthUser, ri.mg.APIKey())
	r.setClient(ri.mg.HTTPClient())

	return getResponseFromJSON(ctx, r, &ri.ListDomainKeysResponse)
}

// DeactivateDomainKey deactivates a domain key for the given domain
func (mg *Client) DeactivateDomainKey(ctx context.Context, domain, dkimSelector string) error {
	uri := generateDeactivateDomainKeyApiUrl(domainsEndpoint, domain, dkimSelector)

	r := newHTTPRequest(generateApiUrl(mg, 4, uri))
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	_, err := makePutRequest(ctx, r, newUrlEncodedPayload())
	return err
}

func (mg *Client) UpdateDomainDkimAuthority(ctx context.Context, domain string, self bool) (mtypes.UpdateDomainDkimAuthorityResponse, error) {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/dkim_authority")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("self", boolToString(self))

	var resp mtypes.UpdateDomainDkimAuthorityResponse

	err := putResponseFromJSON(ctx, r, payload, &resp)

	return resp, err
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

// generateDeleteDomainKeyApiUrl renders a URL fragment relevant for deleting a domain key.
func generateDeleteDomainKeyApiUrl(endpoint, domain, dkimSelector string) string {
	params := url.Values{
		"signing_domain": []string{domain},
		"selector":       []string{dkimSelector},
	}
	return fmt.Sprintf("%s/keys?%s", endpoint, params.Encode())
}

// generateActivateDomainKeyApiUrl renders a URL fragment relevant for deactivating a domain key
func generateActivateDomainKeyApiUrl(endpoint, domain, dkimSelector string) string {
	return fmt.Sprintf("%s/%s/keys/%s/activate", endpoint, domain, dkimSelector)
}

// generateListDomainKeysApiUrl renders a URL fragment relevant for listing a domain's keys
func generateListDomainKeysApiUrl(endpoint, domain string) string {
	return fmt.Sprintf("%s/%s/keys", endpoint, domain)
}

// generateDeactivateDomainKeyApiUrl renders a URL fragment relevant for deactivating a domain key
func generateDeactivateDomainKeyApiUrl(endpoint, domain, dkimSelector string) string {
	return fmt.Sprintf("%s/%s/keys/%s/deactivate", endpoint, domain, dkimSelector)
}
