package mtypes

type IPAddressListResponse struct {
	TotalCount        int                           `json:"total_count"`
	Items             []string                      `json:"items"`
	Details           []IPAddressListResponseDetail `json:"details"`
	AssignableToPools []string                      `json:"assignable_to_pools,omitempty"`
}

type IPAddressListResponseDetail struct {
	IP         string `json:"ip"`
	IsOnWarmup bool   `json:"is_on_warmup,omitempty"`
}

type IPAddress struct {
	IP                string `json:"ip"`                            // Always present
	RDNS              string `json:"rdns,omitempty"`                // Filled by GetIP
	Dedicated         bool   `json:"dedicated,omitempty"`           // Filled by GetIP
	AssignableToPools bool   `json:"assignable_to_pools,omitempty"` // Filled by ListIPS
	IsOnWarmup        bool   `json:"is_on_warmup,omitempty"`        // Filled by ListIPS
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
