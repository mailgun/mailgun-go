package mailgun

import (
	"context"
)

// UpdateDomainDkimSelector updates DKIM authority for a domain
func (mg *MailgunImpl) UpdateDomainDkimSelector(ctx context.Context, domain string, self bool) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/dkim_selector")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("dkim_selector", boolToString(self))
	_, err := makePutRequest(ctx, r, payload)
	return err
}
