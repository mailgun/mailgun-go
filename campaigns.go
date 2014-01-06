package mailgun

import (
	"github.com/mbanzon/simplehttp"
)

type Campaign struct {
	Id                string `json:"id"`
	Name              string `json:"name"`
	CreatedAt         string `json:"created_at"`
	DeliveredCount    int    `json:"delivered_count"`
	ClickedCount      int    `json:"clicked_count"`
	OpenedCount       int    `json:"opened_count"`
	SubmittedCount    int    `json:"submitted_count"`
	UnsubscribedCount int    `json:"unsubscribed_count"`
	BouncedCount      int    `json:"bounced_count"`
	ComplainedCount   int    `json:"complained_count"`
	DroppedCount      int    `json:"dropped_count"`
}

type campaignsEnvelope struct {
	TotalCount int        `json:"total_count"`
	Items      []Campaign `json:"items"`
}

func (m *mailgunImpl) GetCampaigns() (int, []Campaign, error) {
	r := simplehttp.NewGetRequest(generateApiUrl(m, campaignsEndpoint))
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	var envelope campaignsEnvelope
	err := r.MakeJSONRequest(&envelope)
	if err != nil {
		return -1, nil, err
	}
	return envelope.TotalCount, envelope.Items, nil
}

func (m *mailgunImpl) CreateCampaign(name, id string) error {
	r := simplehttp.NewPostRequest(generateApiUrl(m, campaignsEndpoint))
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	r.AddFormValue("name", name)
	if id != "" {
		r.AddFormValue("id", id)
	}
	_, err := r.MakeRequest()
	return err
}

func (m *mailgunImpl) UpdateCampaign(oldId, name, newId string) error {
	r := simplehttp.NewPostRequest(generateApiUrl(m, campaignsEndpoint) + "/" + oldId)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	r.AddFormValue("name", name)
	if newId != "" {
		r.AddFormValue("id", newId)
	}
	_, err := r.MakeRequest()
	return err
}

func (m *mailgunImpl) DeleteCampaign(id string) error {
	r := simplehttp.NewDeleteRequest(generateApiUrl(m, campaignsEndpoint) + "/" + id)
	r.SetBasicAuth(basicAuthUser, m.ApiKey())
	_, err := r.MakeRequest()
	return err
}
