package mailgun

import (
	"context"
	"strconv"
)

// A mailing list may have one of three membership modes.
const (
	// ReadOnly specifies that nobody, including Members, may send messages to
	// the mailing list.  Messages distributed on such lists come from list
	// administrator accounts only.
	AccessLevelReadOnly = "readonly"
	// Members specifies that only those who subscribe to the mailing list may
	// send messages.
	AccessLevelMembers = "members"
	// Everyone specifies that anyone and everyone may both read and submit
	// messages to the mailing list, including non-subscribers.
	AccessLevelEveryone = "everyone"
)

// Specify the access of a mailing list member
type AccessLevel string

// A List structure provides information for a mailing list.
//
// AccessLevel may be one of ReadOnly, Members, or Everyone.
type MailingList struct {
	Address      string      `json:"address,omitempty"`
	Name         string      `json:"name,omitempty"`
	Description  string      `json:"description,omitempty"`
	AccessLevel  AccessLevel `json:"access_level,omitempty"`
	CreatedAt    RFC2822Time `json:"created_at,omitempty"`
	MembersCount int         `json:"members_count,omitempty"`
}

type listsResponse struct {
	Items  []MailingList `json:"items"`
	Paging Paging        `json:"paging"`
}

type mailingListResponse struct {
	MailingList MailingList `json:"member"`
}

type ListsIterator struct {
	listsResponse
	mg  Mailgun
	err error
}

// ListMailingLists returns the specified set of mailing lists administered by your account.
func (mg *MailgunImpl) ListMailingLists(opts *ListOptions) *ListsIterator {
	r := newHTTPRequest(generatePublicApiUrl(mg, listsEndpoint) + "/pages")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	if opts != nil {
		if opts.Limit != 0 {
			r.addParameter("limit", strconv.Itoa(opts.Limit))
		}
	}
	url, err := r.generateUrlWithParameters()
	return &ListsIterator{
		mg:            mg,
		listsResponse: listsResponse{Paging: Paging{Next: url, First: url}},
		err:           err,
	}
}

// If an error occurred during iteration `Err()` will return non nil
func (li *ListsIterator) Err() error {
	return li.err
}

// Next retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (li *ListsIterator) Next(ctx context.Context, items *[]MailingList) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Next)
	if li.err != nil {
		return false
	}
	cpy := make([]MailingList, len(li.Items))
	copy(cpy, li.Items)
	*items = cpy
	if len(li.Items) == 0 {
		return false
	}
	return true
}

// First retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (li *ListsIterator) First(ctx context.Context, items *[]MailingList) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.First)
	if li.err != nil {
		return false
	}
	cpy := make([]MailingList, len(li.Items))
	copy(cpy, li.Items)
	*items = cpy
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (li *ListsIterator) Last(ctx context.Context, items *[]MailingList) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Last)
	if li.err != nil {
		return false
	}
	cpy := make([]MailingList, len(li.Items))
	copy(cpy, li.Items)
	*items = cpy
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (li *ListsIterator) Previous(ctx context.Context, items *[]MailingList) bool {
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
	cpy := make([]MailingList, len(li.Items))
	copy(cpy, li.Items)
	*items = cpy
	if len(li.Items) == 0 {
		return false
	}
	return true
}

func (li *ListsIterator) fetch(ctx context.Context, url string) error {
	li.Items = nil
	r := newHTTPRequest(url)
	r.setClient(li.mg.Client())
	r.setBasicAuth(basicAuthUser, li.mg.APIKey())

	return getResponseFromJSON(ctx, r, &li.listsResponse)
}

// CreateMailingList creates a new mailing list under your Mailgun account.
// You need specify only the Address and Name members of the prototype;
// Description, and AccessLevel are optional.
// If unspecified, Description remains blank,
// while AccessLevel defaults to Everyone.
func (mg *MailgunImpl) CreateMailingList(ctx context.Context, prototype MailingList) (MailingList, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, listsEndpoint))
	r.setClient(mg.Client())
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
	response, err := makePostRequest(ctx, r, p)
	if err != nil {
		return MailingList{}, err
	}
	var l MailingList
	err = response.parseFromJSON(&l)
	return l, err
}

// DeleteMailingList removes all current members of the list, then removes the list itself.
// Attempts to send e-mail to the list will fail subsequent to this call.
func (mg *MailgunImpl) DeleteMailingList(ctx context.Context, addr string) error {
	r := newHTTPRequest(generatePublicApiUrl(mg, listsEndpoint) + "/" + addr)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// GetMailingList allows your application to recover the complete List structure
// representing a mailing list, so long as you have its e-mail address.
func (mg *MailgunImpl) GetMailingList(ctx context.Context, addr string) (MailingList, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, listsEndpoint) + "/" + addr)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	response, err := makeGetRequest(ctx, r)
	if err != nil {
		return MailingList{}, err
	}

	var resp mailingListResponse
	err = response.parseFromJSON(&resp)
	return resp.MailingList, err
}

// UpdateMailingList allows you to change various attributes of a list.
// Address, Name, Description, and AccessLevel are all optional;
// only those fields which are set in the prototype will change.
//
// Be careful!  If changing the address of a mailing list,
// e-mail sent to the old address will not succeed.
// Make sure you account for the change accordingly.
func (mg *MailgunImpl) UpdateMailingList(ctx context.Context, addr string, prototype MailingList) (MailingList, error) {
	r := newHTTPRequest(generatePublicApiUrl(mg, listsEndpoint) + "/" + addr)
	r.setClient(mg.Client())
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
	var l MailingList
	response, err := makePutRequest(ctx, r, p)
	if err != nil {
		return l, err
	}
	err = response.parseFromJSON(&l)
	return l, err
}
