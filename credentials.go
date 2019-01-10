package mailgun

import (
	"context"
	"fmt"
	"strconv"
)

// A Credential structure describes a principle allowed to send or receive mail at the domain.
type Credential struct {
	CreatedAt string `json:"created_at"`
	Login     string `json:"login"`
	Password  string `json:"password"`
}

// ErrEmptyParam results occur when a required parameter is missing.
var ErrEmptyParam = fmt.Errorf("empty or illegal parameter")

// ListCredentials returns the (possibly zero-length) list of credentials associated with your domain.
func (mg *MailgunImpl) ListCredentials(ctx context.Context, opts *ListOptions) ([]Credential, error) {
	r := newHTTPRequest(generateCredentialsUrl(mg, ""))
	r.setClient(mg.Client())
	if opts != nil && opts.Limit != 0 {
		r.addParameter("limit", strconv.Itoa(opts.Limit))
	}

	if opts != nil && opts.Skip != 0 {
		r.addParameter("skip", strconv.Itoa(opts.Skip))
	}

	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var envelope struct {
		TotalCount int          `json:"total_count"`
		Items      []Credential `json:"items"`
	}
	err := getResponseFromJSON(ctx, r, &envelope)
	if err != nil {
		return nil, err
	}
	return envelope.Items, nil
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
func (mg *MailgunImpl) ChangeCredentialPassword(ctx context.Context, id, password string) error {
	if (id == "") || (password == "") {
		return ErrEmptyParam
	}
	r := newHTTPRequest(generateCredentialsUrl(mg, id))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("password", password)
	_, err := makePutRequest(ctx, r, p)
	return err
}

// DeleteCredential attempts to remove the indicated principle from the domain.
func (mg *MailgunImpl) DeleteCredential(ctx context.Context, id string) error {
	if id == "" {
		return ErrEmptyParam
	}
	r := newHTTPRequest(generateCredentialsUrl(mg, id))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}
