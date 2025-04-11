package mtypes

type DomainConnectionResponse struct {
	Connection DomainConnection `json:"connection"`
}

type ListDomainsResponse struct {
	// is -1 if Next() or First() have not been called
	TotalCount int      `json:"total_count"`
	Items      []Domain `json:"items"`
}

// Specify the domain connection options
type DomainConnection struct {
	RequireTLS       bool `json:"require_tls"`
	SkipVerification bool `json:"skip_verification"`
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
