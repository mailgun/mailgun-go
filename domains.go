package mailgun

import (
	"strconv"
	"time"
)

// DefaultLimit and DefaultSkip instruct the SDK to rely on Mailgun's reasonable defaults for Paging settings.
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

// A Domain structure holds information about a domain used when sending mail.
// The SpamAction field must be one of Tag, Disabled, or Delete.
type Domain struct {
	CreatedAt    string `json:"created_at"`
	SMTPLogin    string `json:"smtp_login"`
	Name         string `json:"name"`
	SMTPPassword string `json:"smtp_password"`
	Wildcard     bool   `json:"wildcard"`
	SpamAction   string `json:"spam_action"`
	State        string `json:"state"`
}

// DNSRecord structures describe intended records to properly configure your domain for use with Mailgun.
// Note that Mailgun does not host DNS records.
type DNSRecord struct {
	Priority   string
	RecordType string `json:"record_type"`
	Valid      string
	Name       string
	Value      string
}

type domainResponse struct {
	Domain              Domain            `json:"domain"`
	ReceivingDNSRecords []DNSRecord       `json:"receiving_dns_records"`
	SendingDNSRecords   []DNSRecord       `json:"sending_dns_records"`
	Connection          *DomainConnection `json:"connection,omitempty"`
	Tracking            *DomainTracking   `json:"tracking,omitempty"`
	TagLimits           *TagLimits        `json:"limits,omitempty"`
}

type domainConnectionResponse struct {
	Connection DomainConnection `json:"connection"`
}

type DomainList struct {
	TotalCount int      `json:"total_count"`
	Items      []Domain `json:"items"`
}

type DomainConnection struct {
	RequireTLS       bool `json:"require_tls"`
	SkipVerification bool `json:"skip_verification"`
}

type DomainTracking struct {
	Click       TrackingStatus `json:"click"`
	Open        TrackingStatus `json:"open"`
	Unsubscribe TrackingStatus `json:"unsubscribe"`
}

type TrackingStatus struct {
	Active     bool   `json:"active"`
	HTMLFooter string `json:"html_footer"`
	TextFooter string `json:"text_footer"`
}

type domainTrackingResponse struct {
	Tracking DomainTracking `json:"tracking"`
}

// GetCreatedAt returns the time the domain was created as a normal Go time.Time type.
func (d Domain) GetCreatedAt() (t time.Time, err error) {
	t, err = parseMailgunTime(d.CreatedAt)
	return
}

// GetDomains retrieves a set of domains from Mailgun.
//
// Assuming no error, both the number of items retrieved and a slice of Domain instances.
// The number of items returned may be less than the specified limit, if it's specified.
// Note that zero items and a zero-length slice do not necessarily imply an error occurred.
// Except for the error itself, all results are undefined in the event of an error.
func (mg *MailgunImpl) GetDomains(limit, skip int) (int, []Domain, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint))
	r.setClient(mg.Client())
	if limit != DefaultLimit {
		r.addParameter("limit", strconv.Itoa(limit))
	}
	if skip != DefaultSkip {
		r.addParameter("skip", strconv.Itoa(skip))
	}
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var list DomainList
	err := getResponseFromJSON(r, &list)
	if err != nil {
		return -1, nil, err
	}
	return list.TotalCount, list.Items, nil
}

// Retrieve detailed information about the named domain.
func (mg *MailgunImpl) GetSingleDomain(domain string) (Domain, []DNSRecord, []DNSRecord, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + domain)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var envelope domainResponse
	err := getResponseFromJSON(r, &envelope)
	return envelope.Domain, envelope.ReceivingDNSRecords, envelope.SendingDNSRecords, err
}

// CreateDomain instructs Mailgun to create a new domain for your account.
// The name parameter identifies the domain.
// The smtpPassword parameter provides an access credential for the domain.
// The spamAction domain must be one of Delete, Tag, or Disabled.
// The wildcard parameter instructs Mailgun to treat all subdomains of this domain uniformly if true,
// and as different domains if false.
func (mg *MailgunImpl) CreateDomain(name string, smtpPassword string, spamAction string, wildcard bool) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("name", name)
	payload.addValue("smtp_password", smtpPassword)
	payload.addValue("spam_action", spamAction)
	payload.addValue("wildcard", strconv.FormatBool(wildcard))
	_, err := makePostRequest(r, payload)
	return err
}

// Returns delivery connection settings for the defined domain
func (mg *MailgunImpl) GetDomainConnection(domain string) (DomainConnection, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + domain + "/connection")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp domainConnectionResponse
	err := getResponseFromJSON(r, &resp)
	return resp.Connection, err
}

// Updates the specified delivery connection settings for the defined domain
func (mg *MailgunImpl) UpdateDomainConnection(domain string, settings DomainConnection) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + domain + "/connection")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("require_tls", boolToString(settings.RequireTLS))
	payload.addValue("skip_verification", boolToString(settings.SkipVerification))
	_, err := makePutRequest(r, payload)
	return err
}

// DeleteDomain instructs Mailgun to dispose of the named domain name
func (mg *MailgunImpl) DeleteDomain(name string) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + name)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(r)
	return err
}

// Returns tracking settings for a domain
func (mg *MailgunImpl) GetDomainTracking(domain string) (DomainTracking, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + domain + "/tracking")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp domainTrackingResponse
	err := getResponseFromJSON(r, &resp)
	return resp.Tracking, err
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
