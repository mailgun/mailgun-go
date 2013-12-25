package mailgun

import (
	"fmt"
	"net/http"
)

const (
	apiBase          = "https://api.mailgun.net/v2"
	messagesEndpoint = "messages"
	basicAuthUser    = "key"
)

type Mailgun interface {
	Domain() string
	ApiKey() string
	SendMessage(m *MailgunMessage) error
}

type mailgunImpl struct {
	domain string
	apiKey string
}

func NewMailgun(domain, apiKey string) Mailgun {
	m := mailgunImpl{domain: domain, apiKey: apiKey}
	return &m
}

func (m *mailgunImpl) Domain() string {
	return m.domain
}

func (m *mailgunImpl) ApiKey() string {
	return m.apiKey
}

func (m *mailgunImpl) SendMessage(message *MailgunMessage) error {
	// TODO
	req, err := http.NewRequest("POST", generateApiUrl(m, messagesEndpoint), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(basicAuthUser, m.ApiKey())

	return nil
}

func generateApiUrl(m Mailgun, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", apiBase, m.Domain(), endpoint)
}
