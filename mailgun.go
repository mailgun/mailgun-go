// Package mailgun provides methods for interacting with the Mailgun API.
// For further information please see the Mailgun documentation at
// http://documentation.mailgun.com/
//
// Author: Michael Banzon
package mailgun

import (
	"fmt"
	"time"
	"github.com/mbanzon/simplehttp"
)

const (
	apiBase                 = "https://api.mailgun.net/v2"
	messagesEndpoint        = "messages"
	mimeMessagesEndpoint    = "messages.mime"
	addressValidateEndpoint = "address/validate"
	addressParseEndpoint    = "address/parse"
	bouncesEndpoint         = "bounces"
	statsEndpoint           = "stats"
	domainsEndpoint         = "domains"
	deleteTagEndpoint       = "tags"
	campaignsEndpoint       = "campaigns"
	eventsEndpoint		= "events"
	basicAuthUser           = "api"
)

// Mailgun defines the supported subset of the Mailgun API.
type Mailgun interface {
	Domain() string
	ApiKey() string
	PublicApiKey() string
	Send(m *Message) (string, string, error)
	ValidateEmail(email string) (EmailVerification, error)
	ParseAddresses(addresses ...string) ([]string, []string, error)
	GetBounces(limit, skip int) (int, []Bounce, error)
	GetSingleBounce(address string) (Bounce, error)
	AddBounce(address, code, error string) error
	DeleteBounce(address string) error
	GetStats(limit int, skip int, startDate *time.Time, event ...string) (int, []Stat, error)
	DeleteTag(tag string) error
	GetDomains(limit, skip int) (int, []Domain, error)
	GetSingleDomain(domain string) (Domain, []DNSRecord, []DNSRecord, error)
	CreateDomain(name string, smtpPassword string, spamAction string, wildcard bool) error
	DeleteDomain(name string) error
	GetCampaigns() (int, []Campaign, error)
	CreateCampaign(name, id string) error
	UpdateCampaign(oldId, name, newId string) error
	DeleteCampaign(id string) error
	GetComplaints(limit, skip int) (int, []Complaint, error)
	GetSingleComplaint(address string) (Complaint, error)
	GetStoredMessages() ([]StoredMessage, error)
	GetEvents(GetEventsOptions) ([]Event, Links, error)
}

// Imagine some data needed by a large set of methods in order to interact with the Mailgun API.
// mailgunImpl bundles these data together in a convenient place.
// Colloquially, we refer to instances of this structure as "clients."
type mailgunImpl struct {
	domain       string
	apiKey       string
	publicApiKey string
}

// Creates a new Mailgun client instance.
func NewMailgun(domain, apiKey, publicApiKey string) Mailgun {
	m := mailgunImpl{domain: domain, apiKey: apiKey, publicApiKey: publicApiKey}
	return &m
}

// Returns the domain configured for this client.
func (m *mailgunImpl) Domain() string {
	return m.domain
}

// Returns the API key configured for this client.
func (m *mailgunImpl) ApiKey() string {
	return m.apiKey
}

// Returns the public API key configured for this client.
func (m *mailgunImpl) PublicApiKey() string {
	return m.publicApiKey
}

// Generates the URL for the API using the domain and endpoint.
func generateApiUrl(m Mailgun, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", apiBase, m.Domain(), endpoint)
}

// Generates the URL for the API using the domain and endpoint, under the domains/ namespace.
func generateDomainUrl(m Mailgun, endpoint string) string {
	return fmt.Sprintf("%s/domains/%s/%s", apiBase, m.Domain(), endpoint)
}

// As with generateApiUrl, except that generatePublicApiUrl has no need for the domain.
func generatePublicApiUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", apiBase, endpoint)
}

func generateParameterizedUrl(m Mailgun, endpoint string, payload simplehttp.Payload) (string, error) {
	paramBuffer, err := payload.GetPayloadBuffer()
	if err != nil {
		return "", err
	}
	params := string(paramBuffer.Bytes())
	return fmt.Sprintf("%s?%s", generateApiUrl(m, eventsEndpoint), params), nil
}

// parseMailgunTime translates a timestamp as returned by Mailgun into a Go standard timestamp.
func parseMailgunTime(ts string) (t time.Time, err error) {
	t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 MST", ts)
	return
}

// formatMailgunTime translates a timestamp into a human-readable form.
func formatMailgunTime(t *time.Time) string {
	return t.Format("Mon, 2 Jan 2006 15:04:05 MST")
}
