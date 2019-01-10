package mailgun

import (
	"context"
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

// ListComplaints returns a set of spam complaints registered against your domain.
// Recipients of your messages can click on a link which sends feedback to Mailgun
// indicating that the message they received is, to them, spam.
func (mg *MailgunImpl) ListComplaints(ctx context.Context, opts *ListOptions) ([]Complaint, error) {
	r := newHTTPRequest(generateApiUrl(mg, complaintsEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	if opts != nil && opts.Limit != 0 {
		r.addParameter("limit", strconv.Itoa(opts.Limit))
	}

	if opts != nil && opts.Skip != 0 {
		r.addParameter("skip", strconv.Itoa(opts.Skip))
	}

	var envelope complaintsEnvelope
	err := getResponseFromJSON(ctx, r, &envelope)
	if err != nil {
		return nil, err
	}
	return envelope.Items, nil
}

// GetComplaint returns a single complaint record filed by a recipient at the email address provided.
// If no complaint exists, the Complaint instance returned will be empty.
func (mg *MailgunImpl) GetComplaint(ctx context.Context, address string) (Complaint, error) {
	r := newHTTPRequest(generateApiUrl(mg, complaintsEndpoint) + "/" + address)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())

	var c Complaint
	err := getResponseFromJSON(ctx, r, &c)
	return c, err
}

// CreateComplaint registers the specified address as a recipient who has complained of receiving spam
// from your domain.
func (mg *MailgunImpl) CreateComplaint(ctx context.Context, address string) error {
	r := newHTTPRequest(generateApiUrl(mg, complaintsEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	p.addValue("address", address)
	_, err := makePostRequest(ctx, r, p)
	return err
}

// DeleteComplaint removes a previously registered e-mail address from the list of people who complained
// of receiving spam from your domain.
func (mg *MailgunImpl) DeleteComplaint(ctx context.Context, address string) error {
	r := newHTTPRequest(generateApiUrl(mg, complaintsEndpoint) + "/" + address)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}
