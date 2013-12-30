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
