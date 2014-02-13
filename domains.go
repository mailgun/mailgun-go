package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
	"time"
)

// Holds information about a domain used when sending mail.
type Domain struct {
	CreatedAt    string `json:"created_at"`
	SMTPLogin    string `json:"smtp_login"`
	Name         string `json:"name"`
	SMTPPassword string `json:"smtp_password"`
	Wildcard     bool   `json:"wildcard"`
	SpamAction   bool   `json:"spam_action"`
}

type DNSRecord struct {
	Priority   string
	RecordType string
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

func (m *mailgunImpl) GetDomains(limit, skip int) (int, []Domain, error) {
	r := simplehttp.NewGetRequest(generatePublicApiUrl(domainsEndpoint))
	if limit != -1 {
		r.AddParameter("limit", strconv.Itoa(limit))
	}
	if skip != -1 {
		r.AddParameter("skip", strconv.Itoa(skip))
	}
	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var envelope domainsEnvelope
	err := r.MakeJSONRequest(&envelope)
	if err != nil {
		return -1, nil, err
	}
	return envelope.TotalCount, envelope.Items, nil
}

func (m *mailgunImpl) GetSingleDomain(domain string) (Domain, []DNSRecord, []DNSRecord, error) {
	r := simplehttp.NewGetRequest(generatePublicApiUrl(domainsEndpoint) + "/" + domain)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	var envelope singleDomainEnvelope
	err := r.MakeJSONRequest(&envelope)
	if err != nil {
		return Domain{}, nil, nil, err
	}
	return envelope.Domain, envelope.ReceivingDNSRecords, envelope.SendingDNSRecords, nil
}

func (m *mailgunImpl) CreateDomain(name string, smtpPassword string, spamAction bool, wildcard bool) error {
	r := simplehttp.NewPostRequest(generatePublicApiUrl(domainsEndpoint))
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	r.AddFormValue("name", name)
	r.AddFormValue("smtp_password", smtpPassword)
	if spamAction {
		r.AddFormValue("spam_action", "tag")
	} else {
		r.AddFormValue("spam_action", "disabled")
	}
	r.AddFormValue("wildcard", strconv.FormatBool(wildcard))
	_, err := r.MakeRequest()
	return err
}

func (m *mailgunImpl) DeleteDomain(name string) error {
	r := simplehttp.NewGetRequest(generatePublicApiUrl(domainsEndpoint) + "/" + name)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	_, err := r.MakeRequest()
	return err
}
