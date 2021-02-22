package mailgun

import (
	"context"
	"net/url"
	"strconv"
	"time"
)

type Tag struct {
	Value       string     `json:"tag"`
	Description string     `json:"description"`
	FirstSeen   *time.Time `json:"first-seen,omitempty"`
	LastSeen    *time.Time `json:"last-seen,omitempty"`
}

type tagsResponse struct {
	Items  []Tag  `json:"items"`
	Paging Paging `json:"paging"`
}

type ListTagOptions struct {
	// Restrict the page size to this limit
	Limit int
	// Return only the tags starting with the given prefix
	Prefix string
}

// DeleteTag removes all counters for a particular tag, including the tag itself.
func (mg *MailgunImpl) DeleteTag(ctx context.Context, tag string) error {
	r := newHTTPRequest(generateApiUrl(mg, tagsEndpoint) + "/" + tag)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	_, err := makeDeleteRequest(ctx, r)
	return err
}

// GetTag retrieves metadata about the tag from the api
func (mg *MailgunImpl) GetTag(ctx context.Context, tag string) (Tag, error) {
	r := newHTTPRequest(generateApiUrl(mg, tagsEndpoint) + "/" + tag)
	r.setClient(mg.Client())
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	var tagItem Tag
	return tagItem, getResponseFromJSON(ctx, r, &tagItem)
}

// ListTags returns a cursor used to iterate through a list of tags
//	it := mg.ListTags(nil)
//	var page []mailgun.Tag
//	for it.Next(&page) {
//		for _, tag := range(page) {
//			// Do stuff with tags
//		}
//	}
//	if it.Err() != nil {
//		log.Fatal(it.Err())
//	}
func (mg *MailgunImpl) ListTags(opts *ListTagOptions) *TagIterator {
	req := newHTTPRequest(generateApiUrl(mg, tagsEndpoint))
	if opts != nil {
		if opts.Limit != 0 {
			req.addParameter("limit", strconv.Itoa(opts.Limit))
		}
		if opts.Prefix != "" {
			req.addParameter("prefix", opts.Prefix)
		}
	}

	url, err := req.generateUrlWithParameters()
	return &TagIterator{
		tagsResponse: tagsResponse{Paging: Paging{Next: url, First: url}},
		err:          err,
		mg:           mg,
	}
}

type TagIterator struct {
	tagsResponse
	mg  Mailgun
	err error
}

// Next returns the next page in the list of tags
func (ti *TagIterator) Next(ctx context.Context, items *[]Tag) bool {
	if ti.err != nil {
		return false
	}

	if !canFetchPage(ti.Paging.Next) {
		return false
	}

	ti.err = ti.fetch(ctx, ti.Paging.Next)
	if ti.err != nil {
		return false
	}
	*items = ti.Items
	if len(ti.Items) == 0 {
		return false
	}
	return true
}

// Previous returns the previous page in the list of tags
func (ti *TagIterator) Previous(ctx context.Context, items *[]Tag) bool {
	if ti.err != nil {
		return false
	}

	if ti.Paging.Previous == "" {
		return false
	}

	if !canFetchPage(ti.Paging.Previous) {
		return false
	}

	ti.err = ti.fetch(ctx, ti.Paging.Previous)
	if ti.err != nil {
		return false
	}
	*items = ti.Items
	if len(ti.Items) == 0 {
		return false
	}
	return true
}

// First returns the first page in the list of tags
func (ti *TagIterator) First(ctx context.Context, items *[]Tag) bool {
	if ti.err != nil {
		return false
	}
	ti.err = ti.fetch(ctx, ti.Paging.First)
	if ti.err != nil {
		return false
	}
	*items = ti.Items
	return true
}

// Last returns the last page in the list of tags
func (ti *TagIterator) Last(ctx context.Context, items *[]Tag) bool {
	if ti.err != nil {
		return false
	}
	ti.err = ti.fetch(ctx, ti.Paging.Last)
	if ti.err != nil {
		return false
	}
	*items = ti.Items
	return true
}

// Err returns any error if one occurred
func (ti *TagIterator) Err() error {
	return ti.err
}

func (ti *TagIterator) fetch(ctx context.Context, url string) error {
	ti.Items = nil
	req := newHTTPRequest(url)
	req.setClient(ti.mg.Client())
	req.setBasicAuth(basicAuthUser, ti.mg.APIKey())
	return getResponseFromJSON(ctx, req, &ti.tagsResponse)
}

func canFetchPage(slug string) bool {
	parts, err := url.Parse(slug)
	if err != nil {
		return false
	}
	params, _ := url.ParseQuery(parts.RawQuery)
	if err != nil {
		return false
	}
	value, ok := params["tag"]
	// If tags doesn't exist, it's our first time fetching pages
	if !ok {
		return true
	}
	// If tags has no value, there are no more pages to fetch
	return len(value) == 0
}
