package mailgun

import (
	"io"
	"time"
)

// Message structures contain both the message text and the envelope for an e-mail message.
// TODO(v5): remove v5 from the name
type commonMessageV5 struct {
	domain             string
	to                 []string
	tags               []string
	campaigns          []string
	dkim               *bool
	deliveryTime       time.Time
	stoPeriod          string
	attachments        []string
	readerAttachments  []ReaderAttachment
	inlines            []string
	readerInlines      []ReaderAttachment
	bufferAttachments  []BufferAttachment
	nativeSend         bool
	testMode           bool
	tracking           *bool
	trackingClicks     *string
	trackingOpens      *bool
	headers            map[string]string
	variables          map[string]string
	templateVariables  map[string]interface{}
	recipientVariables map[string]map[string]interface{}
	templateVersionTag string
	templateRenderText bool
	requireTLS         bool
	skipVerification   bool

	// specific featuresV5
}

// plainMessage contains fields relevant to plain API-synthesized messages.
// You're expected to use various setters to set most of these attributes,
// although from, subject, and text are set when the message is created with
// NewMessage.
type plainMessageV5 struct {
	commonMessageV5

	from     string
	cc       []string
	bcc      []string
	subject  string
	text     string
	html     string
	ampHtml  string
	template string
}

// mimeMessage contains fields relevant to pre-packaged MIME messages.
type mimeMessageV5 struct {
	commonMessageV5

	body io.ReadCloser
}

// features abstracts the common characteristics between regular and MIME messages.
// addCC, addBCC, recipientCount, setHtml and setAMPHtml are invoked via the AddCC, AddBCC,
// RecipientCount, SetHTML and SetAMPHtml calls, as these functions are ignored for MIME messages.
// Send() invokes addValues to add message-type-specific MIME headers for the API call
// to Mailgun.
// isValid yields true if and only if the message is valid enough for sending
// through the API.
// Finally, endpoint() tells Send() which endpoint to use to submit the API call.
// TODO(v5): remove?
type featuresV5 interface {
	// AddCC appends a receiver to the carbon-copy header of a message.
	addCC(string)

	addBCC(string)

	setHtml(string)

	setAMPHtml(string)

	addValues(*formDataPayload)

	isValid() bool

	endpoint() string

	// RecipientCount returns the total number of recipients for the message.
	// This includes To:, Cc:, and Bcc: fields.
	//
	// NOTE: At present, this method is reliable only for non-MIME messages, as the
	// Bcc: and Cc: fields are easily accessible.
	// For MIME messages, only the To: field is considered.
	// A fix for this issue is planned for a future release.
	// For now, MIME messages are always assumed to have 10 recipients between Cc: and Bcc: fields.
	// If your MIME messages have more than 10 non-To: field recipients,
	// you may find that some recipients will not receive your e-mail.
	// It's perfectly OK, of course, for a MIME message to not have any Cc: or Bcc: recipients.
	recipientCount() int

	setTemplate(string)
}

type messageIfaceV5 interface {
	Domain() string
	To() []string
	Tags() []string
	Campaigns() []string
	DKIM() *bool
	DeliveryTime() time.Time
	STOPeriod() string
	Attachments() []string
	ReaderAttachments() []ReaderAttachment
	Inlines() []string
	ReaderInlines() []ReaderAttachment
	BufferAttachments() []BufferAttachment
	NativeSend() bool
	TestMode() bool
	Tracking() *bool
	TrackingClicks() *string
	TrackingOpens() *bool
	Headers() map[string]string
	Variables() map[string]string
	TemplateVariables() map[string]interface{}
	RecipientVariables() map[string]map[string]interface{}
	TemplateVersionTag() string
	TemplateRenderText() bool
	RequireTLS() bool
	SkipVerification() bool

	RecipientCount() int
	AddValues(p *formDataPayload)

	featuresV5
}
