package mailgun

import (
	"github.com/mbanzon/simplehttp"
)

const (
	complaintsEndpoint = "complaints"
)

type Complaint struct {
	Count     int    `json:"count"`
	CreatedAt string `json:"created_at"`
	Address   string `json:"address"`
}

type complaintContainer struct {
	TotalCount int         `json:"total_count"`
	Items      []Complaint `json:"items"`
}

func (m *mailgunImpl) GetComplaints(limit, skip int) (int, []interface{}, error) {
	simplehttp.NewGetRequest(generateApiUrl(m, complaintsEndpoint))
	// TODO - this is NOT complete!
	return -1, nil, nil
}
