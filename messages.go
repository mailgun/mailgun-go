package mailgun

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

// MaxNumberOfRecipients represents the largest batch of recipients that Mailgun can support in a single API call.
// This figure includes To:, Cc:, Bcc:, etc. recipients.
const MaxNumberOfRecipients = 1000

// MaxNumberOfTags represents the maximum number of tags that can be added for a message
const MaxNumberOfTags = 3

// Message structures contain both the message text and the envelop for an e-mail message.
type Message struct {
	to                []string
	tags              []string
	campaigns         []string
	dkim              bool
	deliveryTime      time.Time
	attachments       []string
	readerAttachments []ReaderAttachment
	inlines           []string
	readerInlines     []ReaderAttachment
	bufferAttachments []BufferAttachment

	nativeSend         bool
	testMode           bool
	tracking           bool
	trackingClicks     bool
	trackingOpens      bool
	headers            map[string]string
	variables          map[string]string
	templateVariables  map[string]interface{}
	recipientVariables map[string]map[string]interface{}
	domain             string

	dkimSet           bool
	trackingSet       bool
	trackingClicksSet bool
	trackingOpensSet  bool
	requireTLS        bool
	skipVerification  bool

	specific features
	mg       Mailgun
}

type ReaderAttachment struct {
	Filename   string
	ReadCloser io.ReadCloser
}

type BufferAttachment struct {
	Filename string
	Buffer   []byte
}

// StoredMessage structures contain the (parsed) message content for an email
// sent to a Mailgun account.
//
// The MessageHeaders field is special, in that it's formatted as a slice of pairs.
// Each pair consists of a name [0] and value [1].  Array notation is used instead of a map
// because that's how it's sent over the wire, and it's how encoding/json expects this field
// to be.
type StoredMessage struct {
	Recipients        string             `json:"recipients"`
	Sender            string             `json:"sender"`
	From              string             `json:"from"`
	Subject           string             `json:"subject"`
	BodyPlain         string             `json:"body-plain"`
	StrippedText      string             `json:"stripped-text"`
	StrippedSignature string             `json:"stripped-signature"`
	BodyHtml          string             `json:"body-html"`
	StrippedHtml      string             `json:"stripped-html"`
	Attachments       []StoredAttachment `json:"attachments"`
	MessageUrl        string             `json:"message-url"`
	ContentIDMap      map[string]struct {
		Url         string `json:"url"`
		ContentType string `json:"content-type"`
		Name        string `json:"name"`
		Size        int64  `json:"size"`
	} `json:"content-id-map"`
	MessageHeaders [][]string `json:"message-headers"`
}

// StoredAttachment structures contain information on an attachment associated with a stored message.
type StoredAttachment struct {
	Size        int    `json:"size"`
	Url         string `json:"url"`
	Name        string `json:"name"`
	ContentType string `json:"content-type"`
}

type StoredMessageRaw struct {
	Recipients string `json:"recipients"`
	Sender     string `json:"sender"`
	From       string `json:"from"`
	Subject    string `json:"subject"`
	BodyMime   string `json:"body-mime"`
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

// features abstracts the common characteristics between regular and MIME messages.
// addCC, addBCC, recipientCount, setHtml and setAMPHtml are invoked via the package-global AddCC, AddBCC,
// RecipientCount, SetHtml and SetAMPHtml calls, as these functions are ignored for MIME messages.
// Send() invokes addValues to add message-type-specific MIME headers for the API call
// to Mailgun.  isValid yeilds true if and only if the message is valid enough for sending
// through the API.  Finally, endpoint() tells Send() which endpoint to use to submit the API call.
type features interface {
	addCC(string)
	addBCC(string)
	setHtml(string)
	setAMPHtml(string)
	addValues(*formDataPayload)
	isValid() bool
	endpoint() string
	recipientCount() int
	setTemplate(string)
}

// NewMessage returns a new e-mail message with the simplest envelop needed to send.
//
// Unlike the global function,
// this method supports arbitrary-sized recipient lists by
// automatically sending mail in batches of up to MaxNumberOfRecipients.
//
// To support batch sending, you don't want to provide a fixed To: header at this point.
// Pass nil as the to parameter to skip adding the To: header at this stage.
// You can do this explicitly, or implicitly, as follows:
//
//     // Note absence of To parameter(s)!
//     m := mg.NewMessage("me@example.com", "Help save our planet", "Hello world!")
//
// Note that you'll need to invoke the AddRecipientAndVariables or AddRecipient method
// before sending, though.
func (mg *MailgunImpl) NewMessage(from, subject, text string, to ...string) *Message {
	return &Message{
		specific: &plainMessage{
			from:    from,
			subject: subject,
			text:    text,
		},
		to: to,
		mg: mg,
	}
}

// NewMIMEMessage creates a new MIME message.  These messages are largely canned;
// you do not need to invoke setters to set message-related headers.
// However, you do still need to call setters for Mailgun-specific settings.
//
// Unlike the global function,
// this method supports arbitrary-sized recipient lists by
// automatically sending mail in batches of up to MaxNumberOfRecipients.
//
// To support batch sending, you don't want to provide a fixed To: header at this point.
// Pass nil as the to parameter to skip adding the To: header at this stage.
// You can do this explicitly, or implicitly, as follows:
//
//     // Note absence of To parameter(s)!
//     m := mg.NewMessage("me@example.com", "Help save our planet", "Hello world!")
//
// Note that you'll need to invoke the AddRecipientAndVariables or AddRecipient method
// before sending, though.
func (mg *MailgunImpl) NewMIMEMessage(body io.ReadCloser, to ...string) *Message {
	return &Message{
		specific: &mimeMessage{
			body: body,
		},
		to: to,
		mg: mg,
	}
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
	return len(m.to) + m.specific.recipientCount()
}

func (pm *plainMessage) recipientCount() int {
	return len(pm.bcc) + len(pm.cc)
}

func (mm *mimeMessage) recipientCount() int {
	return 10
}

func (m *Message) send(ctx context.Context) (string, string, error) {
	return m.mg.Send(ctx, m)
}

// SetReplyTo sets the receiver who should receive replies
func (m *Message) SetReplyTo(recipient string) {
	m.AddHeader("Reply-To", recipient)
}

// AddCC appends a receiver to the carbon-copy header of a message.
func (m *Message) AddCC(recipient string) {
	m.specific.addCC(recipient)
}

func (pm *plainMessage) addCC(r string) {
	pm.cc = append(pm.cc, r)
}

func (mm *mimeMessage) addCC(_ string) {}

// AddBCC appends a receiver to the blind-carbon-copy header of a message.
func (m *Message) AddBCC(recipient string) {
	m.specific.addBCC(recipient)
}

func (pm *plainMessage) addBCC(r string) {
	pm.bcc = append(pm.bcc, r)
}

func (mm *mimeMessage) addBCC(_ string) {}

// SetHtml is a helper. If you're sending a message that isn't already MIME encoded, SetHtml() will arrange to bundle
// an HTML representation of your message in addition to your plain-text body.
func (m *Message) SetHtml(html string) {
	m.specific.setHtml(html)
}

func (pm *plainMessage) setHtml(h string) {
	pm.html = h
}

func (mm *mimeMessage) setHtml(_ string) {}

// SetAMP is a helper. If you're sending a message that isn't already MIME encoded, SetAMP() will arrange to bundle
// an AMP-For-Email representation of your message in addition to your html & plain-text content.
func (m *Message) SetAMPHtml(html string) {
	m.specific.setAMPHtml(html)
}

func (pm *plainMessage) setAMPHtml(h string) {
	pm.ampHtml = h
}

func (mm *mimeMessage) setAMPHtml(_ string) {}

// AddTag attaches tags to the message.  Tags are useful for metrics gathering and event tracking purposes.
// Refer to the Mailgun documentation for further details.
func (m *Message) AddTag(tag ...string) error {
	if len(m.tags) >= MaxNumberOfTags {
		return fmt.Errorf("cannot add any new tags. Message tag limit (%d) reached", MaxNumberOfTags)
	}

	m.tags = append(m.tags, tag...)
	return nil
}

// SetTemplate sets the name of a template stored via the template API.
// See https://documentation.mailgun.com/en/latest/user_manual.html#templating
func (m *Message) SetTemplate(t string) {
	m.specific.setTemplate(t)
}

func (pm *plainMessage) setTemplate(t string) {
	pm.template = t
}

func (mm *mimeMessage) setTemplate(t string) {}

// AddCampaign is no longer supported and is deprecated for new software.
func (m *Message) AddCampaign(campaign string) {
	m.campaigns = append(m.campaigns, campaign)
}

// SetDKIM arranges to send the o:dkim header with the message, and sets its value accordingly.
// Refer to the Mailgun documentation for more information.
func (m *Message) SetDKIM(dkim bool) {
	m.dkim = dkim
	m.dkimSet = true
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

// SetTracking sets the o:tracking message parameter to adjust, on a message-by-message basis,
// whether or not Mailgun will rewrite URLs to facilitate event tracking.
// Events tracked includes opens, clicks, unsubscribes, etc.
// Note: simply calling this method ensures that the o:tracking header is passed in with the message.
// Its yes/no setting is determined by the call's parameter.
// Note that this header is not passed on to the final recipient(s).
// Refer to the Mailgun documentation for more information.
func (m *Message) SetTracking(tracking bool) {
	m.tracking = tracking
	m.trackingSet = true
}

// SetTrackingClicks information is found in the Mailgun documentation.
func (m *Message) SetTrackingClicks(trackingClicks bool) {
	m.trackingClicks = trackingClicks
	m.trackingClicksSet = true
}

// SetRequireTLS information is found in the Mailgun documentation.
func (m *Message) SetRequireTLS(b bool) {
	m.requireTLS = b
}

// SetSkipVerification information is found in the Mailgun documentation.
func (m *Message) SetSkipVerification(b bool) {
	m.skipVerification = b
}

//SetTrackingOpens information is found in the Mailgun documentation.
func (m *Message) SetTrackingOpens(trackingOpens bool) {
	m.trackingOpens = trackingOpens
	m.trackingOpensSet = true
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

// GetHeaders retrieves the http headers associated with this message
func (m *Message) GetHeaders() map[string]string {
	return m.headers
}

// ErrInvalidMessage is returned by `Send()` when the `mailgun.Message` struct is incomplete
var ErrInvalidMessage = errors.New("message not valid")

// Send attempts to queue a message (see Message, NewMessage, and its methods) for delivery.
// It returns the Mailgun server response, which consists of two components:
// a human-readable status message, and a message ID.  The status and message ID are set only
// if no error occurred.
func (mg *MailgunImpl) Send(ctx context.Context, message *Message) (mes string, id string, err error) {
	if mg.domain == "" {
		err = errors.New("you must provide a valid domain before calling Send()")
		return
	}

	if mg.apiKey == "" {
		err = errors.New("you must provide a valid api-key before calling Send()")
		return
	}

	if !isValid(message) {
		err = ErrInvalidMessage
		return
	}
	payload := newFormDataPayload()

	message.specific.addValues(payload)
	for _, to := range message.to {
		payload.addValue("to", to)
	}
	for _, tag := range message.tags {
		payload.addValue("o:tag", tag)
	}
	for _, campaign := range message.campaigns {
		payload.addValue("o:campaign", campaign)
	}
	if message.dkimSet {
		payload.addValue("o:dkim", yesNo(message.dkim))
	}
	if !message.deliveryTime.IsZero() {
		payload.addValue("o:deliverytime", formatMailgunTime(message.deliveryTime))
	}
	if message.nativeSend {
		payload.addValue("o:native-send", "yes")
	}
	if message.testMode {
		payload.addValue("o:testmode", "yes")
	}
	if message.trackingSet {
		payload.addValue("o:tracking", yesNo(message.tracking))
	}
	if message.trackingClicksSet {
		payload.addValue("o:tracking-clicks", yesNo(message.trackingClicks))
	}
	if message.trackingOpensSet {
		payload.addValue("o:tracking-opens", yesNo(message.trackingOpens))
	}
	if message.requireTLS {
		payload.addValue("o:require-tls", trueFalse(message.requireTLS))
	}
	if message.skipVerification {
		payload.addValue("o:skip-verification", trueFalse(message.skipVerification))
	}
	if message.headers != nil {
		for header, value := range message.headers {
			payload.addValue("h:"+header, value)
		}
	}
	if message.variables != nil {
		for variable, value := range message.variables {
			payload.addValue("v:"+variable, value)
		}
	}
	if message.templateVariables != nil {
		variableString, err := json.Marshal(message.templateVariables)
		if err == nil {
			// the map was marshalled as json so add it
			payload.addValue("h:X-Mailgun-Variables", string(variableString))
		}
	}
	if message.recipientVariables != nil {
		j, err := json.Marshal(message.recipientVariables)
		if err != nil {
			return "", "", err
		}
		payload.addValue("recipient-variables", string(j))
	}
	if message.attachments != nil {
		for _, attachment := range message.attachments {
			payload.addFile("attachment", attachment)
		}
	}
	if message.readerAttachments != nil {
		for _, readerAttachment := range message.readerAttachments {
			payload.addReadCloser("attachment", readerAttachment.Filename, readerAttachment.ReadCloser)
		}
	}
	if message.bufferAttachments != nil {
		for _, bufferAttachment := range message.bufferAttachments {
			payload.addBuffer("attachment", bufferAttachment.Filename, bufferAttachment.Buffer)
		}
	}
	if message.inlines != nil {
		for _, inline := range message.inlines {
			payload.addFile("inline", inline)
		}
	}

	if message.readerInlines != nil {
		for _, readerAttachment := range message.readerInlines {
			payload.addReadCloser("inline", readerAttachment.Filename, readerAttachment.ReadCloser)
		}
	}

	if message.domain == "" {
		message.domain = mg.Domain()
	}

	r := newHTTPRequest(generateApiUrlWithDomain(mg, message.specific.endpoint(), message.domain))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var response sendMessageResponse
	err = postResponseFromJSON(ctx, r, payload, &response)
	if err == nil {
		mes = response.Message
		id = response.Id
	}

	return
}

func (pm *plainMessage) addValues(p *formDataPayload) {
	p.addValue("from", pm.from)
	p.addValue("subject", pm.subject)
	p.addValue("text", pm.text)
	for _, cc := range pm.cc {
		p.addValue("cc", cc)
	}
	for _, bcc := range pm.bcc {
		p.addValue("bcc", bcc)
	}
	if pm.html != "" {
		p.addValue("html", pm.html)
	}
	if pm.template != "" {
		p.addValue("template", pm.template)
	}
	if pm.ampHtml != "" {
		p.addValue("amp-html", pm.ampHtml)
	}
}

func (mm *mimeMessage) addValues(p *formDataPayload) {
	p.addReadCloser("message", "message.mime", mm.body)
}

func (pm *plainMessage) endpoint() string {
	return messagesEndpoint
}

func (mm *mimeMessage) endpoint() string {
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
	if b {
		return "true"
	}
	return "false"
}

// isValid returns true if, and only if,
// a Message instance is sufficiently initialized to send via the Mailgun interface.
func isValid(m *Message) bool {
	if m == nil {
		return false
	}

	if !m.specific.isValid() {
		return false
	}

	if m.RecipientCount() == 0 {
		return false
	}

	if !validateStringList(m.tags, false) {
		return false
	}

	if !validateStringList(m.campaigns, false) || len(m.campaigns) > 3 {
		return false
	}

	return true
}

func (pm *plainMessage) isValid() bool {
	if pm.from == "" {
		return false
	}

	if !validateStringList(pm.cc, false) {
		return false
	}

	if !validateStringList(pm.bcc, false) {
		return false
	}

	if pm.template != "" {
		// pm.text or pm.html not needed if template is supplied
		return true
	}

	if pm.text == "" && pm.html == "" {
		return false
	}

	return true
}

func (mm *mimeMessage) isValid() bool {
	return mm.body != nil
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
				hasOne = hasOne || true
			}
		}
	}

	return hasOne
}

// GetStoredMessage retrieves information about a received e-mail message.
// This provides visibility into, e.g., replies to a message sent to a mailing list.
func (mg *MailgunImpl) GetStoredMessage(ctx context.Context, url string) (StoredMessage, error) {
	r := newHTTPRequest(url)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var response StoredMessage
	err := getResponseFromJSON(ctx, r, &response)
	return response, err
}

// Given a storage id resend the stored message to the specified recipients
func (mg *MailgunImpl) ReSend(ctx context.Context, url string, recipients ...string) (string, string, error) {
	r := newHTTPRequest(url)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newFormDataPayload()

	if len(recipients) == 0 {
		return "", "", errors.New("must provide at least one recipient")
	}

	for _, to := range recipients {
		payload.addValue("to", to)
	}

	var resp sendMessageResponse
	err := postResponseFromJSON(ctx, r, payload, &resp)
	if err != nil {
		return "", "", err
	}
	return resp.Message, resp.Id, nil

}

// GetStoredMessageRaw retrieves the raw MIME body of a received e-mail message.
// Compared to GetStoredMessage, it gives access to the unparsed MIME body, and
// thus delegates to the caller the required parsing.
func (mg *MailgunImpl) GetStoredMessageRaw(ctx context.Context, url string) (StoredMessageRaw, error) {
	r := newHTTPRequest(url)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.addHeader("Accept", "message/rfc2822")

	var response StoredMessageRaw
	err := getResponseFromJSON(ctx, r, &response)
	return response, err
}

// Deprecated: Use GetStoreMessage() instead
func (mg *MailgunImpl) GetStoredMessageForURL(ctx context.Context, url string) (StoredMessage, error) {
	return mg.GetStoredMessage(ctx, url)
}

// Deprecated: Use GetStoreMessageRaw() instead
func (mg *MailgunImpl) GetStoredMessageRawForURL(ctx context.Context, url string) (StoredMessageRaw, error) {
	return mg.GetStoredMessageRaw(ctx, url)
}

// GetStoredAttachment retrieves the raw MIME body of a received e-mail message attachment.
func (mg *MailgunImpl) GetStoredAttachment(ctx context.Context, url string) ([]byte, error) {
	r := newHTTPRequest(url)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.addHeader("Accept", "message/rfc2822")

	response, err := makeGetRequest(ctx, r)

	return response.Data, err
}
