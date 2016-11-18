package mailgun

// Campaigns have been deprecated since development work on this SDK commenced.
// Please refer to http://documentation.mailgun.com/api_reference .
type Campaign struct {
	ID                string `json:"id"`
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

// Campaigns have been deprecated since development work on this SDK commenced.
// Please refer to http://documentation.mailgun.com/api_reference .
func (m *Impl) GetCampaigns() (int, []Campaign, error) {
	r := newHTTPRequest(generateAPIUrl(m, campaignsEndpoint))
	r.setClient(m.Client())
	r.setBasicAuth(basicAuthUser, m.APIKey())

	var envelope campaignsEnvelope
	err := getResponseFromJSON(r, &envelope)
	if err != nil {
		return -1, nil, err
	}
	return envelope.TotalCount, envelope.Items, nil
}

// Campaigns have been deprecated since development work on this SDK commenced.
// Please refer to http://documentation.mailgun.com/api_reference .
func (m *Impl) CreateCampaign(name, id string) error {
	r := newHTTPRequest(generateAPIUrl(m, campaignsEndpoint))
	r.setClient(m.Client())
	r.setBasicAuth(basicAuthUser, m.APIKey())

	payload := newURLEncodedPayload()
	payload.addValue("name", name)
	if id != "" {
		payload.addValue("id", id)
	}
	_, err := makePostRequest(r, payload)
	return err
}

// Campaigns have been deprecated since development work on this SDK commenced.
// Please refer to http://documentation.mailgun.com/api_reference .
func (m *Impl) UpdateCampaign(oldID, name, newID string) error {
	r := newHTTPRequest(generateAPIUrl(m, campaignsEndpoint) + "/" + oldID)
	r.setClient(m.Client())
	r.setBasicAuth(basicAuthUser, m.APIKey())

	payload := newURLEncodedPayload()
	payload.addValue("name", name)
	if newID != "" {
		payload.addValue("id", newID)
	}
	_, err := makePostRequest(r, payload)
	return err
}

// Campaigns have been deprecated since development work on this SDK commenced.
// Please refer to http://documentation.mailgun.com/api_reference .
func (m *Impl) DeleteCampaign(id string) error {
	r := newHTTPRequest(generateAPIUrl(m, campaignsEndpoint) + "/" + id)
	r.setClient(m.Client())
	r.setBasicAuth(basicAuthUser, m.APIKey())
	_, err := makeDeleteRequest(r)
	return err
}
