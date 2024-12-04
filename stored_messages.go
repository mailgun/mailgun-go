package mailgun

import (
	"context"
	"errors"
)

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
func (mg *MailgunImpl) ReSend(ctx context.Context, url string, recipients ...string) (msg, id string, err error) {
	r := newHTTPRequest(url)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := NewFormDataPayload()

	if len(recipients) == 0 {
		return "", "", errors.New("must provide at least one recipient")
	}

	for _, to := range recipients {
		payload.addValue("to", to)
	}

	var resp sendMessageResponse
	err = postResponseFromJSON(ctx, r, payload, &resp)
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
	if err != nil {
		return nil, err
	}

	return response.Data, err
}
