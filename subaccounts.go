package mailgun

import (
	"context"
	"fmt"
	"strconv"
)

type ListSubaccountsOptions struct {
	Limit     int
	Skip      int
	SortArray string
	Enabled   bool
}

type SubaccountsIterator struct {
	subaccountsListResponse

	mg        Mailgun
	limit     int
	offset    int
	skip      int
	sortArray string
	enabled   bool
	url       string
	err       error
}

// A Subaccount structure holds information about a subaccount.
type Subaccount struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type SubaccountResponse struct {
	Item Subaccount `json:"subaccount"`
}

type subaccountsListResponse struct {
	Items []Subaccount `json:"subaccounts"`
	Total int          `json:"total"`
}

// ListSubaccounts retrieves a set of subaccount linked to the primary Mailgun account.
func (mg *MailgunImpl) ListSubaccounts(opts *ListSubaccountsOptions) *SubaccountsIterator {
	r := newHTTPRequest(generateSubaccountsApiUrl(mg))
	r.setClient(mg.client)
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var limit, skip int
	var sortArray string
	var enabled bool
	if opts != nil {
		limit = opts.Limit
		skip = opts.Skip
		sortArray = opts.SortArray
		enabled = opts.Enabled
	}
	if limit == 0 {
		limit = 10
	}

	return &SubaccountsIterator{
		mg:                      mg,
		url:                     generateSubaccountsApiUrl(mg),
		subaccountsListResponse: subaccountsListResponse{Total: -1},
		limit:                   limit,
		skip:                    skip,
		sortArray:               sortArray,
		enabled:                 enabled,
	}
}

// If an error occurred during iteration `Err()` will return non nil
func (ri *SubaccountsIterator) Err() error {
	return ri.err
}

// Offset returns the current offset of the iterator
func (ri *SubaccountsIterator) Offset() int {
	return ri.offset
}

// Next retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ri *SubaccountsIterator) Next(ctx context.Context, items *[]Subaccount) bool {
	if ri.err != nil {
		return false
	}

	ri.err = ri.fetch(ctx, ri.offset, ri.limit)
	if ri.err != nil {
		return false
	}

	cpy := make([]Subaccount, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	if len(ri.Items) == 0 {
		return false
	}
	ri.offset = ri.offset + len(ri.Items)
	return true
}

// First retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (ri *SubaccountsIterator) First(ctx context.Context, items *[]Subaccount) bool {
	if ri.err != nil {
		return false
	}
	ri.err = ri.fetch(ctx, 0, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]Subaccount, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	ri.offset = len(ri.Items)
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ri *SubaccountsIterator) Last(ctx context.Context, items *[]Subaccount) bool {
	if ri.err != nil {
		return false
	}

	if ri.Total == -1 {
		return false
	}

	ri.offset = ri.Total - ri.limit
	if ri.offset < 0 {
		ri.offset = 0
	}

	ri.err = ri.fetch(ctx, ri.offset, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]Subaccount, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ri *SubaccountsIterator) Previous(ctx context.Context, items *[]Subaccount) bool {
	if ri.err != nil {
		return false
	}

	if ri.Total == -1 {
		return false
	}

	ri.offset = ri.offset - (ri.limit * 2)
	if ri.offset < 0 {
		ri.offset = 0
	}

	ri.err = ri.fetch(ctx, ri.offset, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]Subaccount, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	if len(ri.Items) == 0 {
		return false
	}
	return true
}

func (ri *SubaccountsIterator) fetch(ctx context.Context, skip, limit int) error {
	ri.Items = nil
	r := newHTTPRequest(ri.url)
	r.setBasicAuth(basicAuthUser, ri.mg.APIKey())
	r.setClient(ri.mg.Client())

	if skip != 0 {
		r.addParameter("skip", strconv.Itoa(skip))
	}
	if limit != 0 {
		r.addParameter("limit", strconv.Itoa(limit))
	}

	return getResponseFromJSON(ctx, r, &ri.subaccountsListResponse)
}

// CreateSubaccount instructs Mailgun to create a new account (Subaccount) that is linked to the primary account.
// Subaccounts are child accounts that share the same plan and usage allocations as the primary, but have their own
// assets (sending domains, unique users, API key, SMTP credentials, settings, statistics and site login).
// All you need is the name of the subaccount.
func (mg *MailgunImpl) CreateSubaccount(ctx context.Context, subaccountName string) (SubaccountResponse, error) {
	r := newHTTPRequest(generateSubaccountsApiUrl(mg))
	r.setClient(mg.client)
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("name", subaccountName)
	resp := SubaccountResponse{}
	err := postResponseFromJSON(ctx, r, payload, &resp)
	return resp, err
}

// SubaccountDetails retrieves detailed information about subaccount using subaccountId.
func (mg *MailgunImpl) SubaccountDetails(ctx context.Context, subaccountId string) (SubaccountResponse, error) {
	r := newHTTPRequest(generateSubaccountsApiUrl(mg) + "/" + subaccountId)
	r.setClient(mg.client)
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp SubaccountResponse
	err := getResponseFromJSON(ctx, r, &resp)
	return resp, err
}

// EnableSubaccount instructs Mailgun to enable subaccount.
func (mg *MailgunImpl) EnableSubaccount(ctx context.Context, subaccountId string) (SubaccountResponse, error) {
	r := newHTTPRequest(generateSubaccountsApiUrl(mg) + "/" + subaccountId + "/" + "enable")
	r.setClient(mg.client)
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	resp := SubaccountResponse{}
	err := postResponseFromJSON(ctx, r, nil, &resp)
	return resp, err
}

// DisableSubaccount instructs Mailgun to disable subaccount.
func (mg *MailgunImpl) DisableSubaccount(ctx context.Context, subaccountId string) (SubaccountResponse, error) {
	r := newHTTPRequest(generateSubaccountsApiUrl(mg) + "/" + subaccountId + "/" + "disable")
	r.setClient(mg.client)
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	resp := SubaccountResponse{}
	err := postResponseFromJSON(ctx, r, nil, &resp)
	return resp, err
}

func generateSubaccountsApiUrl(m Mailgun) string {
	return fmt.Sprintf("%s/%s/%s", m.APIBase(), accountsEndpoint, subaccountsEndpoint)
}
