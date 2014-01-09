// Package mailgun provides methods for interacting with the Mailgun API.
// For further information please see the Mailgun documentation at
// http://documentation.mailgun.com/
//
// Author: Michael Banzon
package mailgun

import (
	"fmt"
	"time"
)

const (
	apiBase                 = "https://api.mailgun.net/v2"
	messagesEndpoint        = "messages"
	addressValidateEndpoint = "address/validate"
	addressParseEndpoint    = "address/parse"
	bouncesEndpoint         = "bounces"
	statsEndpoint           = "stats"
	domainsEndpoint         = "domains"
	deleteTagEndpoint       = "tags"
	campaignsEndpoint       = "campaigns"
	basicAuthUser           = "api"
)

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
	GetStats(limit int, skip int, startDate time.Time, event ...string) (int, []Stat, error)
	DeleteTag(tag string) error
	GetDomains(limit, skip int) (int, []Domain, error)
	GetSingleDomain(domain string) (Domain, []DNSRecord, []DNSRecord, error)
	CreateDomain(name string, smtpPassword string, spamAction bool, wildcard bool) error
	DeleteDomain(name string) error
	GetCampaigns() (int, []Campaign, error)
	CreateCampaign(name, id string) error
	UpdateCampaign(oldId, name, newId string) error
	DeleteCampaign(id string) error
	GetComplaints(limit, skip int) (int, []Complaint, error)
	GetSingleComplaint(address string) (Complaint, error)
}

type mailgunImpl struct {
	domain       string
	apiKey       string
	publicApiKey string
}

func NewMailgun(domain, apiKey, publicApiKey string) Mailgun {
	m := mailgunImpl{domain: domain, apiKey: apiKey, publicApiKey: publicApiKey}
	return &m
}

func (m *mailgunImpl) Domain() string {
	return m.domain
}

func (m *mailgunImpl) ApiKey() string {
	return m.apiKey
}

func (m *mailgunImpl) PublicApiKey() string {
	return m.publicApiKey
}

func generateApiUrl(m Mailgun, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", apiBase, m.Domain(), endpoint)
}

func generatePublicApiUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", apiBase, endpoint)
}

func parseMailgunTime(ts string) (t time.Time, err error) {
	t, err = time.Parse("Mon, 2 Jan 2006 15:04:05 MST", ts)
	return
}
