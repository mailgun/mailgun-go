package mailgun

import (
	"context"
	"iter"
	"strconv"

	"github.com/mailgun/mailgun-go/v5/mtypes"
)

// ListDomainsIter returns an iterator over pages of domains.
// Iteration stops after the first error or after fetching all items.
// https://documentation.mailgun.com/docs/mailgun/api-reference/send/mailgun/domains/get-v4-domains
func (mg *Client) ListDomainsIter(ctx context.Context, opts *ListDomainsOptions) iter.Seq2[[]mtypes.Domain, error] {
	limit := 100
	if opts != nil && opts.Limit != 0 {
		limit = opts.Limit
	}

	url := generateApiUrl(mg, 4, domainsEndpoint)

	return func(yield func([]mtypes.Domain, error) bool) {
		var skip int
		for {
			resp, err := mg.fetchDomains(ctx, url, opts, skip, limit)
			if err != nil {
				yield(nil, err)
				return
			}

			if len(resp.Items) == 0 {
				return
			}

			if !yield(resp.Items, nil) {
				return
			}

			skip += len(resp.Items)
			if skip >= resp.TotalCount {
				return
			}
		}
	}
}

func (mg *Client) fetchDomains(ctx context.Context, url string, opts *ListDomainsOptions, skip, limit int,
) (mtypes.ListDomainsResponse, error) {
	r := newHTTPRequest(url)
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.setClient(mg.HTTPClient())

	if skip != 0 {
		r.addParameter("skip", strconv.Itoa(skip))
	}

	if limit != 0 {
		r.addParameter("limit", strconv.Itoa(limit))
	}

	if opts != nil {
		if opts.State != nil {
			r.addParameter("state", string(*opts.State))
		}
		if opts.Sort != nil {
			r.addParameter("sort", *opts.Sort)
		}
		if opts.Authority != nil {
			r.addParameter("authority", *opts.Authority)
		}
		if opts.Search != nil {
			r.addParameter("search", *opts.Search)
		}
		if opts.IncludeSubaccounts != nil {
			r.addParameter("include_subaccounts", strconv.FormatBool(*opts.IncludeSubaccounts))
		}
	}

	var resp mtypes.ListDomainsResponse
	err := getResponseFromJSON(ctx, r, &resp)

	return resp, err
}
