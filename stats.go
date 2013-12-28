package mailgun

type Stat struct {
	Event      string `json:"event"`
	TotalCount int    `json:"total_count"`
	CreatedAt  string `json:"created_at"`
	Id         string `json:"id"`
	Tags       string `json:"tags"`
}

type statsEnvelope struct {
	TotalCount int    `json:"total_count"`
	Items      []Stat `json:"items"`
}
