package mailgun

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
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

	data := generateUrlValues(message)

	req, err := http.NewRequest("POST", generateApiUrl(m, messagesEndpoint), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return SendMessageResponse{}, err
	}
	req.SetBasicAuth(basicAuthUser, m.ApiKey())
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return SendMessageResponse{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return SendMessageResponse{}, errors.New(fmt.Sprintf("Status is not 200. It was %d", resp.StatusCode))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return SendMessageResponse{}, err
	}

	var response SendMessageResponse
	err2 := json.Unmarshal(body, &response)
	if err2 != nil {
		return SendMessageResponse{}, err2
	}

	return response, nil
}

func generateUrlValues(message *MailgunMessage) url.Values {
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
	
	return data
}

func generateApiUrl(m Mailgun, endpoint string) string {
	return fmt.Sprintf("%s/%s/%s", apiBase, m.Domain(), endpoint)
}

func generatePublicApiUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", apiBase, endpoint)
}
