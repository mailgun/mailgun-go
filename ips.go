package mailgun

import (
	"context"

	"github.com/mailgun/mailgun-go/v5/mtypes"
)

// ListIPs returns a list of IPs assigned to your account
func (mg *Client) ListIPs(ctx context.Context, dedicated, enabled bool) ([]mtypes.IPAddress, error) {
	r := newHTTPRequest(generateApiUrl(mg, 3, ipsEndpoint))
	r.setClient(mg.HTTPClient())
	if dedicated {
		r.addParameter("dedicated", "true")
	}
	if enabled {
		r.addParameter("enabled", "true")
	}
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp mtypes.IPAddressListResponse
	if err := getResponseFromJSON(ctx, r, &resp); err != nil {
		return nil, err
	}
	var result []mtypes.IPAddress
	for _, ip := range resp.Items {
		result = append(result, mtypes.IPAddress{IP: ip})
	}
	return result, nil
}

// GetIP returns information about the specified IP
func (mg *Client) GetIP(ctx context.Context, ip string) (mtypes.IPAddress, error) {
	r := newHTTPRequest(generateApiUrl(mg, 3, ipsEndpoint) + "/" + ip)
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var resp mtypes.IPAddress
	err := getResponseFromJSON(ctx, r, &resp)
	return resp, err
}

// ListDomainIPs returns a list of IPs currently assigned to the specified domain.
func (mg *Client) ListDomainIPs(ctx context.Context, domain string) ([]mtypes.IPAddress, error) {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/ips")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var resp mtypes.IPAddressListResponse
	if err := getResponseFromJSON(ctx, r, &resp); err != nil {
		return nil, err
	}
	var result []mtypes.IPAddress
	for _, ip := range resp.Items {
		result = append(result, mtypes.IPAddress{IP: ip})
	}
	return result, nil
}

// AddDomainIP Assign a dedicated IP to the domain specified.
func (mg *Client) AddDomainIP(ctx context.Context, domain, ip string) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/ips")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	payload := newUrlEncodedPayload()
	payload.addValue("ip", ip)
	_, err := makePostRequest(ctx, r, payload)
	return err
}

// DeleteDomainIP Unassign an IP from the domain specified.
func (mg *Client) DeleteDomainIP(ctx context.Context, domain, ip string) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, domainsEndpoint) + "/" + domain + "/ips/" + ip)
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}
