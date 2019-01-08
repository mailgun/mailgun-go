package mailgun

import (
	"context"
	"strconv"
)

type Unsubscription struct {
	CreatedAt string   `json:"created_at"`
	Tags      []string `json:"tags"`
	ID        string   `json:"id"`
	Address   string   `json:"address"`
}

// ListUnsubscribes retrieves a list of unsubscriptions issued by recipients of mail from your domain.
// Zero is a valid list length.
func (mg *MailgunImpl) ListUnsubscribes(ctx context.Context, limit, skip int) (int, []Unsubscription, error) {
	r := newHTTPRequest(generateApiUrl(mg, unsubscribesEndpoint))
	if limit != DefaultLimit {
		r.addParameter("limit", strconv.Itoa(limit))
	}
	if skip != DefaultSkip {
		r.addParameter("skip", strconv.Itoa(skip))
	}
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var envelope struct {
		TotalCount int              `json:"total_count"`
		Items      []Unsubscription `json:"items"`
	}
	err := getResponseFromJSON(ctx, r, &envelope)
	return envelope.TotalCount, envelope.Items, err
}

// GetUnsubscribes retrieves a list of unsubscriptions by recipient address.
// Zero is a valid list length.
func (mg *MailgunImpl) GetUnsubscribes(ctx context.Context, address string) (int, []Unsubscription, error) {
	r := newHTTPRequest(generateApiUrlWithTarget(mg, unsubscribesEndpoint, address))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var envelope struct {
		TotalCount int              `json:"total_count"`
		Items      []Unsubscription `json:"items"`
	}
	err := getResponseFromJSON(ctx, r, &envelope)
	return envelope.TotalCount, envelope.Items, err
}

// Unsubscribe adds an e-mail address to the domain's unsubscription table.
func (mg *MailgunImpl) Unsubscribe(ctx context.Context, a, t string) error {
	r := newHTTPRequest(generateApiUrl(mg, unsubscribesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("address", a)
	p.addValue("tag", t)
	_, err := makePostRequest(ctx, r, p)
	return err
}

// RemoveUnsubscribe removes the e-mail address given from the domain's unsubscription table.
// If passing in an ID (discoverable from, e.g., ListUnsubscribes()), the e-mail address associated
// with the given ID will be removed.
func (mg *MailgunImpl) RemoveUnsubscribe(ctx context.Context, a string) error {
	r := newHTTPRequest(generateApiUrlWithTarget(mg, unsubscribesEndpoint, a))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// RemoveUnsubscribe removes the e-mail address given from the domain's unsubscription table with a matching tag.
// If passing in an ID (discoverable from, e.g., ListUnsubscribes()), the e-mail address associated
// with the given ID will be removed.
func (mg *MailgunImpl) RemoveUnsubscribeWithTag(ctx context.Context, a, t string) error {
	r := newHTTPRequest(generateApiUrlWithTarget(mg, unsubscribesEndpoint, a))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.addParameter("tag", t)
	_, err := makeDeleteRequest(ctx, r)
	return err
}
