package mailgun

// https://documentation.mailgun.com/docs/inboxready/openapi-final/tag/Domains/

import (
	"context"
	"strconv"

	"github.com/mailgun/errors"
	"github.com/mailgun/mailgun-go/v5/internal/types/inboxready"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

type ListMonitoredDomainsOptions = inboxready.GETV1InboxreadyDomainsParams

type MonitoredDomainsIterator struct {
	mg      Mailgun
	opts    inboxready.GETV1InboxreadyDomainsParams
	req     *httpRequest
	resp    inboxready.InboxReadyGithubComMailgunInboxreadyAPIDomainListResponse
	isFirst bool
	err     error
}

func (mg *Client) ListMonitoredDomains(opts ListMonitoredDomainsOptions) (*MonitoredDomainsIterator, error) {
	// TODO(vtopc): support opts.Domain

	if opts.Limit == nil || *opts.Limit == 0 {
		opts.Limit = ptr(10)
	}

	req := newHTTPRequest(generateApiUrl(mg, 1, inboxreadyDomainsEndpoint))
	req.addParameter("limit", strconv.Itoa(*opts.Limit))
	req.setClient(mg.HTTPClient())
	req.setBasicAuth(basicAuthUser, mg.APIKey())

	return &MonitoredDomainsIterator{
		opts:    opts,
		req:     req,
		isFirst: true,
	}, nil
}

func (iter *MonitoredDomainsIterator) Err() error {
	return iter.err
}

// Next retrieves the next page of items from the api. Returns false when there are
// no more pages to retrieve or if there was an error.
// Use `.Err()` to retrieve the error
func (iter *MonitoredDomainsIterator) Next(ctx context.Context, resp []mtypes.MonitoredDomain) (more bool) {
	if iter.err != nil {
		return false
	}

	if iter.isFirst {
		iter.isFirst = false
	} else {
		iter.req.URL = iter.resp.Paging.Next
	}

	v, err := iter.fetch(ctx)
	if err != nil {
		iter.err = err

		return false
	}

	iter.resp = *v

	return len(resp) == *iter.opts.Limit
}

func (iter *MonitoredDomainsIterator) fetch(ctx context.Context,
) (*inboxready.InboxReadyGithubComMailgunInboxreadyAPIDomainListResponse, error) {
	httpResp, err := makeGetRequest(ctx, iter.req)
	if err != nil {
		return nil, err
	}

	var resp inboxready.InboxReadyGithubComMailgunInboxreadyAPIDomainListResponse
	err = httpResp.parseFromJSON(&resp)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}

	return &resp, nil
}

// TODO:
// AddDomainToMonitoring
// DeleteMonitoredDomain
