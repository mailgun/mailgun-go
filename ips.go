package mailgun

import "context"

type ipAddressListResponse struct {
	TotalCount int      `json:"total_count"`
	Items      []string `json:"items"`
}

type IPAddress struct {
	IP        string `json:"ip"`
	RDNS      string `json:"rdns"`
	Dedicated bool   `json:"dedicated"`
}

type okResp struct {
	ID      string `json:"id,omitempty"`
	Message string `json:"message"`
}

// ListIPS returns a list of IPs assigned to your account
func (mg *MailgunImpl) ListIPS(ctx context.Context, dedicated bool) ([]IPAddress, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, ipsEndpoint))
	r.setClient(mg.Client())
	if dedicated {
		r.addParameter("dedicated", "true")
	}
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp ipAddressListResponse
	if err := getResponseFromJSON(ctx, r, &resp); err != nil {
		return nil, err
	}
	var result []IPAddress
	for _, ip := range resp.Items {
		result = append(result, IPAddress{IP: ip})
	}
	return result, nil
}

// GetIP returns information about the specified IP
func (mg *MailgunImpl) GetIP(ctx context.Context, ip string) (IPAddress, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, ipsEndpoint) + "/" + ip)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp IPAddress
	err := getResponseFromJSON(ctx, r, &resp)
	return resp, err
}

// ListDomainIPS returns a list of IPs currently assigned to the specified domain.
func (mg *MailgunImpl) ListDomainIPS(ctx context.Context) ([]IPAddress, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + mg.domain + "/ips")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp ipAddressListResponse
	if err := getResponseFromJSON(ctx, r, &resp); err != nil {
		return nil, err
	}
	var result []IPAddress
	for _, ip := range resp.Items {
		result = append(result, IPAddress{IP: ip})
	}
	return result, nil
}

// Assign a dedicated IP to the domain specified.
func (mg *MailgunImpl) AddDomainIP(ctx context.Context, ip string) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + mg.domain + "/ips")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("ip", ip)
	_, err := makePostRequest(ctx, r, payload)
	return err
}

// Unassign an IP from the domain specified.
func (mg *MailgunImpl) DeleteDomainIP(ctx context.Context, ip string) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + mg.domain + "/ips/" + ip)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}
