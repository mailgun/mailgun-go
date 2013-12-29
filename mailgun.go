package mailgun

import (
	"errors"
	"fmt"
	"time"
	"github.com/mbanzon/simplehttp"
)

const (
	apiBase                 = "https://api.mailgun.net/v2"
	messagesEndpoint        = "messages"
	addressValidateEndpoint = "address/validate"
	addressParseEndpoint    = "address/parse"
	bouncesEndpoint         = "bounces"
	basicAuthUser           = "api"
)

type Mailgun interface {
	Domain() string
	ApiKey() string
	PublicApiKey() string
	SendMessage(m *MailgunMessage) (SendMessageResponse, error)
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

type SendMessageResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
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

func (m *mailgunImpl) SendMessage(message *MailgunMessage) (SendMessageResponse, error) {
	if !message.validateMessage() {
		return SendMessageResponse{}, errors.New("Message not valid")
	}

	r := simplehttp.NewSimpleHTTPRequest("POST", generateApiUrl(m, messagesEndpoint))
	r.AddFormValue("from", message.From.String())
	r.AddFormValue("subject", message.Subject)
	r.AddFormValue("text", message.Text)
	for _, to := range message.To {
		r.AddFormValue("to", to.String())
	}
	for _, cc := range message.Cc {
		r.AddFormValue("cc", cc.String())
	}
	for _, bcc := range message.Bcc {
		r.AddFormValue("bcc", bcc.String())
	}
	if message.Html != "" {
		r.AddFormValue("html", message.Html)
	}
	r.SetBasicAuth(basicAuthUser, m.ApiKey())

	var response SendMessageResponse
	err := r.MakeJSONRequest(&response)

	return response, err
}

func generateApiUrl(m Mailgun, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", apiBase, m.Domain(), endpoint)
}

func generatePublicApiUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", apiBase, endpoint)
}
