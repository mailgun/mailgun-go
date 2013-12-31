package mailgun

import (
	"errors"
	"github.com/mbanzon/simplehttp"
)

type Message struct {
	from     string
	to       []string
	cc       []string
	bcc      []string
	subject  string
	text     string
	html     string
	tags     []string
	campaigns []string
}

type sendMessageResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

func NewMessage(from string, subject string, text string, to ...string) *Message {
	return &Message{from: from, subject: subject, text: text, to: to}
}

func (m *Message) AddRecipient(recipient string) {
	m.to = append(m.to, recipient)
}

func (m *Message) AddCC(recipient string) {
	m.cc = append(m.cc, recipient)
}

func (m *Message) AddBCC(recipient string) {
	m.bcc = append(m.bcc, recipient)
}

func (m *Message) SetHtml(html string) {
	m.html = html
}
func (m *Message) AddTag(tag string) {
	m.tags = append(m.tags, tag)
}

func (m *Message) AddCampaign(campain string) {
	m.campaigns = append(m.campaigns, campain)
}

func (m *mailgunImpl) Send(message *Message) (mes string, id string, err error) {
	if !message.validateMessage() {
		err = errors.New("Message not valid")
	} else {
		r := simplehttp.NewPostRequest(generateApiUrl(m, messagesEndpoint))
		r.AddFormValue("from", message.from)
		r.AddFormValue("subject", message.subject)
		r.AddFormValue("text", message.text)
		for _, to := range message.to {
			r.AddFormValue("to", to)
		}
		for _, cc := range message.cc {
			r.AddFormValue("cc", cc)
		}
		for _, bcc := range message.bcc {
			r.AddFormValue("bcc", bcc)
		}
		for _, tag := range message.tags {
			r.AddFormValue("o:tag", tag)
		}
		for _, campain := range message.campaigns {
			r.AddFormValue("o:campain", campain)
		}
		if message.html != "" {
			r.AddFormValue("html", message.html)
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

	if !validateStringList(m.to, true) {
		return false
	}

	if !validateStringList(m.cc, false) {
		return false
	}

	if !validateStringList(m.bcc, false) {
		return false
	}

	if !validateStringList(m.tags, false) {
		return false
	}

	if !validateStringList(m.campaigns, false) || len(m.campaigns) > 3 {
		return false
	}

	if m.text == "" {
		return false
	}

	return true
}

func validateStringList(list []string, requireOne bool) bool {
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
