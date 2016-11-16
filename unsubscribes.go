package mailgun

import (
	"strconv"
)

type Unsubscription struct {
	CreatedAt string `json:"created_at"`
	Tag       string `json:"tag"`
	ID        string `json:"id"`
	Address   string `json:"address"`
}

// GetUnsubscribes retrieves a list of unsubscriptions issued by recipients of mail from your domain.
// Zero is a valid list length.
func (mg *Impl) GetUnsubscribes(limit, skip int) (int, []Unsubscription, error) {
	r := newHTTPRequest(generateAPIUrl(mg, unsubscribesEndpoint))
	if limit != DefaultLimit {
		r.addParameter("limit", strconv.Itoa(limit))
	}
	if skip != DefaultSkip {
		r.addParameter("skip", strconv.Itoa(skip))
	}
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var envelope struct {
		TotalCount int              `json:"total_count"`
		Items      []Unsubscription `json:"items"`
	}
	err := getResponseFromJSON(r, &envelope)
	return envelope.TotalCount, envelope.Items, err
}

// GetUnsubscribesByAddress retrieves a list of unsubscriptions by recipient address.
// Zero is a valid list length.
func (mg *Impl) GetUnsubscribesByAddress(a string) (int, []Unsubscription, error) {
	r := newHTTPRequest(generateAPIUrlWithTarget(mg, unsubscribesEndpoint, a))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var envelope struct {
		TotalCount int              `json:"total_count"`
		Items      []Unsubscription `json:"items"`
	}
	err := getResponseFromJSON(r, &envelope)
	return envelope.TotalCount, envelope.Items, err
}

// Unsubscribe adds an e-mail address to the domain's unsubscription table.
func (mg *Impl) Unsubscribe(a, t string) error {
	r := newHTTPRequest(generateAPIUrl(mg, unsubscribesEndpoint))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newURLEncodedPayload()
	p.addValue("address", a)
	p.addValue("tag", t)
	_, err := makePostRequest(r, p)
	return err
}

// RemoveUnsubscribe removes the e-mail address given from the domain's unsubscription table.
// If passing in an ID (discoverable from, e.g., GetUnsubscribes()), the e-mail address associated
// with the given ID will be removed.
func (mg *Impl) RemoveUnsubscribe(a string) error {
	r := newHTTPRequest(generateAPIUrlWithTarget(mg, unsubscribesEndpoint, a))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(r)
	return err
}

// RemoveUnsubscribeWithTag removes the e-mail address given from the domain's unsubscription table with a matching tag.
// If passing in an ID (discoverable from, e.g., GetUnsubscribes()), the e-mail address associated
// with the given ID will be removed.
func (mg *Impl) RemoveUnsubscribeWithTag(a, t string) error {
	r := newHTTPRequest(generateAPIUrlWithTarget(mg, unsubscribesEndpoint, a))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.addParameter("tag", t)
	_, err := makeDeleteRequest(r)
	return err
}
