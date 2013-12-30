package mailgun

import (
	"errors"
	"github.com/mbanzon/simplehttp"
)

type Message struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Text    string
	Html    string
}

type sendMessageResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

func (m *Message) AddRecipient(recipient string) {
	m.To = append(m.To, recipient)
}

func (m *Message) AddCC(recipient string) {
	m.Cc = append(m.Cc, recipient)
}

func (m *Message) AddBCC(recipient string) {
	m.Bcc = append(m.Bcc, recipient)
}

func (m *mailgunImpl) Send(message *Message) (mes string, id string, err error) {
	if !message.validateMessage() {
		err = errors.New("Message not valid")
	} else {
		r := simplehttp.NewSimpleHTTPRequest("POST", generateApiUrl(m, messagesEndpoint))
		r.AddFormValue("from", message.From)
		r.AddFormValue("subject", message.Subject)
		r.AddFormValue("text", message.Text)
		for _, to := range message.To {
			r.AddFormValue("to", to)
		}
		for _, cc := range message.Cc {
			r.AddFormValue("cc", cc)
		}
		for _, bcc := range message.Bcc {
			r.AddFormValue("bcc", bcc)
		}
		if message.Html != "" {
			r.AddFormValue("html", message.Html)
		}
		r.SetBasicAuth(basicAuthUser, m.ApiKey())

		var response sendMessageResponse
		err = r.MakeJSONRequest(&response)
		if err == nil {
			mes = response.Message
			id = response.Id
		}
	}

	return
}

func (m *Message) validateMessage() bool {
	if m == nil {
		return false
	}

	if m.From == "" {
		return false
	}

	if !validateAddressList(m.To, true) {
		return false
	}

	if !validateAddressList(m.Cc, false) {
		return false
	}

	if !validateAddressList(m.Bcc, false) {
		return false
	}

	if m.Text == "" {
		return false
	}

	return true
}

func validateAddressList(list []string, requireOne bool) bool {
	hasOne := false

	if list == nil {
		return !requireOne
	} else {
		for _, a := range list {
			if a == "" {
				return false
			} else {
				hasOne = hasOne || true
			}
		}
	}

	return hasOne
}
