package mailgun

import (
	"context"
)

// GetDomainConnection returns delivery connection settings for the defined domain
func (mg *MailgunImpl) GetDomainConnection(ctx context.Context, domain string) (DomainConnection, error) {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/connection")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp domainConnectionResponse
	err := getResponseFromJSON(ctx, r, &resp)
	return resp.Connection, err
}

// UpdateDomainConnection updates the specified delivery connection settings for the defined domain
func (mg *MailgunImpl) UpdateDomainConnection(ctx context.Context, domain string, settings DomainConnection) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/connection")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("require_tls", boolToString(settings.RequireTLS))
	payload.addValue("skip_verification", boolToString(settings.SkipVerification))
	_, err := makePutRequest(ctx, r, payload)
	return err
}
