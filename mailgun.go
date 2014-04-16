// Package mailgun provides methods for interacting with the Mailgun API.
// For further information please see the Mailgun documentation at
// http://documentation.mailgun.com/
//
// Author: Michael Banzon
package mailgun

import (
	"fmt"
	"github.com/mbanzon/simplehttp"
	"time"
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
	eventsEndpoint          = "events"
	credentialsEndpoint     = "credentials"
	unsubscribesEndpoint    = "unsubscribes"
	routesEndpoint          = "routes"
	webhooksEndpoint        = "webhooks"
	listsEndpoint           = "lists"
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
	GetStoredMessage(id string) (StoredMessage, error)
	DeleteStoredMessage(id string) error
	GetEvents(GetEventsOptions) ([]Event, Links, error)
	GetCredentials(limit, skip int) (int, []Credential, error)
	CreateCredential(login, password string) error
	ChangeCredentialPassword(id, password string) error
	DeleteCredential(id string) error
	GetUnsubscribes(limit, skip int) (int, []Unsubscription, error)
	GetUnsubscribesByAddress(string) (int, []Unsubscription, error)
	Unsubscribe(address, tag string) error
	RemoveUnsubscribe(string) error
	CreateComplaint(string) error
	DeleteComplaint(string) error
	GetRoutes(limit, skip int) (int, []Route, error)
	GetRouteByID(string) (Route, error)
	CreateRoute(Route) (Route, error)
	DeleteRoute(string) error
	UpdateRoute(string, Route) (Route, error)
	GetWebhooks() (map[string]string, error)
	CreateWebhook(kind, url string) error
	DeleteWebhook(kind string) error
	GetWebhookByType(kind string) (string, error)
	UpdateWebhook(kind, url string) error
	GetLists(limit, skip int, filter string) (int, []List, error)
	CreateList(List) (List, error)
	DeleteList(string) error
	GetListByAddress(string) (List, error)
	UpdateList(string, List) (List, error)
	GetMembers(limit, skip int, subfilter *bool, address string) (int, []Member, error)
	GetMemberByAddress(MemberAddr, listAddr string) (Member, error)
	CreateMember(merge bool, addr string, prototype Member) error
	CreateMemberList(subscribed *bool, addr string, newMembers []interface{}) error
	UpdateMember(Member, list string, prototype Member) (Member, error)
	DeleteMember(Member, list string) error
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

// Generates the URL to a mailing list subscriber endpoint.
func generateMemberApiUrl(endpoint, address string) string {
	return fmt.Sprintf("%s/%s/%s/members", apiBase, endpoint, address)
}

// Generates a targetted URL for the API using the domain and endpoint.
func generateApiUrlWithTarget(m Mailgun, endpoint, target string) string {
	tail := ""
	if target != "" {
		tail = fmt.Sprintf("/%s", target)
	}
	return fmt.Sprintf("%s%s", generateApiUrl(m, endpoint), tail)
}

func generateDomainApiUrl(m Mailgun, endpoint string) string {
	return fmt.Sprintf("%s/domains/%s/%s", apiBase, m.Domain(), endpoint)
}

func generateCredentialsUrl(m Mailgun, id string) string {
	tail := ""
	if id != "" {
		tail = fmt.Sprintf("/%s", id)
	}
	return generateDomainApiUrl(m, fmt.Sprintf("credentials%s", tail))
	// return fmt.Sprintf("%s/domains/%s/credentials%s", apiBase, m.Domain(), tail)
}

// Generates the URL needed to acquire a copy of a stored message.
func generateStoredMessageUrl(m Mailgun, endpoint, id string) string {
	return generateDomainApiUrl(m, fmt.Sprintf("%s/%s", endpoint, id))
	// return fmt.Sprintf("%s/domains/%s/%s/%s", apiBase, m.Domain(), endpoint, id)
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
