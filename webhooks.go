package mailgun

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

	"github.com/mailgun/mailgun-go/v3/events"
)

type UrlOrUrls struct {
	Urls []string `json:"urls"`
	Url  string   `json:"url"`
}

type WebHooksListResponse struct {
	Webhooks map[string]UrlOrUrls `json:"webhooks"`
}

type WebHookResponse struct {
	Webhook UrlOrUrls `json:"webhook"`
}

// ListWebhooks returns the complete set of webhooks configured for your domain.
// Note that a zero-length mapping is not an error.
func (mg *MailgunImpl) ListWebhooks(ctx context.Context) (map[string][]string, error) {
	r := newHTTPRequest(generateDomainApiUrl(mg, webhooksEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var body WebHooksListResponse
	err := getResponseFromJSON(ctx, r, &body)
	if err != nil {
		return nil, err
	}

	hooks := make(map[string][]string, 0)
	for k, v := range body.Webhooks {
		if v.Url != "" {
			hooks[k] = []string{v.Url}
		}
		if len(v.Urls) != 0 {
			hooks[k] = append(hooks[k], v.Urls...)
		}
	}
	return hooks, nil
}

// CreateWebhook installs a new webhook for your domain.
func (mg *MailgunImpl) CreateWebhook(ctx context.Context, kind string, urls []string) error {
	r := newHTTPRequest(generateDomainApiUrl(mg, webhooksEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("id", kind)
	for _, url := range urls {
		p.addValue("url", url)
	}
	_, err := makePostRequest(ctx, r, p)
	return err
}

// DeleteWebhook removes the specified webhook from your domain's configuration.
func (mg *MailgunImpl) DeleteWebhook(ctx context.Context, kind string) error {
	r := newHTTPRequest(generateDomainApiUrl(mg, webhooksEndpoint) + "/" + kind)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// GetWebhook retrieves the currently assigned webhook URL associated with the provided type of webhook.
func (mg *MailgunImpl) GetWebhook(ctx context.Context, kind string) ([]string, error) {
	r := newHTTPRequest(generateDomainApiUrl(mg, webhooksEndpoint) + "/" + kind)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var body WebHookResponse
	if err := getResponseFromJSON(ctx, r, &body); err != nil {
		return nil, err
	}

	if body.Webhook.Url != "" {
		return []string{body.Webhook.Url}, nil
	}
	if len(body.Webhook.Urls) != 0 {
		return body.Webhook.Urls, nil
	}
	return nil, fmt.Errorf("webhook '%s' returned no urls", kind)
}

// UpdateWebhook replaces one webhook setting for another.
func (mg *MailgunImpl) UpdateWebhook(ctx context.Context, kind string, urls []string) error {
	r := newHTTPRequest(generateDomainApiUrl(mg, webhooksEndpoint) + "/" + kind)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	for _, url := range urls {
		p.addValue("url", url)
	}
	_, err := makePutRequest(ctx, r, p)
	return err
}

// Represents the signature portion of the webhook POST body
type Signature struct {
	TimeStamp string `json:"timestamp"`
	Token     string `json:"token"`
	Signature string `json:"signature"`
}

// Represents the JSON payload provided when a Webhook is called by mailgun
type WebhookPayload struct {
	Signature Signature      `json:"signature"`
	EventData events.RawJSON `json:"event-data"`
}

// Use this method to parse the webhook signature given as JSON in the webhook response
func (mg *MailgunImpl) VerifyWebhookSignature(sig Signature) (verified bool, err error) {
	h := hmac.New(sha256.New, []byte(mg.APIKey()))
	io.WriteString(h, sig.TimeStamp)
	io.WriteString(h, sig.Token)

	calculatedSignature := h.Sum(nil)
	signature, err := hex.DecodeString(sig.Signature)
	if err != nil {
		return false, err
	}
	if len(calculatedSignature) != len(signature) {
		return false, nil
	}

	return subtle.ConstantTimeCompare(signature, calculatedSignature) == 1, nil
}

// Deprecated: Please use the VerifyWebhookSignature() to parse the latest
// version of WebHooks from mailgun
func (mg *MailgunImpl) VerifyWebhookRequest(req *http.Request) (verified bool, err error) {
	h := hmac.New(sha256.New, []byte(mg.APIKey()))
	io.WriteString(h, req.FormValue("timestamp"))
	io.WriteString(h, req.FormValue("token"))

	calculatedSignature := h.Sum(nil)
	signature, err := hex.DecodeString(req.FormValue("signature"))
	if err != nil {
		return false, err
	}
	if len(calculatedSignature) != len(signature) {
		return false, nil
	}

	return subtle.ConstantTimeCompare(signature, calculatedSignature) == 1, nil
}
