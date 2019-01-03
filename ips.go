package mailgun

import "github.com/mailgun/mailgun-go/schema"

func (mg *MailgunImpl) ListIPS(dedicated bool) ([]schema.IPAddress, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, ipsEndpoint))
	r.setClient(mg.Client())
	if dedicated {
		r.addParameter("dedicated", "true")
	}
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp schema.IPAddressList
	if err := getResponseFromJSON(r, &resp); err != nil {
		return nil, err
	}
	var result []schema.IPAddress
	for _, ip := range resp.Items {
		result = append(result, schema.IPAddress{IP: ip})
	}
	return result, nil
}

func (mg *MailgunImpl) GetIP(ip string) (schema.IPAddress, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, ipsEndpoint) + "/" + ip)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp schema.IPAddress
	err := getResponseFromJSON(r, &resp)
	return resp, err
}

func (mg *MailgunImpl) ListDomainIPS() ([]schema.IPAddress, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + mg.domain + "/ips")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp schema.IPAddressList
	if err := getResponseFromJSON(r, &resp); err != nil {
		return nil, err
	}
	var result []schema.IPAddress
	for _, ip := range resp.Items {
		result = append(result, schema.IPAddress{IP: ip})
	}
	return result, nil
}

// Assign a dedicated IP to the domain specified.
func (mg *MailgunImpl) AddDomainIP(ip string) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + mg.domain + "/ips")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("ip", ip)
	_, err := makePostRequest(r, payload)
	return err
}

// Unassign an IP from the domain specified.
func (mg *MailgunImpl) DeleteDomainIP(ip string) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + mg.domain + "/ips/" + ip)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(r)
	return err
}
