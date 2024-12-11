package mailgun

// This file contains a draft for new v5 messages.

import (
	"io"
	"time"
)

// Message structures contain both the message text and the envelope for an e-mail message.
// TODO(v5): rename to CommonMessage
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
	templateVariables  map[string]any
	recipientVariables map[string]map[string]any
	templateVersionTag string
	templateRenderText bool
	requireTLS         bool
	skipVerification   bool
}

// PlainMessage contains fields relevant to plain API-synthesized messages.
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

// MimeMessage contains fields relevant to pre-packaged MIME messages.
// TODO(v5): rename to MimeMessage
type mimeMessageV5 struct {
	commonMessageV5

	body io.ReadCloser
}
