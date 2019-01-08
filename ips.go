package mailgun

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

// Returns a list of IPs assigned to your account
func (mg *MailgunImpl) ListIPS(dedicated bool) ([]IPAddress, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, ipsEndpoint))
	r.setClient(mg.Client())
	if dedicated {
		r.addParameter("dedicated", "true")
	}
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp ipAddressListResponse
	if err := getResponseFromJSON(r, &resp); err != nil {
		return nil, err
	}
	var result []IPAddress
	for _, ip := range resp.Items {
		result = append(result, IPAddress{IP: ip})
	}
	return result, nil
}

// Returns information about the specified IP
func (mg *MailgunImpl) GetIP(ip string) (IPAddress, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, ipsEndpoint) + "/" + ip)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp IPAddress
	err := getResponseFromJSON(r, &resp)
	return resp, err
}

// Returns a list of IPs currently assigned to the specified domain.
func (mg *MailgunImpl) ListDomainIPS() ([]IPAddress, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, domainsEndpoint) + "/" + mg.domain + "/ips")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp ipAddressListResponse
	if err := getResponseFromJSON(r, &resp); err != nil {
		return nil, err
	}
	var result []IPAddress
	for _, ip := range resp.Items {
		result = append(result, IPAddress{IP: ip})
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
