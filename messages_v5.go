package mailgun

// This file contains methods for new v5 messages.

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// NewMessage returns a new e-mail message with the simplest envelop needed to send.
//
// Supports arbitrary-sized recipient lists by
// automatically sending mail in batches of up to MaxNumberOfRecipients.
//
// To support batch sending, do not provide `to` at this point.
// You can do this explicitly, or implicitly, as follows:
//
//	// Note absence of `to` parameter(s)!
//	m := NewMessage("me@example.com", "Help save our planet", "Hello world!")
//
// Note that you'll need to invoke the AddRecipientAndVariables or AddRecipient method
// before sending, though.
// TODO(v5): rename to NewMessage
func newMessageV5(domain, from, subject, text string, to ...string) *plainMessageV5 {
	return &plainMessageV5{
		commonMessageV5: commonMessageV5{
			domain: domain,
			to:     to,
		},

		from:    from,
		subject: subject,
		text:    text,
	}
}

// NewMIMEMessage creates a new MIME message. These messages are largely canned;
// you do not need to invoke setters to set message-related headers.
// However, you do still need to call setters for Mailgun-specific settings.
//
// Supports arbitrary-sized recipient lists by
// automatically sending mail in batches of up to MaxNumberOfRecipients.
//
// To support batch sending, do not provide `to` at this point.
// You can do this explicitly, or implicitly, as follows:
//
//	// Note absence of `to` parameter(s)!
//	m := NewMIMEMessage(body)
//
// Note that you'll need to invoke the AddRecipientAndVariables or AddRecipient method
// before sending, though.
// TODO(v5): rename to NewMIMEMessage
func newMIMEMessage(domain string, body io.ReadCloser, to ...string) *mimeMessageV5 {
	return &mimeMessageV5{
		commonMessageV5: commonMessageV5{
			domain: domain,
			to:     to,
		},
		body: body,
	}
}

// AddReaderAttachment arranges to send a file along with the e-mail message.
// File contents are read from a io.ReadCloser.
// The filename parameter is the resulting filename of the attachment.
// The readCloser parameter is the io.ReadCloser which reads the actual bytes to be used
// as the contents of the attached file.
func (m *commonMessageV5) AddReaderAttachment(filename string, readCloser io.ReadCloser) {
	ra := ReaderAttachment{Filename: filename, ReadCloser: readCloser}
	m.readerAttachments = append(m.readerAttachments, ra)
}

// AddBufferAttachment arranges to send a file along with the e-mail message.
// File contents are read from the []byte array provided
// The filename parameter is the resulting filename of the attachment.
// The buffer parameter is the []byte array which contains the actual bytes to be used
// as the contents of the attached file.
func (m *commonMessageV5) AddBufferAttachment(filename string, buffer []byte) {
	ba := BufferAttachment{Filename: filename, Buffer: buffer}
	m.bufferAttachments = append(m.bufferAttachments, ba)
}

// AddAttachment arranges to send a file from the filesystem along with the e-mail message.
// The attachment parameter is a filename, which must refer to a file which actually resides
// in the local filesystem.
func (m *commonMessageV5) AddAttachment(attachment string) {
	m.attachments = append(m.attachments, attachment)
}

// AddReaderInline arranges to send a file along with the e-mail message.
// File contents are read from a io.ReadCloser.
// The filename parameter is the resulting filename of the attachment.
// The readCloser parameter is the io.ReadCloser which reads the actual bytes to be used
// as the contents of the attached file.
func (m *commonMessageV5) AddReaderInline(filename string, readCloser io.ReadCloser) {
	ra := ReaderAttachment{Filename: filename, ReadCloser: readCloser}
	m.readerInlines = append(m.readerInlines, ra)
}

// AddInline arranges to send a file along with the e-mail message, but does so
// in a way that its data remains "inline" with the rest of the message.  This
// can be used to send image or font data along with an HTML-encoded message body.
// The attachment parameter is a filename, which must refer to a file which actually resides
// in the local filesystem.
func (m *commonMessageV5) AddInline(inline string) {
	m.inlines = append(m.inlines, inline)
}

// SetDKIM arranges to send the o:dkim header with the message, and sets its value accordingly.
// Refer to the Mailgun documentation for more information.
func (m *commonMessageV5) SetDKIM(dkim bool) {
	m.dkim = &dkim
}

// EnableNativeSend allows the return path to match the address in the Message.Headers.From:
// field when sending from Mailgun rather than the usual bounce+ address in the return path.
func (m *commonMessageV5) EnableNativeSend() {
	m.nativeSend = true
}

// EnableTestMode allows submittal of a message, such that it will be discarded by Mailgun.
// This facilitates testing client-side software without actually consuming e-mail resources.
func (m *commonMessageV5) EnableTestMode() {
	m.testMode = true
}

// SetDeliveryTime schedules the message for transmission at the indicated time.
// Pass nil to remove any installed schedule.
// Refer to the Mailgun documentation for more information.
func (m *commonMessageV5) SetDeliveryTime(dt time.Time) {
	m.deliveryTime = dt
}

// SetSTOPeriod toggles Send Time Optimization (STO) on a per-message basis.
// String should be set to the number of hours in [0-9]+h format,
// with the minimum being 24h and the maximum being 72h.
// Refer to the Mailgun documentation for more information.
func (m *commonMessageV5) SetSTOPeriod(stoPeriod string) error {
	validPattern := `^([2-6][4-9]|[3-6][0-9]|7[0-2])h$`
	// TODO(vtopc): regexp.Compile, which is called by regexp.MatchString, is a heave operation, move into global variable
	// or just parse using time.ParseDuration().
	match, err := regexp.MatchString(validPattern, stoPeriod)
	if err != nil {
		return err
	}

	if !match {
		return errors.New("STO period is invalid. Valid range is 24h to 72h")
	}

	m.stoPeriod = stoPeriod
	return nil
}

// SetTracking sets the o:tracking message parameter to adjust, on a message-by-message basis,
// whether or not Mailgun will rewrite URLs to facilitate event tracking.
// Events tracked includes opens, clicks, unsubscribes, etc.
// Note: simply calling this method ensures that the o:tracking header is passed in with the message.
// Its yes/no setting is determined by the call's parameter.
// Note that this header is not passed on to the final recipient(s).
// Refer to the Mailgun documentation for more information.
func (m *commonMessageV5) SetTracking(tracking bool) {
	m.tracking = &tracking
}

// SetTrackingClicks information is found in the Mailgun documentation.
func (m *commonMessageV5) SetTrackingClicks(trackingClicks bool) {
	m.trackingClicks = ptr(yesNo(trackingClicks))
}

// SetTrackingOptions sets the o:tracking, o:tracking-clicks and o:tracking-opens at once.
func (m *commonMessageV5) SetTrackingOptions(options *TrackingOptions) {
	m.tracking = &options.Tracking

	m.trackingClicks = &options.TrackingClicks

	m.trackingOpens = &options.TrackingOpens
}

// SetRequireTLS information is found in the Mailgun documentation.
func (m *commonMessageV5) SetRequireTLS(b bool) {
	m.requireTLS = b
}

// SetSkipVerification information is found in the Mailgun documentation.
func (m *commonMessageV5) SetSkipVerification(b bool) {
	m.skipVerification = b
}

// SetTrackingOpens information is found in the Mailgun documentation.
func (m *commonMessageV5) SetTrackingOpens(trackingOpens bool) {
	m.trackingOpens = &trackingOpens
}

// SetTemplateVersion information is found in the Mailgun documentation.
func (m *commonMessageV5) SetTemplateVersion(tag string) {
	m.templateVersionTag = tag
}

// SetTemplateRenderText information is found in the Mailgun documentation.
func (m *commonMessageV5) SetTemplateRenderText(render bool) {
	m.templateRenderText = render
}

// AddHeader allows you to send custom MIME headers with the message.
func (m *commonMessageV5) AddHeader(header, value string) {
	if m.headers == nil {
		m.headers = make(map[string]string)
	}
	m.headers[header] = value
}

// AddVariable lets you associate a set of variables with messages you send,
// which Mailgun can use to, in essence, complete form-mail.
// Refer to the Mailgun documentation for more information.
func (m *commonMessageV5) AddVariable(variable string, value any) error {
	if m.variables == nil {
		m.variables = make(map[string]string)
	}

	j, err := json.Marshal(value)
	if err != nil {
		return err
	}

	encoded := string(j)
	v, err := strconv.Unquote(encoded)
	if err != nil {
		v = encoded
	}

	m.variables[variable] = v
	return nil
}

// AddTemplateVariable adds a template variable to the map of template variables, replacing the variable if it is already there.
// This is used for server-side message templates and can nest arbitrary values. At send time, the resulting map will be converted into
// a JSON string and sent as a header in the X-Mailgun-Variables header.
func (m *commonMessageV5) AddTemplateVariable(variable string, value any) error {
	if m.templateVariables == nil {
		m.templateVariables = make(map[string]any)
	}
	m.templateVariables[variable] = value
	return nil
}

// AddDomain allows you to use a separate domain for the type of messages you are sending.
func (m *commonMessageV5) AddDomain(domain string) {
	m.domain = domain
}

// GetHeaders retrieves the http headers associated with this message
func (m *commonMessageV5) GetHeaders() map[string]string {
	return m.headers
}

// specific message methods

func (m *plainMessageV5) AddRecipient(recipient string) error {
	return m.AddRecipientAndVariables(recipient, nil)
}

// AddRecipientAndVariables appends a receiver to the To: header of a message,
// and as well attaches a set of variables relevant for this recipient.
// It will return an error if the limit of recipients have been exceeded for this message
func (m *plainMessageV5) AddRecipientAndVariables(r string, vars map[string]any) error {
	if m.RecipientCount() >= MaxNumberOfRecipients {
		return fmt.Errorf("recipient limit exceeded (max %d)", MaxNumberOfRecipients)
	}
	m.to = append(m.to, r)
	if vars != nil {
		if m.recipientVariables == nil {
			m.recipientVariables = make(map[string]map[string]any)
		}
		m.recipientVariables[r] = vars
	}

	return nil
}

func (m *mimeMessageV5) AddRecipient(recipient string) error {
	if m.RecipientCount() >= MaxNumberOfRecipients {
		return fmt.Errorf("recipient limit exceeded (max %d)", MaxNumberOfRecipients)
	}
	m.to = append(m.to, recipient)

	return nil
}

func (m *plainMessageV5) RecipientCount() int {
	return len(m.to) + len(m.bcc) + len(m.cc)
}

func (m *mimeMessageV5) RecipientCount() int {
	return len(m.to)
}

// SetReplyTo sets the receiver who should receive replies
func (m *commonMessageV5) SetReplyTo(recipient string) {
	m.AddHeader("Reply-To", recipient)
}

func (m *plainMessageV5) AddCC(r string) {
	m.cc = append(m.cc, r)
}

func (*mimeMessageV5) AddCC(_ string) {}

func (m *plainMessageV5) AddBCC(r string) {
	m.bcc = append(m.bcc, r)
}

func (*mimeMessageV5) AddBCC(_ string) {}

func (m *plainMessageV5) SetHTML(h string) {
	m.html = h
}

func (*mimeMessageV5) SetHTML(_ string) {}

func (m *plainMessageV5) SetAmpHTML(h string) {
	m.ampHtml = h
}

func (*mimeMessageV5) SetAmpHTML(_ string) {}

// AddTag attaches tags to the message.  Tags are useful for metrics gathering and event tracking purposes.
// Refer to the Mailgun documentation for further details.
func (m *commonMessageV5) AddTag(tag ...string) error {
	if len(m.tags) >= MaxNumberOfTags {
		return fmt.Errorf("cannot add any new tags. Message tag limit (%d) reached", MaxNumberOfTags)
	}

	m.tags = append(m.tags, tag...)
	return nil
}

func (m *plainMessageV5) SetTemplate(t string) {
	m.template = t
}

func (*mimeMessageV5) SetTemplate(_ string) {}

func (m *plainMessageV5) AddValues(p *FormDataPayload) {
	p.addValue("from", m.from)
	p.addValue("subject", m.subject)
	p.addValue("text", m.text)
	for _, cc := range m.cc {
		p.addValue("cc", cc)
	}
	for _, bcc := range m.bcc {
		p.addValue("bcc", bcc)
	}
	if m.html != "" {
		p.addValue("html", m.html)
	}
	if m.template != "" {
		p.addValue("template", m.template)
	}
	if m.ampHtml != "" {
		p.addValue("amp-html", m.ampHtml)
	}
}

func (m *mimeMessageV5) AddValues(p *FormDataPayload) {
	p.addReadCloser("message", "message.mime", m.body)
}

func (*plainMessageV5) Endpoint() string {
	return messagesEndpoint
}

func (*mimeMessageV5) Endpoint() string {
	return mimeMessagesEndpoint
}

// Send attempts to queue a message (see Message, NewMessage, and its methods) for delivery.
// It returns the Mailgun server response, which consists of two components:
//   - A human-readable status message, typically "Queued. Thank you."
//   - A Message ID, which is the id used to track the queued message. The message id is useful
//     when contacting support to report an issue with a specific message or to relate a
//     delivered, accepted or failed event back to specific message.
//
// The status and message ID are only returned if no error occurred.
//
// Error returns can be of type `error.Error` which wrap internal and standard
// golang errors like `url.Error`. The error can also be of type
// mailgun.UnexpectedResponseError which contains the error returned by the mailgun API.
//
//	mailgun.UnexpectedResponseError {
//	  URL:      "https://api.mailgun.com/v3/messages",
//	  Expected: 200,
//	  Actual:   400,
//	  Data:     "Domain not found: example.com",
//	}
//
//	See the public mailgun documentation for all possible return codes and error messages
func (mg *MailgunImpl) sendV5(ctx context.Context, m SendableMessage) (mes, id string, err error) {
	// TODO(vtopc): move domain checks into NewMessage and NewMIMEMessage?
	if m.Domain() == "" {
		err = errors.New("you must provide a valid domain before calling Send()")
		return "", "", err
	}

	invalidChars := ":&'@(),!?#;%+=<>"
	if i := strings.ContainsAny(m.Domain(), invalidChars); i {
		err = fmt.Errorf("you called Send() with a domain that contains invalid characters")
		return "", "", err
	}

	if mg.apiKey == "" {
		err = errors.New("you must provide a valid api-key before calling Send()")
		return "", "", err
	}

	if !isValid(m) {
		err = ErrInvalidMessage
		return "", "", err
	}

	if m.STOPeriod() != "" && m.RecipientCount() > 1 {
		err = errors.New("STO can only be used on a per-message basis")
		return "", "", err
	}
	payload := NewFormDataPayload()

	m.AddValues(payload)

	// TODO: make (CommonMessage).AddValues():
	err = addMessageValues(payload, m)
	if err != nil {
		return "", "", err
	}

	r := newHTTPRequest(generateApiUrlWithDomain(mg, m.Endpoint(), m.Domain()))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	// Override any HTTP headers if provided
	for k, v := range mg.overrideHeaders {
		r.addHeader(k, v)
	}

	var response sendMessageResponse
	err = postResponseFromJSON(ctx, r, payload, &response)
	if err == nil {
		mes = response.Message
		id = response.Id
	}

	return mes, id, err
}

func (m *plainMessageV5) isValid() bool {
	if !validateStringList(m.cc, false) {
		return false
	}

	if !validateStringList(m.bcc, false) {
		return false
	}

	if m.from == "" {
		return false
	}

	if m.template != "" {
		// m.text or m.html not needed if template is supplied
		return true
	}

	if m.text == "" && m.html == "" {
		return false
	}

	return true
}

func (m *mimeMessageV5) isValid() bool {
	return m.body != nil
}
