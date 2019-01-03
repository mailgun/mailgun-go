package mailgun

import (
	"strconv"
)

const (
	complaintsEndpoint = "complaints"
)

// Complaint structures track how many times one of your emails have been marked as spam.
// CreatedAt indicates when the first report arrives from a given recipient, identified by Address.
// Count provides a running counter of how many times
// the recipient thought your messages were not solicited.
type Complaint struct {
	Count     int    `json:"count"`
	CreatedAt string `json:"created_at"`
	Address   string `json:"address"`
}

type complaintsEnvelope struct {
	Items  []Complaint `json:"items"`
	Paging Paging      `json:"paging"`
}

// GetComplaints returns a set of spam complaints registered against your domain.
// Recipients of your messages can click on a link which sends feedback to Mailgun
// indicating that the message they received is, to them, spam.
func (mg *MailgunImpl) GetComplaints(limit, skip int) (int, []Complaint, error) {
	r := newHTTPRequest(generateApiUrl(mg, complaintsEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	if limit != -1 {
		r.addParameter("limit", strconv.Itoa(limit))
	}
	if skip != -1 {
		r.addParameter("skip", strconv.Itoa(skip))
	}

	var envelope complaintsEnvelope
	err := getResponseFromJSON(r, &envelope)
	if err != nil {
		return -1, nil, err
	}
	return len(envelope.Items), envelope.Items, nil
}

// GetSingleComplaint returns a single complaint record filed by a recipient at the email address provided.
// If no complaint exists, the Complaint instance returned will be empty.
func (mg *MailgunImpl) GetSingleComplaint(address string) (Complaint, error) {
	r := newHTTPRequest(generateApiUrl(mg, complaintsEndpoint) + "/" + address)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var c Complaint
	err := getResponseFromJSON(r, &c)
	return c, err
}

// CreateComplaint registers the specified address as a recipient who has complained of receiving spam
// from your domain.
func (mg *MailgunImpl) CreateComplaint(address string) error {
	r := newHTTPRequest(generateApiUrl(mg, complaintsEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("address", address)
	_, err := makePostRequest(r, p)
	return err
}

// DeleteComplaint removes a previously registered e-mail address from the list of people who complained
// of receiving spam from your domain.
func (mg *MailgunImpl) DeleteComplaint(address string) error {
	r := newHTTPRequest(generateApiUrl(mg, complaintsEndpoint) + "/" + address)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(r)
	return err
}
