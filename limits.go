package mailgun

import "context"

type TagLimits struct {
	Limit int `json:"limit"`
	Count int `json:"count"`
}

// GetTagLimits returns tracking settings for a domain
func (mg *MailgunImpl) GetTagLimits(ctx context.Context, domain string) (TagLimits, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + domain + "/limits/tag")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp TagLimits
	err := getResponseFromJSON(ctx, r, &resp)
	return resp, err
}
