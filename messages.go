package mailgun

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

// MaxNumberOfRecipients represents the largest batch of recipients that Mailgun can support in a single API call.
// This figure includes To:, Cc:, Bcc:, etc. recipients.
const MaxNumberOfRecipients = 1000

// MaxNumberOfTags represents the maximum number of tags that can be added for a message
const MaxNumberOfTags = 3

// Message structures contain both the message text and the envelope for an e-mail message.
type Message struct {
	to                []string
	tags              []string
	campaigns         []string
	dkim              *bool
	deliveryTime      time.Time
	stoPeriod         string
	attachments       []string
	readerAttachments []ReaderAttachment
	inlines           []string
	readerInlines     []ReaderAttachment
	bufferAttachments []BufferAttachment

	nativeSend         bool
	testMode           bool
	tracking           *bool
	trackingClicks     *string
	trackingOpens      *bool
	headers            map[string]string
	variables          map[string]string
	templateVariables  map[string]interface{}
	recipientVariables map[string]map[string]interface{}
	domain             string
	templateVersionTag string
	templateRenderText bool

	requireTLS       bool
	skipVerification bool

	Specific
}

type ReaderAttachment struct {
	Filename   string
	ReadCloser io.ReadCloser
}

type BufferAttachment struct {
	Filename string
	Buffer   []byte
}

// plainMessage contains fields relevant to plain API-synthesized messages.
// You're expected to use various setters to set most of these attributes,
// although from, subject, and text are set when the message is created with
// NewMessage.
type plainMessage struct {
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
type mimeMessage struct {
	body io.ReadCloser
}

type sendMessageResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}

// TrackingOptions contains fields relevant to trackings.
type TrackingOptions struct {
	Tracking       bool
	TrackingClicks string
	TrackingOpens  bool
}

// Specific abstracts the common characteristics between regular and MIME messages.
type Specific interface {
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
	AddValues(*FormDataPayload)

	// IsValid yields true if and only if the message is valid enough for sending
	// through the API.
	IsValid() bool

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
	// See https://documentation.mailgun.com/docs/mailgun/user-manual/sending-messages/#templates
	SetTemplate(string)

	// // TODO(v5):
	// // AddRecipient appends a receiver to the To: header of a message.
	// // It will return an error if the limit of recipients have been exceeded for this message
	// AddRecipient(recipient string) error
}

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
func NewMessage(from, subject, text string, to ...string) *Message {
	return &Message{
		Specific: &plainMessage{
			from:    from,
			subject: subject,
			text:    text,
		},
		to: to,
	}
}

// Deprecated: use func NewMessage instead of method.
//
// TODO(v5): remove this method
func (*MailgunImpl) NewMessage(from, subject, text string, to ...string) *Message {
	return NewMessage(from, subject, text, to...)
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
func NewMIMEMessage(body io.ReadCloser, to ...string) *Message {
	return &Message{
		Specific: &mimeMessage{
			body: body,
		},
		to: to,
	}
}

// Deprecated: use func NewMIMEMessage instead of method.
//
// TODO(v5): remove this method
func (*MailgunImpl) NewMIMEMessage(body io.ReadCloser, to ...string) *Message {
	return NewMIMEMessage(body, to...)
}

func (m *Message) Domain() string {
	return m.domain
}

func (m *Message) To() []string {
	return m.to
}

func (m *Message) Tags() []string {
	return m.tags
}

func (m *Message) Campaigns() []string {
	return m.campaigns
}

func (m *Message) DKIM() *bool {
	return m.dkim
}

func (m *Message) DeliveryTime() time.Time {
	return m.deliveryTime
}

func (m *Message) STOPeriod() string {
	return m.stoPeriod
}

func (m *Message) Attachments() []string {
	return m.attachments
}

func (m *Message) ReaderAttachments() []ReaderAttachment {
	return m.readerAttachments
}

func (m *Message) Inlines() []string {
	return m.inlines
}

func (m *Message) ReaderInlines() []ReaderAttachment {
	return m.readerInlines
}

func (m *Message) BufferAttachments() []BufferAttachment {
	return m.bufferAttachments
}

func (m *Message) NativeSend() bool {
	return m.nativeSend
}

func (m *Message) TestMode() bool {
	return m.testMode
}

func (m *Message) Tracking() *bool {
	return m.tracking
}

func (m *Message) TrackingClicks() *string {
	return m.trackingClicks
}

func (m *Message) TrackingOpens() *bool {
	return m.trackingOpens
}

func (m *Message) Variables() map[string]string {
	return m.variables
}

func (m *Message) TemplateVariables() map[string]interface{} {
	return m.templateVariables
}

func (m *Message) RecipientVariables() map[string]map[string]interface{} {
	return m.recipientVariables
}

func (m *Message) TemplateVersionTag() string {
	return m.templateVersionTag
}

func (m *Message) TemplateRenderText() bool {
	return m.templateRenderText
}

func (m *Message) RequireTLS() bool {
	return m.requireTLS
}

func (m *Message) SkipVerification() bool {
	return m.skipVerification
}

// AddReaderAttachment arranges to send a file along with the e-mail message.
// File contents are read from a io.ReadCloser.
// The filename parameter is the resulting filename of the attachment.
// The readCloser parameter is the io.ReadCloser which reads the actual bytes to be used
// as the contents of the attached file.
func (m *Message) AddReaderAttachment(filename string, readCloser io.ReadCloser) {
	ra := ReaderAttachment{Filename: filename, ReadCloser: readCloser}
	m.readerAttachments = append(m.readerAttachments, ra)
}

// AddBufferAttachment arranges to send a file along with the e-mail message.
// File contents are read from the []byte array provided
// The filename parameter is the resulting filename of the attachment.
// The buffer parameter is the []byte array which contains the actual bytes to be used
// as the contents of the attached file.
func (m *Message) AddBufferAttachment(filename string, buffer []byte) {
	ba := BufferAttachment{Filename: filename, Buffer: buffer}
	m.bufferAttachments = append(m.bufferAttachments, ba)
}

// AddAttachment arranges to send a file from the filesystem along with the e-mail message.
// The attachment parameter is a filename, which must refer to a file which actually resides
// in the local filesystem.
func (m *Message) AddAttachment(attachment string) {
	m.attachments = append(m.attachments, attachment)
}

// AddReaderInline arranges to send a file along with the e-mail message.
// File contents are read from a io.ReadCloser.
// The filename parameter is the resulting filename of the attachment.
// The readCloser parameter is the io.ReadCloser which reads the actual bytes to be used
// as the contents of the attached file.
func (m *Message) AddReaderInline(filename string, readCloser io.ReadCloser) {
	ra := ReaderAttachment{Filename: filename, ReadCloser: readCloser}
	m.readerInlines = append(m.readerInlines, ra)
}

// AddInline arranges to send a file along with the e-mail message, but does so
// in a way that its data remains "inline" with the rest of the message.  This
// can be used to send image or font data along with an HTML-encoded message body.
// The attachment parameter is a filename, which must refer to a file which actually resides
// in the local filesystem.
func (m *Message) AddInline(inline string) {
	m.inlines = append(m.inlines, inline)
}

// AddRecipient appends a receiver to the To: header of a message.
// It will return an error if the limit of recipients have been exceeded for this message
func (m *Message) AddRecipient(recipient string) error {
	return m.AddRecipientAndVariables(recipient, nil)
}

// AddRecipientAndVariables appends a receiver to the To: header of a message,
// and as well attaches a set of variables relevant for this recipient.
// It will return an error if the limit of recipients have been exceeded for this message
func (m *Message) AddRecipientAndVariables(r string, vars map[string]interface{}) error {
	if m.RecipientCount() >= MaxNumberOfRecipients {
		return fmt.Errorf("recipient limit exceeded (max %d)", MaxNumberOfRecipients)
	}
	m.to = append(m.to, r)
	if vars != nil {
		if m.recipientVariables == nil {
			m.recipientVariables = make(map[string]map[string]interface{})
		}
		m.recipientVariables[r] = vars
	}
	return nil
}

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
func (m *Message) RecipientCount() int {
	return len(m.to) + m.Specific.RecipientCount()
}

func (m *plainMessage) RecipientCount() int {
	return len(m.bcc) + len(m.cc)
}

func (m *mimeMessage) RecipientCount() int {
	// TODO(v5): len(m.to)
	return 10
}

// SetReplyTo sets the receiver who should receive replies
func (m *Message) SetReplyTo(recipient string) {
	m.AddHeader("Reply-To", recipient)
}

func (m *plainMessage) AddCC(r string) {
	m.cc = append(m.cc, r)
}

func (m *mimeMessage) AddCC(_ string) {}

func (m *plainMessage) AddBCC(r string) {
	m.bcc = append(m.bcc, r)
}

func (m *mimeMessage) AddBCC(_ string) {}

// Deprecated: use SetHTML instead.
//
// TODO(v5): remove this method
func (m *Message) SetHtml(html string) {
	m.SetHTML(html)
}

func (m *plainMessage) SetHTML(h string) {
	m.html = h
}

func (m *mimeMessage) SetHTML(_ string) {}

// Deprecated: use SetAmpHTML instead.
// TODO(v5): remove this method
func (m *Message) SetAMPHtml(html string) {
	m.SetAmpHTML(html)
}

func (m *plainMessage) SetAmpHTML(h string) {
	m.ampHtml = h
}

func (m *mimeMessage) SetAmpHTML(_ string) {}

// AddTag attaches tags to the message.  Tags are useful for metrics gathering and event tracking purposes.
// Refer to the Mailgun documentation for further details.
func (m *Message) AddTag(tag ...string) error {
	if len(m.tags) >= MaxNumberOfTags {
		return fmt.Errorf("cannot add any new tags. Message tag limit (%d) reached", MaxNumberOfTags)
	}

	m.tags = append(m.tags, tag...)
	return nil
}

func (m *plainMessage) SetTemplate(t string) {
	m.template = t
}

func (m *mimeMessage) SetTemplate(_ string) {}

// Deprecated: is no longer supported and is deprecated for new software.
// TODO(v5): remove this method.
func (m *Message) AddCampaign(campaign string) {
	m.campaigns = append(m.campaigns, campaign)
}

// SetDKIM arranges to send the o:dkim header with the message, and sets its value accordingly.
// Refer to the Mailgun documentation for more information.
func (m *Message) SetDKIM(dkim bool) {
	m.dkim = &dkim
}

// EnableNativeSend allows the return path to match the address in the Message.Headers.From:
// field when sending from Mailgun rather than the usual bounce+ address in the return path.
func (m *Message) EnableNativeSend() {
	m.nativeSend = true
}

// EnableTestMode allows submittal of a message, such that it will be discarded by Mailgun.
// This facilitates testing client-side software without actually consuming e-mail resources.
func (m *Message) EnableTestMode() {
	m.testMode = true
}

// SetDeliveryTime schedules the message for transmission at the indicated time.
// Pass nil to remove any installed schedule.
// Refer to the Mailgun documentation for more information.
func (m *Message) SetDeliveryTime(dt time.Time) {
	m.deliveryTime = dt
}

// SetSTOPeriod toggles Send Time Optimization (STO) on a per-message basis.
// String should be set to the number of hours in [0-9]+h format,
// with the minimum being 24h and the maximum being 72h.
// Refer to the Mailgun documentation for more information.
func (m *Message) SetSTOPeriod(stoPeriod string) error {
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
func (m *Message) SetTracking(tracking bool) {
	m.tracking = &tracking
}

// SetTrackingClicks information is found in the Mailgun documentation.
func (m *Message) SetTrackingClicks(trackingClicks bool) {
	m.trackingClicks = ptr(yesNo(trackingClicks))
}

// SetTrackingOptions sets the o:tracking, o:tracking-clicks and o:tracking-opens at once.
func (m *Message) SetTrackingOptions(options *TrackingOptions) {
	m.tracking = &options.Tracking

	m.trackingClicks = &options.TrackingClicks

	m.trackingOpens = &options.TrackingOpens
}

// SetRequireTLS information is found in the Mailgun documentation.
func (m *Message) SetRequireTLS(b bool) {
	m.requireTLS = b
}

// SetSkipVerification information is found in the Mailgun documentation.
func (m *Message) SetSkipVerification(b bool) {
	m.skipVerification = b
}

// SetTrackingOpens information is found in the Mailgun documentation.
func (m *Message) SetTrackingOpens(trackingOpens bool) {
	m.trackingOpens = &trackingOpens
}

// SetTemplateVersion information is found in the Mailgun documentation.
func (m *Message) SetTemplateVersion(tag string) {
	m.templateVersionTag = tag
}

// SetTemplateRenderText information is found in the Mailgun documentation.
func (m *Message) SetTemplateRenderText(render bool) {
	m.templateRenderText = render
}

// AddHeader allows you to send custom MIME headers with the message.
func (m *Message) AddHeader(header, value string) {
	if m.headers == nil {
		m.headers = make(map[string]string)
	}
	m.headers[header] = value
}

// AddVariable lets you associate a set of variables with messages you send,
// which Mailgun can use to, in essence, complete form-mail.
// Refer to the Mailgun documentation for more information.
func (m *Message) AddVariable(variable string, value interface{}) error {
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
func (m *Message) AddTemplateVariable(variable string, value interface{}) error {
	if m.templateVariables == nil {
		m.templateVariables = make(map[string]interface{})
	}
	m.templateVariables[variable] = value
	return nil
}

// AddDomain allows you to use a separate domain for the type of messages you are sending.
func (m *Message) AddDomain(domain string) {
	m.domain = domain
}

// Headers retrieves the http headers associated with this message
func (m *Message) Headers() map[string]string {
	return m.headers
}

// Deprecated: use func Headers() instead.
// TODO(v5): remove this method, it violates https://go.dev/doc/effective_go#Getters
func (m *Message) GetHeaders() map[string]string {
	return m.headers
}

// ErrInvalidMessage is returned by `Send()` when the `mailgun.Message` struct is incomplete
var ErrInvalidMessage = errors.New("message not valid")

type SendableMessage interface {
	Domain() string
	To() []string
	Tags() []string
	// Deprecated: is no longer supported and is deprecated for new software.
	// TODO(v5): remove this method
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

	Specific
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
// See the public mailgun documentation for all possible return codes and error messages
// TODO(v5): switch m to SendableMessage interface - https://bryanftan.medium.com/accept-interfaces-return-structs-in-go-d4cab29a301b
func (mg *MailgunImpl) Send(ctx context.Context, m *Message) (mes string, id string, err error) {
	if mg.domain == "" {
		err = errors.New("you must provide a valid domain before calling Send()")
		return
	}

	invalidChars := ":&'@(),!?#;%+=<>"
	if i := strings.ContainsAny(mg.domain, invalidChars); i {
		err = fmt.Errorf("you called Send() with a domain that contains invalid characters")
		return
	}

	if mg.apiKey == "" {
		err = errors.New("you must provide a valid api-key before calling Send()")
		return
	}

	if !isValid(m) {
		err = ErrInvalidMessage
		return
	}

	if m.STOPeriod() != "" && m.RecipientCount() > 1 {
		err = errors.New("STO can only be used on a per-message basis")
		return
	}
	payload := NewFormDataPayload()

	m.AddValues(payload)
	for _, to := range m.To() {
		payload.addValue("to", to)
	}
	for _, tag := range m.Tags() {
		payload.addValue("o:tag", tag)
	}
	for _, campaign := range m.Campaigns() {
		payload.addValue("o:campaign", campaign)
	}
	if m.DKIM() != nil {
		payload.addValue("o:dkim", yesNo(*m.DKIM()))
	}
	if !m.DeliveryTime().IsZero() {
		payload.addValue("o:deliverytime", formatMailgunTime(m.DeliveryTime()))
	}
	if m.STOPeriod() != "" {
		payload.addValue("o:deliverytime-optimize-period", m.STOPeriod())
	}
	if m.NativeSend() {
		payload.addValue("o:native-send", "yes")
	}
	if m.TestMode() {
		payload.addValue("o:testmode", "yes")
	}
	if m.Tracking() != nil {
		payload.addValue("o:tracking", yesNo(*m.Tracking()))
	}
	if m.TrackingClicks() != nil {
		payload.addValue("o:tracking-clicks", *m.TrackingClicks())
	}
	if m.TrackingOpens() != nil {
		payload.addValue("o:tracking-opens", yesNo(*m.TrackingOpens()))
	}
	if m.RequireTLS() {
		payload.addValue("o:require-tls", trueFalse(m.RequireTLS()))
	}
	if m.SkipVerification() {
		payload.addValue("o:skip-verification", trueFalse(m.SkipVerification()))
	}
	if m.Headers() != nil {
		for header, value := range m.Headers() {
			payload.addValue("h:"+header, value)
		}
	}
	if m.Variables() != nil {
		for variable, value := range m.Variables() {
			payload.addValue("v:"+variable, value)
		}
	}
	if m.TemplateVariables() != nil {
		variableString, err := json.Marshal(m.TemplateVariables())
		if err == nil {
			// the map was marshalled as json so add it
			payload.addValue("h:X-Mailgun-Variables", string(variableString))
		}
	}
	if m.RecipientVariables() != nil {
		j, err := json.Marshal(m.RecipientVariables())
		if err != nil {
			return "", "", err
		}
		payload.addValue("recipient-variables", string(j))
	}
	if m.Attachments() != nil {
		for _, attachment := range m.Attachments() {
			payload.addFile("attachment", attachment)
		}
	}
	if m.ReaderAttachments() != nil {
		for _, readerAttachment := range m.ReaderAttachments() {
			payload.addReadCloser("attachment", readerAttachment.Filename, readerAttachment.ReadCloser)
		}
	}
	if m.BufferAttachments() != nil {
		for _, bufferAttachment := range m.BufferAttachments() {
			payload.addBuffer("attachment", bufferAttachment.Filename, bufferAttachment.Buffer)
		}
	}
	if m.Inlines() != nil {
		for _, inline := range m.Inlines() {
			payload.addFile("inline", inline)
		}
	}

	if m.ReaderInlines() != nil {
		for _, readerAttachment := range m.ReaderInlines() {
			payload.addReadCloser("inline", readerAttachment.Filename, readerAttachment.ReadCloser)
		}
	}

	if m.domain == "" {
		m.domain = mg.Domain()
	}

	if m.TemplateVersionTag() != "" {
		payload.addValue("t:version", m.TemplateVersionTag())
	}

	if m.TemplateRenderText() {
		payload.addValue("t:text", yesNo(m.TemplateRenderText()))
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

	r.mu.RLock()
	defer r.mu.RUnlock()
	if r.capturedCurlOutput != "" {
		mg.mu.Lock()
		defer mg.mu.Unlock()
		mg.capturedCurlOutput = r.capturedCurlOutput
	}

	return
}

func (m *plainMessage) AddValues(p *FormDataPayload) {
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

func (m *mimeMessage) AddValues(p *FormDataPayload) {
	p.addReadCloser("message", "message.mime", m.body)
}

func (m *plainMessage) Endpoint() string {
	return messagesEndpoint
}

func (m *mimeMessage) Endpoint() string {
	return mimeMessagesEndpoint
}

// yesNo translates a true/false boolean value into a yes/no setting suitable for the Mailgun API.
func yesNo(b bool) string {
	if b {
		return "yes"
	}
	return "no"
}

func trueFalse(b bool) string {
	return strconv.FormatBool(b)
}

// isValid returns true if, and only if,
// a Message instance is sufficiently initialized to send via the Mailgun interface.
func isValid(m SendableMessage) bool {
	if m == nil {
		return false
	}

	if !m.IsValid() {
		return false
	}

	if m.RecipientCount() == 0 {
		return false
	}

	if !validateStringList(m.Tags(), false) {
		return false
	}

	if !validateStringList(m.Campaigns(), false) || len(m.Campaigns()) > 3 {
		return false
	}

	return true
}

func (m *plainMessage) IsValid() bool {
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

func (m *mimeMessage) IsValid() bool {
	return m.body != nil
}

// validateStringList returns true if, and only if,
// a slice of strings exists AND all of its elements exist,
// OR if the slice doesn't exist AND it's not required to exist.
// The requireOne parameter indicates whether the list is required to exist.
func validateStringList(list []string, requireOne bool) bool {
	hasOne := false

	if list == nil {
		return !requireOne
	} else {
		for _, a := range list {
			if a == "" {
				return false
			} else {
				// TODO(vtopc): hasOne is always true:
				hasOne = hasOne || true
			}
		}
	}

	return hasOne
}
