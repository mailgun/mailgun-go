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
// TODO(v5): rename to CommonMessage and remove Specific
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
	templateVariables  map[string]any
	recipientVariables map[string]map[string]any
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

// PlainMessage contains fields relevant to plain API-synthesized messages.
// You're expected to use various setters to set most of these attributes,
// although from, subject, and text are set when the message is created with
// NewMessage.
// TODO(v5): embed CommonMessage
type PlainMessage struct {
	from     string
	cc       []string
	bcc      []string
	subject  string
	text     string
	html     string
	ampHtml  string
	template string
}

func (m *PlainMessage) From() string {
	return m.from
}

func (m *PlainMessage) CC() []string {
	return m.cc
}

func (m *PlainMessage) BCC() []string {
	return m.bcc
}

func (m *PlainMessage) Subject() string {
	return m.subject
}

func (m *PlainMessage) Text() string {
	return m.text
}

func (m *PlainMessage) HTML() string {
	return m.html
}

func (m *PlainMessage) AmpHTML() string {
	return m.ampHtml
}

func (m *PlainMessage) Template() string {
	return m.template
}

// MimeMessage contains fields relevant to pre-packaged MIME messages.
// TODO(v5): embed CommonMessage
type MimeMessage struct {
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
// TODO(v5): should return PlainMessage
func NewMessage(from, subject, text string, to ...string) *Message {
	return &Message{
		Specific: &PlainMessage{
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
// TODO(v5): should return MimeMessage
func NewMIMEMessage(body io.ReadCloser, to ...string) *Message {
	return &Message{
		Specific: &MimeMessage{
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

func (m *Message) TemplateVariables() map[string]any {
	return m.templateVariables
}

func (m *Message) RecipientVariables() map[string]map[string]any {
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
func (m *Message) AddRecipientAndVariables(r string, vars map[string]any) error {
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
	return len(m.To()) + m.Specific.RecipientCount()
}

func (m *PlainMessage) RecipientCount() int {
	return len(m.BCC()) + len(m.CC())
}

func (*MimeMessage) RecipientCount() int {
	// TODO(v5): 10 + len(m.To())
	return 10
}

// SetReplyTo sets the receiver who should receive replies
func (m *Message) SetReplyTo(recipient string) {
	m.AddHeader("Reply-To", recipient)
}

func (m *PlainMessage) AddCC(r string) {
	m.cc = append(m.CC(), r)
}

func (*MimeMessage) AddCC(_ string) {}

func (m *PlainMessage) AddBCC(r string) {
	m.bcc = append(m.BCC(), r)
}

func (*MimeMessage) AddBCC(_ string) {}

// Deprecated: use SetHTML instead.
//
// TODO(v5): remove this method
func (m *Message) SetHtml(html string) {
	m.SetHTML(html)
}

func (m *PlainMessage) SetHTML(h string) {
	m.html = h
}

func (*MimeMessage) SetHTML(_ string) {}

// Deprecated: use SetAmpHTML instead.
// TODO(v5): remove this method
func (m *Message) SetAMPHtml(html string) {
	m.SetAmpHTML(html)
}

func (m *PlainMessage) SetAmpHTML(h string) {
	m.ampHtml = h
}

func (*MimeMessage) SetAmpHTML(_ string) {}

// AddTag attaches tags to the message.  Tags are useful for metrics gathering and event tracking purposes.
// Refer to the Mailgun documentation for further details.
func (m *Message) AddTag(tag ...string) error {
	if len(m.Tags()) >= MaxNumberOfTags {
		return fmt.Errorf("cannot add any new tags. Message tag limit (%d) reached", MaxNumberOfTags)
	}

	m.tags = append(m.Tags(), tag...)
	return nil
}

func (m *PlainMessage) SetTemplate(t string) {
	m.template = t
}

func (*MimeMessage) SetTemplate(_ string) {}

// Deprecated: is no longer supported and is deprecated for new software.
// TODO(v5): remove this method.
func (m *Message) AddCampaign(campaign string) {
	m.campaigns = append(m.Campaigns(), campaign)
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
func (m *Message) AddVariable(variable string, value any) error {
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
func (m *Message) AddTemplateVariable(variable string, value any) error {
	if m.templateVariables == nil {
		m.templateVariables = make(map[string]any)
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
	TemplateVariables() map[string]any
	RecipientVariables() map[string]map[string]any
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
func (mg *MailgunImpl) Send(ctx context.Context, m *Message) (mes, id string, err error) {
	if mg.domain == "" {
		err = errors.New("you must provide a valid domain before calling Send()")
		return "", "", err
	}

	invalidChars := ":&'@(),!?#;%+=<>"
	if i := strings.ContainsAny(mg.domain, invalidChars); i {
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

	m.Specific.AddValues(payload)

	err = addMessageValues(payload, m)
	if err != nil {
		return "", "", err
	}

	// TODO(v5): remove due for domain agnostic API:
	if m.Domain() == "" {
		m.domain = mg.Domain()
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

	return mes, id, err
}

func addMessageValues(dst *FormDataPayload, src SendableMessage) error {
	addMessageOptions(dst, src)
	addMessageHeaders(dst, src)

	err := addMessageVariables(dst, src)
	if err != nil {
		return err
	}

	addMessageAttachment(dst, src)

	return nil
}

func addMessageOptions(dst *FormDataPayload, src SendableMessage) {
	for _, to := range src.To() {
		dst.addValue("to", to)
	}

	for _, tag := range src.Tags() {
		dst.addValue("o:tag", tag)
	}
	for _, campaign := range src.Campaigns() {
		// TODO(v5): deprecated - https://documentation.mailgun.com/docs/mailgun/api-reference/openapi-final/tag/Messages/
		dst.addValue("o:campaign", campaign)
	}
	if src.DKIM() != nil {
		dst.addValue("o:dkim", yesNo(*src.DKIM()))
	}
	if !src.DeliveryTime().IsZero() {
		dst.addValue("o:deliverytime", formatMailgunTime(src.DeliveryTime()))
	}
	if src.STOPeriod() != "" {
		dst.addValue("o:deliverytime-optimize-period", src.STOPeriod())
	}
	if src.NativeSend() {
		dst.addValue("o:native-send", "yes")
	}
	if src.TestMode() {
		dst.addValue("o:testmode", "yes")
	}
	if src.Tracking() != nil {
		dst.addValue("o:tracking", yesNo(*src.Tracking()))
	}
	if src.TrackingClicks() != nil {
		dst.addValue("o:tracking-clicks", *src.TrackingClicks())
	}
	if src.TrackingOpens() != nil {
		dst.addValue("o:tracking-opens", yesNo(*src.TrackingOpens()))
	}
	if src.RequireTLS() {
		dst.addValue("o:require-tls", trueFalse(src.RequireTLS()))
	}
	if src.SkipVerification() {
		dst.addValue("o:skip-verification", trueFalse(src.SkipVerification()))
	}

	if src.TemplateVersionTag() != "" {
		dst.addValue("t:version", src.TemplateVersionTag())
	}
	if src.TemplateRenderText() {
		dst.addValue("t:text", yesNo(src.TemplateRenderText()))
	}
}

func addMessageHeaders(dst *FormDataPayload, src SendableMessage) {
	if src.Headers() != nil {
		for header, value := range src.Headers() {
			dst.addValue("h:"+header, value)
		}
	}
}

func addMessageVariables(dst *FormDataPayload, src SendableMessage) error {
	if src.Variables() != nil {
		for variable, value := range src.Variables() {
			dst.addValue("v:"+variable, value)
		}
	}
	if src.TemplateVariables() != nil {
		variableString, err := json.Marshal(src.TemplateVariables())
		if err == nil {
			// the map was marshalled as json so add it
			dst.addValue("h:X-Mailgun-Variables", string(variableString))
		}
	}
	if src.RecipientVariables() != nil {
		j, err := json.Marshal(src.RecipientVariables())
		if err != nil {
			return err
		}
		dst.addValue("recipient-variables", string(j))
	}

	return nil
}

func addMessageAttachment(dst *FormDataPayload, src SendableMessage) {
	if src.Attachments() != nil {
		for _, attachment := range src.Attachments() {
			dst.addFile("attachment", attachment)
		}
	}
	if src.ReaderAttachments() != nil {
		for _, readerAttachment := range src.ReaderAttachments() {
			dst.addReadCloser("attachment", readerAttachment.Filename, readerAttachment.ReadCloser)
		}
	}
	if src.BufferAttachments() != nil {
		for _, bufferAttachment := range src.BufferAttachments() {
			dst.addBuffer("attachment", bufferAttachment.Filename, bufferAttachment.Buffer)
		}
	}
	if src.Inlines() != nil {
		for _, inline := range src.Inlines() {
			dst.addFile("inline", inline)
		}
	}
	if src.ReaderInlines() != nil {
		for _, readerAttachment := range src.ReaderInlines() {
			dst.addReadCloser("inline", readerAttachment.Filename, readerAttachment.ReadCloser)
		}
	}
}

func (m *PlainMessage) AddValues(p *FormDataPayload) {
	p.addValue("from", m.From())
	p.addValue("subject", m.Subject())
	p.addValue("text", m.Text())
	for _, cc := range m.CC() {
		p.addValue("cc", cc)
	}
	for _, bcc := range m.BCC() {
		p.addValue("bcc", bcc)
	}
	if m.HTML() != "" {
		p.addValue("html", m.HTML())
	}
	if m.Template() != "" {
		p.addValue("template", m.Template())
	}
	if m.AmpHTML() != "" {
		p.addValue("amp-html", m.AmpHTML())
	}
}

func (m *MimeMessage) AddValues(p *FormDataPayload) {
	p.addReadCloser("message", "message.mime", m.body)
}

func (*PlainMessage) Endpoint() string {
	return messagesEndpoint
}

func (*MimeMessage) Endpoint() string {
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

func (m *PlainMessage) IsValid() bool {
	if !validateStringList(m.CC(), false) {
		return false
	}

	if !validateStringList(m.BCC(), false) {
		return false
	}

	if m.From() == "" {
		return false
	}

	if m.Template() != "" {
		// m.text or m.html not needed if template is supplied
		return true
	}

	if m.Text() == "" && m.HTML() == "" {
		return false
	}

	return true
}

func (m *MimeMessage) IsValid() bool {
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
	}

	for _, a := range list {
		if a == "" {
			return false
		}

		hasOne = true
	}

	return hasOne
}
