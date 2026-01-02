package mtypes

type ListAllDomainsKeysResponse struct {
	TotalCount int         `json:"total_count"`
	Items      []DomainKey `json:"items"`
}

type ListDomainKeysResponse struct {
	Items []DomainKey `json:"items"`
}

type DomainKey struct {
	SigningDomain string    `json:"signing_domain"`
	Selector      string    `json:"selector"`
	DNSRecord     DNSRecord `json:"dns_record"`
}
