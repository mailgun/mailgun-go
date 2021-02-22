package mailgun

import (
	"context"
	"fmt"
	"strconv"
)

// A Credential structure describes a principle allowed to send or receive mail at the domain.
type Credential struct {
	CreatedAt RFC2822Time `json:"created_at"`
	Login     string      `json:"login"`
	Password  string      `json:"password"`
}

type credentialsListResponse struct {
	// is -1 if Next() or First() have not been called
	TotalCount int          `json:"total_count"`
	Items      []Credential `json:"items"`
}

// Returned when a required parameter is missing.
var ErrEmptyParam = fmt.Errorf("empty or illegal parameter")

// ListCredentials returns the (possibly zero-length) list of credentials associated with your domain.
func (mg *MailgunImpl) ListCredentials(opts *ListOptions) *CredentialsIterator {
	var limit int
	if opts != nil {
		limit = opts.Limit
	}

	if limit == 0 {
		limit = 100
	}
	return &CredentialsIterator{
		mg:                      mg,
		url:                     generateCredentialsUrl(mg, ""),
		credentialsListResponse: credentialsListResponse{TotalCount: -1},
		limit:                   limit,
	}
}

type CredentialsIterator struct {
	credentialsListResponse

	limit  int
	mg     Mailgun
	offset int
	url    string
	err    error
}

// If an error occurred during iteration `Err()` will return non nil
func (ri *CredentialsIterator) Err() error {
	return ri.err
}

// Offset returns the current offset of the iterator
func (ri *CredentialsIterator) Offset() int {
	return ri.offset
}

// Next retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ri *CredentialsIterator) Next(ctx context.Context, items *[]Credential) bool {
	if ri.err != nil {
		return false
	}

	ri.err = ri.fetch(ctx, ri.offset, ri.limit)
	if ri.err != nil {
		return false
	}

	cpy := make([]Credential, len(ri.Items))
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
func (ri *CredentialsIterator) First(ctx context.Context, items *[]Credential) bool {
	if ri.err != nil {
		return false
	}
	ri.err = ri.fetch(ctx, 0, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]Credential, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	ri.offset = len(ri.Items)
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ri *CredentialsIterator) Last(ctx context.Context, items *[]Credential) bool {
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

	ri.err = ri.fetch(ctx, ri.offset, ri.limit)
	if ri.err != nil {
		return false
	}
	cpy := make([]Credential, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ri *CredentialsIterator) Previous(ctx context.Context, items *[]Credential) bool {
	if ri.err != nil {
		return false
	}

	if ri.TotalCount == -1 {
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
	cpy := make([]Credential, len(ri.Items))
	copy(cpy, ri.Items)
	*items = cpy
	if len(ri.Items) == 0 {
		return false
	}
	return true
}

func (ri *CredentialsIterator) fetch(ctx context.Context, skip, limit int) error {
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

	return getResponseFromJSON(ctx, r, &ri.credentialsListResponse)
}

// CreateCredential attempts to create associate a new principle with your domain.
func (mg *MailgunImpl) CreateCredential(ctx context.Context, login, password string) error {
	if (login == "") || (password == "") {
		return ErrEmptyParam
	}
	r := newHTTPRequest(generateCredentialsUrl(mg, ""))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("login", login)
	p.addValue("password", password)
	_, err := makePostRequest(ctx, r, p)
	return err
}

// ChangeCredentialPassword attempts to alter the indicated credential's password.
func (mg *MailgunImpl) ChangeCredentialPassword(ctx context.Context, login, password string) error {
	if (login == "") || (password == "") {
		return ErrEmptyParam
	}
	r := newHTTPRequest(generateCredentialsUrl(mg, login))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("password", password)
	_, err := makePutRequest(ctx, r, p)
	return err
}

// DeleteCredential attempts to remove the indicated principle from the domain.
func (mg *MailgunImpl) DeleteCredential(ctx context.Context, login string) error {
	if login == "" {
		return ErrEmptyParam
	}
	r := newHTTPRequest(generateCredentialsUrl(mg, login))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}
