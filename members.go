package mailgun

import (
	"context"
	"encoding/json"
	"strconv"
)

// yes and no are variables which provide us the ability to take their addresses.
// Subscribed and Unsubscribed are pointers to these booleans.
//
// We use a pointer to boolean as a kind of trinary data type:
// if nil, the relevant data type remains unspecified.
// Otherwise, its value is either true or false.
var (
	yes bool = true
	no  bool = false
)

// Mailing list members have an attribute that determines if they've subscribed to the mailing list or not.
// This attribute may be used to filter the results returned by GetSubscribers().
// All, Subscribed, and Unsubscribed provides a convenient and readable syntax for specifying the scope of the search.
var (
	All          *bool = nil
	Subscribed   *bool = &yes
	Unsubscribed *bool = &no
)

// A Member structure represents a member of the mailing list.
// The Vars field can represent any JSON-encodable data.
type Member struct {
	Address    string                 `json:"address,omitempty"`
	Name       string                 `json:"name,omitempty"`
	Subscribed *bool                  `json:"subscribed,omitempty"`
	Vars       map[string]interface{} `json:"vars,omitempty"`
}

type memberListResponse struct {
	Lists  []Member `json:"items"`
	Paging Paging   `json:"paging"`
}

type memberResponse struct {
	Member Member `json:"member"`
}

type MemberListIterator struct {
	memberListResponse
	mg  Mailgun
	err error
}

// Used by List methods to specify what list parameters to send to the mailgun API
type ListOptions struct {
	Limit int
}

func (mg *MailgunImpl) ListMembers(address string, opts *ListOptions) *MemberListIterator {
	r := newHTTPRequest(generateMemberApiUrl(mg, listsEndpoint, address) + "/pages")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	if opts != nil {
		if opts.Limit != 0 {
			r.addParameter("limit", strconv.Itoa(opts.Limit))
		}
	}
	url, err := r.generateUrlWithParameters()
	return &MemberListIterator{
		mg:                 mg,
		memberListResponse: memberListResponse{Paging: Paging{Next: url, First: url}},
		err:                err,
	}
}

// If an error occurred during iteration `Err()` will return non nil
func (li *MemberListIterator) Err() error {
	return li.err
}

// Next retrieves the next page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (li *MemberListIterator) Next(ctx context.Context, items *[]Member) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Next)
	if li.err != nil {
		return false
	}
	*items = li.Lists
	if len(li.Lists) == 0 {
		return false
	}
	return true
}

// First retrieves the first page of items from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (li *MemberListIterator) First(ctx context.Context, items *[]Member) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.First)
	if li.err != nil {
		return false
	}
	*items = li.Lists
	return true
}

// Last retrieves the last page of items from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (li *MemberListIterator) Last(ctx context.Context, items *[]Member) bool {
	if li.err != nil {
		return false
	}
	li.err = li.fetch(ctx, li.Paging.Last)
	if li.err != nil {
		return false
	}
	*items = li.Lists
	return true
}

// Previous retrieves the previous page of items from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (li *MemberListIterator) Previous(ctx context.Context, items *[]Member) bool {
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
	*items = li.Lists
	if len(li.Lists) == 0 {
		return false
	}
	return true
}

func (li *MemberListIterator) fetch(ctx context.Context, url string) error {
	r := newHTTPRequest(url)
	r.setClient(li.mg.Client())
	r.setBasicAuth(basicAuthUser, li.mg.APIKey())

	return getResponseFromJSON(ctx, r, &li.memberListResponse)
}

// GetMember returns a complete Member structure for a member of a mailing list,
// given only their subscription e-mail address.
func (mg *MailgunImpl) GetMember(ctx context.Context, s, l string) (Member, error) {
	r := newHTTPRequest(generateMemberApiUrl(mg, listsEndpoint, l) + "/" + s)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	response, err := makeGetRequest(ctx, r)
	if err != nil {
		return Member{}, err
	}
	var resp memberResponse
	err = response.parseFromJSON(&resp)
	return resp.Member, err
}

// CreateMember registers a new member of the indicated mailing list.
// If merge is set to true, then the registration may update an existing Member's settings.
// Otherwise, an error will occur if you attempt to add a member with a duplicate e-mail address.
func (mg *MailgunImpl) CreateMember(ctx context.Context, merge bool, addr string, prototype Member) error {
	vs, err := json.Marshal(prototype.Vars)
	if err != nil {
		return err
	}

	r := newHTTPRequest(generateMemberApiUrl(mg, listsEndpoint, addr))
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newFormDataPayload()
	p.addValue("upsert", yesNo(merge))
	p.addValue("address", prototype.Address)
	p.addValue("name", prototype.Name)
	p.addValue("vars", string(vs))
	if prototype.Subscribed != nil {
		p.addValue("subscribed", yesNo(*prototype.Subscribed))
	}
	_, err = makePostRequest(ctx, r, p)
	return err
}

// UpdateMember lets you change certain details about the indicated mailing list member.
// Address, Name, Vars, and Subscribed fields may be changed.
func (mg *MailgunImpl) UpdateMember(ctx context.Context, s, l string, prototype Member) (Member, error) {
	r := newHTTPRequest(generateMemberApiUrl(mg, listsEndpoint, l) + "/" + s)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newFormDataPayload()
	if prototype.Address != "" {
		p.addValue("address", prototype.Address)
	}
	if prototype.Name != "" {
		p.addValue("name", prototype.Name)
	}
	if prototype.Vars != nil {
		vs, err := json.Marshal(prototype.Vars)
		if err != nil {
			return Member{}, err
		}
		p.addValue("vars", string(vs))
	}
	if prototype.Subscribed != nil {
		p.addValue("subscribed", yesNo(*prototype.Subscribed))
	}
	response, err := makePutRequest(ctx, r, p)
	if err != nil {
		return Member{}, err
	}
	var envelope struct {
		Member Member `json:"member"`
	}
	err = response.parseFromJSON(&envelope)
	return envelope.Member, err
}

// DeleteMember removes the member from the list.
func (mg *MailgunImpl) DeleteMember(ctx context.Context, member, addr string) error {
	r := newHTTPRequest(generateMemberApiUrl(mg, listsEndpoint, addr) + "/" + member)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// CreateMemberList registers multiple Members and non-Member members to a single mailing list
// in a single round-trip.
// u indicates if the existing members should be updated or duplicates should be updated.
// Use All to elect not to provide a default.
// The newMembers list can take one of two JSON-encodable forms: an slice of strings, or
// a slice of Member structures.
// If a simple slice of strings is passed, each string refers to the member's e-mail address.
// Otherwise, each Member needs to have at least the Address field filled out.
// Other fields are optional, but may be set according to your needs.
func (mg *MailgunImpl) CreateMemberList(ctx context.Context, u *bool, addr string, newMembers []interface{}) error {
	r := newHTTPRequest(generateMemberApiUrl(mg, listsEndpoint, addr) + ".json")
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	p := newFormDataPayload()
	if u != nil {
		p.addValue("upsert", yesNo(*u))
	}
	bs, err := json.Marshal(newMembers)
	if err != nil {
		return err
	}
	p.addValue("members", string(bs))
	_, err = makePostRequest(ctx, r, p)
	return err
}
