package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"strconv"
)

const (
	complaintsEndpoint = "complaints"
)

type Complaint struct {
	Count     int    `json:"count"`
	CreatedAt string `json:"created_at"`
	Address   string `json:"address"`
}

type complaintsEnvelope struct {
	TotalCount int         `json:"total_count"`
	Items      []Complaint `json:"items"`
}

func (m *mailgunImpl) GetComplaints(limit, skip int) (int, []Complaint, error) {
	r := simplehttp.NewGetRequest(generateApiUrl(m, complaintsEndpoint))
	if limit != -1 {
		r.AddParameter("limit", strconv.Itoa(limit))
	}
	if skip != -1 {
		r.AddParameter("skip", strconv.Itoa(skip))
	}
	var envelope complaintsEnvelope
	err := r.MakeJSONRequest(&envelope)
	if err != nil {
		return -1, nil, err
	}
	return envelope.TotalCount, envelope.Items, nil
}

func (m *mailgunImpl) GetSingleComplaint(address string) (Complaint, error) {
	r := simplehttp.NewGetRequest(generateApiUrl(m, complaintsEndpoint) + "/" + address)
	var c Complaint
	err := r.MakeJSONRequest(&c)
	if err != nil {
		return Complaint{}, err
	}
	return c, nil
}
