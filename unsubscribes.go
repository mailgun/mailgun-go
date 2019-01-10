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

// Fetches the list of unsubscribes.
func (mg *MailgunImpl) ListUnsubscribes(ctx context.Context, opts *ListOptions) ([]Unsubscription, error) {
	r := newHTTPRequest(generateApiUrl(mg, unsubscribesEndpoint))
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.setClient(mg.Client())

	if opts != nil && opts.Limit != 0 {
		r.addParameter("limit", strconv.Itoa(opts.Limit))
	}

	if opts != nil && opts.Skip != 0 {
		r.addParameter("skip", strconv.Itoa(opts.Skip))
	}

	var envelope struct {
		TotalCount int              `json:"total_count"`
		Items      []Unsubscription `json:"items"`
	}
	err := getResponseFromJSON(ctx, r, &envelope)
	return envelope.Items, err
}

// Retreives a single unsubscribe record. Can be used to check if a given address is present in the list of unsubscribed users.
func (mg *MailgunImpl) GetUnsubscribe(ctx context.Context, address string) (Unsubscription, error) {
	// TODO: Test this method!
	r := newHTTPRequest(generateApiUrlWithTarget(mg, unsubscribesEndpoint, address))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var envelope struct {
		Unsubscribe Unsubscription `json:"unsubscribe"`
	}
	err := getResponseFromJSON(ctx, r, &envelope)
	return envelope.Unsubscribe, err
}

// Unsubscribe adds an e-mail address to the domain's unsubscription table.
func (mg *MailgunImpl) CreateUnsubscribe(ctx context.Context, address, tag string) error {
	r := newHTTPRequest(generateApiUrl(mg, unsubscribesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("address", address)
	p.addValue("tag", tag)
	_, err := makePostRequest(ctx, r, p)
	return err
}

// Removes the e-mail address given from the domain's unsubscription table.
// If passing in an ID (discoverable from, e.g., ListUnsubscribes()), the e-mail address associated
// with the given ID will be removed.
func (mg *MailgunImpl) DeleteUnsubscribe(ctx context.Context, address string) error {
	r := newHTTPRequest(generateApiUrlWithTarget(mg, unsubscribesEndpoint, address))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// Removes the e-mail address given from the domain's unsubscription table with a matching tag.
// If passing in an ID (discoverable from, e.g., ListUnsubscribes()), the e-mail address associated
// with the given ID will be removed.
func (mg *MailgunImpl) DeleteUnsubscribeWithTag(ctx context.Context, a, t string) error {
	r := newHTTPRequest(generateApiUrlWithTarget(mg, unsubscribesEndpoint, a))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.addParameter("tag", t)
	_, err := makeDeleteRequest(ctx, r)
	return err
}
