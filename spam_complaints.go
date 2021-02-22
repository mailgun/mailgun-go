package mailgun

import (
	"context"
	"strconv"
)

const (
	complaintsEndpoint = "complaints"
)

// Complaint structures track how many times one of your emails have been marked as spam.
// the recipient thought your messages were not solicited.
type Complaint struct {
	Count     int         `json:"count"`
	CreatedAt RFC2822Time `json:"created_at"`
	Address   string      `json:"address"`
}

type complaintsResponse struct {
	Paging Paging      `json:"paging"`
	Items  []Complaint `json:"items"`
}

// ListComplaints returns a set of spam complaints registered against your domain.
// Recipients of your messages can click on a link which sends feedback to Mailgun
// indicating that the message they received is, to them, spam.
func (mg *MailgunImpl) ListComplaints(opts *ListOptions) *ComplaintsIterator {
	r := newHTTPRequest(generateApiUrl(mg, complaintsEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	if opts != nil {
		if opts.Limit != 0 {
			r.addParameter("limit", strconv.Itoa(opts.Limit))
		}
	}
	url, err := r.generateUrlWithParameters()
	return &ComplaintsIterator{
		mg:                 mg,
		complaintsResponse: complaintsResponse{Paging: Paging{Next: url, First: url}},
		err:                err,
	}
}

type ComplaintsIterator struct {
	complaintsResponse
	mg  Mailgun
	err error
}

// If an error occurred during iteration `Err()` will return non nil
func (ci *ComplaintsIterator) Err() error {
	return ci.err
}

// Next retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ci *ComplaintsIterator) Next(ctx context.Context, items *[]Complaint) bool {
	if ci.err != nil {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.Next)
	if ci.err != nil {
		return false
	}
	cpy := make([]Complaint, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	if len(ci.Items) == 0 {
		return false
	}
	return true
}

// First retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (ci *ComplaintsIterator) First(ctx context.Context, items *[]Complaint) bool {
	if ci.err != nil {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.First)
	if ci.err != nil {
		return false
	}
	cpy := make([]Complaint, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ci *ComplaintsIterator) Last(ctx context.Context, items *[]Complaint) bool {
	if ci.err != nil {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.Last)
	if ci.err != nil {
		return false
	}
	cpy := make([]Complaint, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ci *ComplaintsIterator) Previous(ctx context.Context, items *[]Complaint) bool {
	if ci.err != nil {
		return false
	}
	if ci.Paging.Previous == "" {
		return false
	}
	ci.err = ci.fetch(ctx, ci.Paging.Previous)
	if ci.err != nil {
		return false
	}
	cpy := make([]Complaint, len(ci.Items))
	copy(cpy, ci.Items)
	*items = cpy
	if len(ci.Items) == 0 {
		return false
	}
	return true
}

func (ci *ComplaintsIterator) fetch(ctx context.Context, url string) error {
	ci.Items = nil
	r := newHTTPRequest(url)
	r.setClient(ci.mg.Client())
	r.setBasicAuth(basicAuthUser, ci.mg.APIKey())

	return getResponseFromJSON(ctx, r, &ci.complaintsResponse)
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
