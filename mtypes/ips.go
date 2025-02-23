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
