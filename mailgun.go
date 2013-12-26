package mailgun

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
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
	if !message.validateMessage() {
		return errors.New("Message not valid")
	}

	req, err := http.NewRequest("POST", generateApiUrl(m, messagesEndpoint), nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(basicAuthUser, m.ApiKey())

	data := url.Values{}
	data.Add("from", message.From.String())
	data.Add("subject", message.Subject)
	data.Add("text", message.Text)

	for _, to := range message.To {
		data.Add("to", to.String())
	}
	for _, cc := range message.Cc {
		data.Add("cc", cc.String())
	}
	for _, bcc := range message.Bcc {
		data.Add("bcc", bcc.String())
	}

	if message.Html != "" {
		data.Add("html", message.Html)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New("Status is not 200")
	}

	return nil
}

func generateApiUrl(m Mailgun, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", apiBase, m.Domain(), endpoint)
}
