package mtypes

type IPAddressListResponse struct {
	TotalCount int      `json:"total_count"`
	Items      []string `json:"items"`
}

type IPAddress struct {
	IP        string `json:"ip"`
	RDNS      string `json:"rdns"`
	Dedicated bool   `json:"dedicated"`
}

type ListIPDomainsResponse struct {
	// is -1 if Next() or First() have not been called
	TotalCount int         `json:"total_count"`
	Items      []DomainIPs `json:"items"`
}

type DomainIPs struct {
	Domain string   `json:"domain"`
	IPs    []string `json:"ips"`
}
