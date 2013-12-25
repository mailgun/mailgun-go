package mailgun

type Mailgun interface {
	ApiKey() string
	SendMessage(m *MailgunMessage)
}

type mailgunImpl struct {
	apiKey string
}

func NewMailgun(apiKey string) Mailgun {
	m := mailgunImpl{apiKey: apiKey}
	return &m
}

func (m *mailgunImpl) ApiKey() string {
	return m.apiKey
}

func (m *mailgunImpl) SendMessage(message *MailgunMessage) {
	// TODO
}
