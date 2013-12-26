package mailgun

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"bytes"
)

const (
	apiBase          = "https://api.mailgun.net/v2"
	messagesEndpoint = "messages"
	basicAuthUser    = "api"
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

	req, err := http.NewRequest("POST", generateApiUrl(m, messagesEndpoint), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(basicAuthUser, m.ApiKey())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Status is not 200. It was %d", resp.StatusCode))
	}

	return nil
}

func generateApiUrl(m Mailgun, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", apiBase, m.Domain(), endpoint)
}
