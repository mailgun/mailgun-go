package mailgun

import (
	"encoding/json"
	"errors"
	"github.com/mbanzon/simplehttp"
	"time"
)

type Message struct {
	from         string
	to           []string
	cc           []string
	bcc          []string
	subject      string
	text         string
	html         string
	tags         []string
	campaigns    []string
	dkim         bool
	deliveryTime *time.Time

	testMode       bool
	tracking       bool
	trackingClicks bool
	trackingOpens  bool
	headers        map[string]string
	variables      map[string]string

	dkimSet           bool
	trackingSet       bool
	trackingClicksSet bool
	trackingOpensSet  bool
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

func (m *Message) AddCampaign(campaign string) {
	m.campaigns = append(m.campaigns, campaign)
}

func (m *Message) SetDKIM(dkim bool) {
	m.dkim = dkim
	m.dkimSet = true
}

func (m *Message) EnableTestMode() {
	m.testMode = true
}

func (m *Message) SetTracking(tracking bool) {
	m.tracking = tracking
	m.trackingSet = true
}

func (m *Message) SetTrackingClicks(trackingClicks bool) {
	m.trackingClicks = trackingClicks
	m.trackingClicksSet = true
}

func (m *Message) SetTrackingOpens(trackingOpens bool) {
	m.trackingOpens = trackingOpens
	m.trackingOpensSet = true
}

func (m *Message) AddHeader(header, value string) {
	if m.headers == nil {
		m.headers = make(map[string]string)
	}
	m.headers[header] = value
}

func (m *Message) AddVariable(variable string, value interface{}) error {
	j, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if m.variables == nil {
		m.variables = make(map[string]string)
	}
	m.variables[variable] = string(j)
	return nil
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
		for _, campaign := range message.campaigns {
			r.AddFormValue("o:campaign", campaign)
		}
		if message.html != "" {
			r.AddFormValue("html", message.html)
		}
		if message.dkimSet {
			r.AddFormValue("o:dkim", yesNo(message.dkim))
		}
		if message.deliveryTime != nil {
			r.AddFormValue("o:deliverytime", message.deliveryTime.Format("Mon, 2 Jan 2006 15:04:05 MST"))
		}
		if message.testMode {
			r.AddFormValue("o:testmode", "yes")
		}
		if message.trackingSet {
			r.AddFormValue("o:tracking", yesNo(message.tracking))
		}
		if message.trackingClicksSet {
			r.AddFormValue("o:tracking-clicks", yesNo(message.trackingClicks))
		}
		if message.trackingOpensSet {
			r.AddFormValue("o:tracking-opens", yesNo(message.trackingOpens))
		}
		if message.headers != nil {
			for header, value := range message.headers {
				r.AddFormValue("h:"+header, value)
			}
		}
		if message.variables != nil {
			for variable, value := range message.variables {
				r.AddFormValue("v:"+variable, value)
			}
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

func yesNo(b bool) string {
	if b {
		return "yes"
	} else {
		return "no"
	}
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
