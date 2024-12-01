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
}

// plainMessage contains fields relevant to plain API-synthesized messages.
// You're expected to use various setters to set most of these attributes,
// although from, subject, and text are set when the message is created with
// NewMessage.
// TODO(v5): rename to PlainMessage
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
// TODO(v5): rename to MimeMessage
type mimeMessageV5 struct {
	commonMessageV5

	body io.ReadCloser
}

// features abstracts the common characteristics between regular and MIME messages.
type specificV5 interface {
	// AddCC appends a receiver to the carbon-copy header of a message.
	AddCC(string)

	// AddBCC appends a receiver to the blind-carbon-copy header of a message.
	AddBCC(string)

	// SetHTML If you're sending a message that isn't already MIME encoded, it will arrange to bundle
	// an HTML representation of your message in addition to your plain-text body.
	SetHTML(string)

	// SetAmpHTML If you're sending a message that isn't already MIME encoded, it will arrange to bundle
	// an AMP-For-Email representation of your message in addition to your html & plain-text content.
	SetAmpHTML(string)

	// AddValues invoked by Send() to add message-type-specific MIME headers for the API call
	// to Mailgun.
	AddValues(*formDataPayload)

	// Endpoint tells Send() which endpoint to use to submit the API call.
	Endpoint() string

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
	RecipientCount() int

	// SetTemplate sets the name of a template stored via the template API.
	// See https://documentation.mailgun.com/en/latest/user_manual.html#templating
	SetTemplate(string)

	// AddRecipient appends a receiver to the To: header of a message.
	// It will return an error if the limit of recipients have been exceeded for this message
	AddRecipient(recipient string)

	// isValid yields true if and only if the message is valid enough for sending
	// through the API.
	isValid() bool
}

// TODO(v5): implement for plain and MIME messages
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

	specificV5
}
