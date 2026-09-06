package mailgun

import (
	"context"
	"iter"
	"strconv"
	"strings"

	"github.com/mailgun/mailgun-go/v5/mtypes"
)

type ListDomainsOptions struct {
	Limit int

	// Get only domains with a specific state.
	State *mtypes.DomainState

	// If sorting is not specified domains are returned in reverse creation date order.
	Sort *string

	// Get only domains with a specific authority.
	Authority *string

	// Search domains by the given partial or complete name. Does not support wildcards.
	Search *string

	// Search on every domain that belongs to any subaccounts under this account.
	IncludeSubaccounts *bool
}

// ListDomainsIter returns an iterator over pages of domains.
// Iteration stops after the first error or after fetching all items.
// https://documentation.mailgun.com/docs/mailgun/api-reference/send/mailgun/domains/get-v4-domains
func (mg *Client) ListDomainsIter(ctx context.Context, opts *ListDomainsOptions) iter.Seq2[[]mtypes.Domain, error] {
	limit := 100
	if opts != nil && opts.Limit != 0 {
		limit = opts.Limit
	}

	url := generateApiUrl(mg, 4, domainsEndpoint)

	return func(yield func([]mtypes.Domain, error) bool) {
		var skip int
		for {
			resp, err := mg.fetchDomains(ctx, url, opts, skip, limit)
			if err != nil {
				yield(nil, err)
				return
			}

			if len(resp.Items) == 0 {
				return
			}

			if !yield(resp.Items, nil) {
				return
			}

			skip += limit
			if len(resp.Items) < limit || skip >= resp.TotalCount {
				return
			}
		}
	}
}

func (mg *Client) fetchDomains(ctx context.Context, url string, opts *ListDomainsOptions, skip, limit int,
) (mtypes.ListDomainsResponse, error) {
	r := newHTTPRequest(url)
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.setClient(mg.HTTPClient())

	if skip != 0 {
		r.addParameter("skip", strconv.Itoa(skip))
	}

	if limit != 0 {
		r.addParameter("limit", strconv.Itoa(limit))
	}

	if opts != nil {
		if opts.State != nil {
			r.addParameter("state", string(*opts.State))
		}
		if opts.Sort != nil {
			r.addParameter("sort", *opts.Sort)
		}
		if opts.Authority != nil {
			r.addParameter("authority", *opts.Authority)
		}
		if opts.Search != nil {
			r.addParameter("search", *opts.Search)
		}
		if opts.IncludeSubaccounts != nil {
			r.addParameter("include_subaccounts", strconv.FormatBool(*opts.IncludeSubaccounts))
		}
	}

	var resp mtypes.ListDomainsResponse
	err := getResponseFromJSON(ctx, r, &resp)

	return resp, err
}

type GetDomainOptions struct {
	// If set to true, domain payload will include dkim_host, mailfrom_host and pod
	Extended *bool

	// Domain payload will include sending and receiving dns records payload
	WithDNS *bool
}

// GetDomain retrieves detailed information about the named domain.
// https://documentation.mailgun.com/docs/mailgun/api-reference/send/mailgun/domains/get-v4-domains--name-
func (mg *Client) GetDomain(ctx context.Context, domain string, opts *GetDomainOptions) (mtypes.GetDomainResponse, error) {
	r := newHTTPRequest(generateApiUrl(mg, 4, domainsEndpoint) + "/" + domain)
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	if opts != nil {
		if opts.Extended != nil {
			r.addParameter("h:extended", strconv.FormatBool(*opts.Extended))
		}

		if opts.WithDNS != nil {
			r.addParameter("h:with_dns", strconv.FormatBool(*opts.WithDNS))
		}
	}

	var resp mtypes.GetDomainResponse
	err := getResponseFromJSON(ctx, r, &resp)
	return resp, err
}

// VerifyDomain verifies the domains DNS records (includes A, CNAME, SPF,
// DKIM and MX records) to ensure the domain is ready and able to send.
func (mg *Client) VerifyDomain(ctx context.Context, domain string) (mtypes.GetDomainResponse, error) {
	r := newHTTPRequest(generateApiUrl(mg, 4, domainsEndpoint) + "/" + domain + "/verify")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	// TODO(vtopc): why newUrlEncodedPayload()?
	payload := newUrlEncodedPayload()
	var resp mtypes.GetDomainResponse
	err := putResponseFromJSON(ctx, r, payload, &resp)
	return resp, err
}

// VerifyAndReturnDomain verifies the domains DNS records (includes A, CNAME, SPF,
// DKIM and MX records) to ensure the domain is ready and able to send.
//
// Deprecated: use VerifyDomain instead.
//
// TODO(v6): remove this method
func (mg *Client) VerifyAndReturnDomain(ctx context.Context, domain string) (mtypes.GetDomainResponse, error) {
	return mg.VerifyDomain(ctx, domain)
}

// CreateDomainOptions - optional parameters when creating a domain
// https://documentation.mailgun.com/docs/mailgun/api-reference/openapi-final/tag/Domains/#tag/Domains/operation/POST-v4-domains
type CreateDomainOptions struct {
	Password   string
	SpamAction mtypes.SpamAction
	// Wildcard parameter instructs Mailgun to treat all subdomains of this domain uniformly if true,
	// and as different domains if false.
	Wildcard                   bool
	ForceDKIMAuthority         bool
	DKIMKeySize                int
	IPs                        []string
	WebScheme                  string
	UseAutomaticSenderSecurity bool
	ArchiveTo                  string
	DKIMHostName               string
	DKIMSelector               string
	ForceRootDKIMHost          bool
	EncryptIncomingMessage     bool
	PoolID                     string
	RequireTLS                 bool
	SkipVerification           bool
	WebPrefix                  string
	MessageTTL                 int
}

// CreateDomain instructs Mailgun to create a new domain for your account.
// The domain parameter identifies the domain.
// https://documentation.mailgun.com/docs/mailgun/api-reference/send/mailgun/domains/post-v4-domains
func (mg *Client) CreateDomain(ctx context.Context, domain string, opts *CreateDomainOptions) (mtypes.GetDomainResponse, error) {
	r := newHTTPRequest(generateApiUrl(mg, 4, domainsEndpoint))
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := NewFormDataPayload()
	payload.addValue("name", domain)

	if opts != nil {
		if opts.SpamAction != "" {
			payload.addValue("spam_action", string(opts.SpamAction))
		}
		if opts.Wildcard {
			payload.addValue("wildcard", boolToString(opts.Wildcard))
		}
		if opts.ForceDKIMAuthority {
			payload.addValue("force_dkim_authority", boolToString(opts.ForceDKIMAuthority))
		}
		if opts.DKIMKeySize != 0 {
			payload.addValue("dkim_key_size", strconv.Itoa(opts.DKIMKeySize))
		}
		if len(opts.IPs) != 0 {
			payload.addValue("ips", strings.Join(opts.IPs, ","))
		}
		if opts.Password != "" {
			payload.addValue("smtp_password", opts.Password)
		}
		if opts.WebScheme != "" {
			payload.addValue("web_scheme", opts.WebScheme)
		}
		if opts.UseAutomaticSenderSecurity {
			payload.addValue("use_automatic_sender_security", boolToString(opts.UseAutomaticSenderSecurity))
		}
		if opts.ArchiveTo != "" {
			payload.addValue("archive_to", opts.ArchiveTo)
		}
		if opts.DKIMHostName != "" {
			payload.addValue("dkim_host_name", opts.DKIMHostName)
		}
		if opts.DKIMSelector != "" {
			payload.addValue("dkim_selector", opts.DKIMSelector)
		}
		if opts.ForceRootDKIMHost {
			payload.addValue("force_root_dkim_host", boolToString(opts.ForceRootDKIMHost))
		}
		if opts.EncryptIncomingMessage {
			payload.addValue("encrypt_incoming_message", boolToString(opts.EncryptIncomingMessage))
		}
		if opts.PoolID != "" {
			payload.addValue("pool_id", opts.PoolID)
		}
		if opts.RequireTLS {
			payload.addValue("require_tls", boolToString(opts.RequireTLS))
		}
		if opts.SkipVerification {
			payload.addValue("skip_verification", boolToString(opts.SkipVerification))
		}
		if opts.WebPrefix != "" {
			payload.addValue("web_prefix", opts.WebPrefix)
		}
		if opts.MessageTTL != 0 {
			payload.addValue("message_ttl", strconv.Itoa(opts.MessageTTL))
		}
	}
	var resp mtypes.GetDomainResponse
	err := postResponseFromJSON(ctx, r, payload, &resp)
	return resp, err
}

// DeleteDomain instructs Mailgun to dispose of the named domain name
func (mg *Client) DeleteDomain(ctx context.Context, domain string) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain)
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// UpdateDomainOptions options for updating a domain
type UpdateDomainOptions struct {
	Password                   string
	SpamAction                 mtypes.SpamAction
	Wildcard                   *bool
	WebScheme                  string
	WebPrefix                  string
	RequireTLS                 *bool
	SkipVerification           *bool
	UseAutomaticSenderSecurity *bool
	ArchiveTo                  string
	MailFromHost               string
	MessageTTL                 *int
}

// UpdateDomain updates a domain's attributes.
func (mg *Client) UpdateDomain(ctx context.Context, domain string, opts *UpdateDomainOptions) error {
	r := newHTTPRequest(generateApiUrl(mg, 4, domainsEndpoint) + "/" + domain)
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := NewFormDataPayload()

	if opts != nil {
		if opts.Password != "" {
			payload.addValue("smtp_password", opts.Password)
		}
		if opts.SpamAction != "" {
			payload.addValue("spam_action", string(opts.SpamAction))
		}
		if opts.Wildcard != nil {
			payload.addValue("wildcard", boolToString(*opts.Wildcard))
		}
		if opts.WebScheme != "" {
			payload.addValue("web_scheme", opts.WebScheme)
		}
		if opts.WebPrefix != "" {
			payload.addValue("web_prefix", opts.WebPrefix)
		}
		if opts.RequireTLS != nil {
			payload.addValue("require_tls", boolToString(*opts.RequireTLS))
		}
		if opts.SkipVerification != nil {
			payload.addValue("skip_verification", boolToString(*opts.SkipVerification))
		}
		if opts.UseAutomaticSenderSecurity != nil {
			payload.addValue("use_automatic_sender_security", boolToString(*opts.UseAutomaticSenderSecurity))
		}
		if opts.ArchiveTo != "" {
			payload.addValue("archive_to", opts.ArchiveTo)
		}
		if opts.MailFromHost != "" {
			payload.addValue("mailfrom_host", opts.MailFromHost)
		}
		if opts.MessageTTL != nil {
			payload.addValue("message_ttl", strconv.Itoa(*opts.MessageTTL))
		}
	}

	_, err := makePutRequest(ctx, r, payload)

	return err
}
