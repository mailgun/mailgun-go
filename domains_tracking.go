package mailgun

import (
	"context"
)

// GetDomainTracking returns tracking settings for a domain
func (mg *MailgunImpl) GetDomainTracking(ctx context.Context, domain string) (DomainTracking, error) {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/tracking")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp domainTrackingResponse
	err := getResponseFromJSON(ctx, r, &resp)
	return resp.Tracking, err
}

func (mg *MailgunImpl) UpdateClickTracking(ctx context.Context, domain, active string) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/tracking/click")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("active", active)
	_, err := makePutRequest(ctx, r, payload)
	return err
}

func (mg *MailgunImpl) UpdateUnsubscribeTracking(ctx context.Context, domain, active, htmlFooter, textFooter string) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/tracking/unsubscribe")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("active", active)
	payload.addValue("html_footer", htmlFooter)
	payload.addValue("text_footer", textFooter)
	_, err := makePutRequest(ctx, r, payload)
	return err
}

func (mg *MailgunImpl) UpdateOpenTracking(ctx context.Context, domain, active string) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/tracking/open")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("active", active)
	_, err := makePutRequest(ctx, r, payload)
	return err
}
