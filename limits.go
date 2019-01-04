package mailgun

type TagLimits struct {
	Limit int `json:"limit"`
	Count int `json:"count"`
}

// Returns tracking settings for a domain
func (mg *MailgunImpl) GetTagLimits(domain string) (TagLimits, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + domain + "/limits/tag")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp TagLimits
	err := getResponseFromJSON(r, &resp)
	return resp, err
}
