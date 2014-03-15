package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
	"time"
)

// DefaultLimit and DefaultSkip instruct the SDK to rely on Mailgun's reasonable defaults for pagination settings.
const (
	DefaultLimit = -1
	DefaultSkip  = -1
)

// Disabled, Tag, and Delete indicate spam actions.
// Disabled prevents Mailgun from taking any action on what it perceives to be spam.
// Tag instruments the received message with headers providing a measure of its spamness.
// Delete instructs Mailgun to just block or delete the message all-together.
const (
	Tag      = "tag"
	Disabled = "disabled"
	Delete   = "delete"
)

// Holds information about a domain used when sending mail.
// The SpamAction field must be one of Tag, Disabled, or Delete.
type Domain struct {
	CreatedAt    string `json:"created_at"`
	SMTPLogin    string `json:"smtp_login"`
	Name         string `json:"name"`
	SMTPPassword string `json:"smtp_password"`
	Wildcard     bool   `json:"wildcard"`
	SpamAction   string `json:"spam_action"`
}

type DNSRecord struct {
	Priority   string
	RecordType string `json:"record_type"`
	Valid      string
	Value      string
}

// Used to decode the domains JSON response.
type domainsEnvelope struct {
	TotalCount int      `json:"total_count"`
	Items      []Domain `json:"items"`
}

type singleDomainEnvelope struct {
	Domain              Domain      `json:"domain"`
	ReceivingDNSRecords []DNSRecord `json:"receiving_dns_records"`
	SendingDNSRecords   []DNSRecord `json:"sending_dns_records"`
}

func (d Domain) GetCreatedAt() (t time.Time, err error) {
	t, err = parseMailgunTime(d.CreatedAt)
	return
}

// GetDomains retrieves a set of domains from Mailgun.
// The limit parameter indicates how many domains to constrain the result to.
// If set to DefaultLimit, it will defer to Mailgun's best judgement (currently, 100).
// The skip parameter indicates where to start the retrieval of domains from.
// If set to DefaultSkip, it will start at the beginning of the list (skip of 0).
//
// Assuming no error, both the number of items retrieved and a slice of Domain instances.
// The number of items returned may be less than the specified limit, if it's specified.
// Note that zero items and a zero-length slice do not necessarily imply an error occurred.
// Except for the error itself, all results are undefined in the event of an error.
func (m *mailgunImpl) GetDomains(limit, skip int) (int, []Domain, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(domainsEndpoint))
	if limit != DefaultLimit {
		r.AddParameter("limit", strconv.Itoa(limit))
	}
	if skip != DefaultSkip {
		r.AddParameter("skip", strconv.Itoa(skip))
	}
	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var envelope domainsEnvelope
	err := r.GetResponseFromJSON(&envelope)
	if err != nil {
		return -1, nil, err
	}
	return envelope.TotalCount, envelope.Items, nil
}

// Retrieve detailed information about the named domain.
func (m *mailgunImpl) GetSingleDomain(domain string) (Domain, []DNSRecord, []DNSRecord, error) {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(domainsEndpoint) + "/" + domain)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	var envelope singleDomainEnvelope
	err := r.GetResponseFromJSON(&envelope)
	if err != nil {
		return Domain{}, nil, nil, err
	}
	return envelope.Domain, envelope.ReceivingDNSRecords, envelope.SendingDNSRecords, nil
}

// CreateDomain instructs Mailgun to create a new domain for your account.
// The name parameter identifies the domain.
// The smtpPassword parameter provides an access credential for the domain.
// The spamAction domain must be one of Delete, Tag, or Disabled.
// The wildcard parameter instructs Mailgun to treat all subdomains of this domain uniformly if true,
// and as different domains if false.
func (m *mailgunImpl) CreateDomain(name string, smtpPassword string, spamAction string, wildcard bool) error {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(domainsEndpoint))
	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	payload := simplehttp.NewUrlEncodedPayload()
	payload.AddValue("name", name)
	payload.AddValue("smtp_password", smtpPassword)
	payload.AddValue("spam_action", spamAction)
	payload.AddValue("wildcard", strconv.FormatBool(wildcard))
	_, err := r.MakePostRequest(payload)
	return err
}

// DeleteDomain instructs Mailgun to dispose of the named domain name.
func (m *mailgunImpl) DeleteDomain(name string) error {
	r := simplehttp.NewHTTPRequest(generatePublicApiUrl(domainsEndpoint) + "/" + name)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	_, err := r.MakeDeleteRequest()
	return err
}
