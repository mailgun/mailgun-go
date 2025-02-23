package mailgun

import (
	"context"
	"strconv"

	"github.com/mailgun/mailgun-go/v4/mtypes"
)

type ListsIterator struct {
	mtypes.ListMailingListsResponse
	mg  Mailgun
	err error
}

// ListMailingLists returns the specified set of mailing lists administered by your account.
func (mg *MailgunImpl) ListMailingLists(opts *ListOptions) *ListsIterator {
	r := newHTTPRequest(generateApiUrl(mg, 3, listsEndpoint) + "/pages")
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	if opts != nil {
		if opts.Limit != 0 {
			r.addParameter("limit", strconv.Itoa(opts.Limit))
		}
	}
	url, err := r.generateUrlWithParameters()
	return &ListsIterator{
		mg:                       mg,
		ListMailingListsResponse: mtypes.ListMailingListsResponse{Paging: mtypes.Paging{Next: url, First: url}},
		err:                      err,
	}
}

// If an error occurred during iteration `Err()` will return non nil
func (li *ListsIterator) Err() error {
	return li.err
}

// Next retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (li *ListsIterator) Next(ctx context.Context, items *[]mtypes.MailingList) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Next)
	if li.err != nil {
		return false
	}
	cpy := make([]mtypes.MailingList, len(li.Items))
	copy(cpy, li.Items)
	*items = cpy

	return len(li.Items) != 0
}

// First retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (li *ListsIterator) First(ctx context.Context, items *[]mtypes.MailingList) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.First)
	if li.err != nil {
		return false
	}
	cpy := make([]mtypes.MailingList, len(li.Items))
	copy(cpy, li.Items)
	*items = cpy
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (li *ListsIterator) Last(ctx context.Context, items *[]mtypes.MailingList) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Last)
	if li.err != nil {
		return false
	}
	cpy := make([]mtypes.MailingList, len(li.Items))
	copy(cpy, li.Items)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (li *ListsIterator) Previous(ctx context.Context, items *[]mtypes.MailingList) bool {
	if li.err != nil {
		return false
	}
	if li.Paging.Previous == "" {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Previous)
	if li.err != nil {
		return false
	}
	cpy := make([]mtypes.MailingList, len(li.Items))
	copy(cpy, li.Items)
	*items = cpy

	return len(li.Items) != 0
}

func (li *ListsIterator) fetch(ctx context.Context, url string) error {
	li.Items = nil
	r := newHTTPRequest(url)
	r.setClient(li.mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, li.mg.APIKey())

	return getResponseFromJSON(ctx, r, &li.ListMailingListsResponse)
}

// CreateMailingList creates a new mailing list under your Mailgun account.
// You need specify only the Address and Name members of the prototype;
// Description, AccessLevel and ReplyPreference are optional.
// If unspecified, Description remains blank,
// while AccessLevel defaults to Everyone
// and ReplyPreference defaults to List.
func (mg *MailgunImpl) CreateMailingList(ctx context.Context, prototype mtypes.MailingList) (mtypes.MailingList, error) {
	r := newHTTPRequest(generateApiUrl(mg, 3, listsEndpoint))
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	if prototype.Address != "" {
		p.addValue("address", prototype.Address)
	}
	if prototype.Name != "" {
		p.addValue("name", prototype.Name)
	}
	if prototype.Description != "" {
		p.addValue("description", prototype.Description)
	}
	if prototype.AccessLevel != "" {
		p.addValue("access_level", string(prototype.AccessLevel))
	}
	if prototype.ReplyPreference != "" {
		p.addValue("reply_preference", string(prototype.ReplyPreference))
	}
	response, err := makePostRequest(ctx, r, p)
	if err != nil {
		return mtypes.MailingList{}, err
	}
	var l mtypes.MailingList
	err = response.parseFromJSON(&l)
	return l, err
}

// DeleteMailingList removes all current members of the list, then removes the list itself.
// Attempts to send e-mail to the list will fail subsequent to this call.
func (mg *MailgunImpl) DeleteMailingList(ctx context.Context, addr string) error {
	r := newHTTPRequest(generateApiUrl(mg, 3, listsEndpoint) + "/" + addr)
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// GetMailingList allows your application to recover the complete List structure
// representing a mailing list, so long as you have its e-mail address.
func (mg *MailgunImpl) GetMailingList(ctx context.Context, addr string) (mtypes.MailingList, error) {
	r := newHTTPRequest(generateApiUrl(mg, 3, listsEndpoint) + "/" + addr)
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	response, err := makeGetRequest(ctx, r)
	if err != nil {
		return mtypes.MailingList{}, err
	}

	var resp mtypes.GetMailingListResponse
	err = response.parseFromJSON(&resp)
	return resp.MailingList, err
}

// UpdateMailingList allows you to change various attributes of a list.
// Address, Name, Description, AccessLevel and ReplyPreference are all optional;
// only those fields which are set in the prototype will change.
//
// Be careful!  If changing the address of a mailing list,
// e-mail sent to the old address will not succeed.
// Make sure you account for the change accordingly.
func (mg *MailgunImpl) UpdateMailingList(ctx context.Context, addr string, prototype mtypes.MailingList) (mtypes.MailingList, error) {
	r := newHTTPRequest(generateApiUrl(mg, 3, listsEndpoint) + "/" + addr)
	r.setClient(mg.HTTPClient())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newUrlEncodedPayload()
	if prototype.Address != "" {
		p.addValue("address", prototype.Address)
	}
	if prototype.Name != "" {
		p.addValue("name", prototype.Name)
	}
	if prototype.Description != "" {
		p.addValue("description", prototype.Description)
	}
	if prototype.AccessLevel != "" {
		p.addValue("access_level", string(prototype.AccessLevel))
	}
	if prototype.ReplyPreference != "" {
		p.addValue("reply_preference", string(prototype.ReplyPreference))
	}
	var l mtypes.MailingList
	response, err := makePutRequest(ctx, r, p)
	if err != nil {
		return l, err
	}
	err = response.parseFromJSON(&l)
	return l, err
}
