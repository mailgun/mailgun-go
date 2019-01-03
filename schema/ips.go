package schema

type IPAddressList struct {
	TotalCount int      `json:"total_count"`
	Items      []string `json:"items"`
}

type IPAddress struct {
	IP        string `json:"ip"`
	RDNS      string `json:"rdns"`
	Dedicated bool   `json:"dedicated"`
}

type ExportList struct {
	Items []Export `json:"items"`
}

type Export struct {
	ID     string `json:"id"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

type OK struct {
	Message string `json:"message"`
}
