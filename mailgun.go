package mailgun

const (
	MAILGUN_API_ADDRESS = "https://api.mailgun.net/v2/samples.mailgun.org/messages"
)

type Mailgun interface {
	Domain() string
	ApiKey() string
	SendMessage(m *MailgunMessage)
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

func (m *mailgunImpl) SendMessage(message *MailgunMessage) {
	// TODO
}
