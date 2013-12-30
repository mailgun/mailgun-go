package mailgun

import (
	"errors"
	"github.com/mbanzon/simplehttp"
)

type Message struct {
	from    string
	To      []string
	Cc      []string
	Bcc     []string
	subject string
	text    string
	Html    string
}

type sendMessageResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

func NewMessage(from, subject, text string) *Message {
	return &Message{ from: from, subject: subject, text: text }
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
		r := simplehttp.NewPostRequest(generateApiUrl(m, messagesEndpoint))
		r.AddFormValue("from", message.from)
		r.AddFormValue("subject", message.subject)
		r.AddFormValue("text", message.text)
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

	if m.from == "" {
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

	if m.text == "" {
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
