package mailgun

// https://documentation.mailgun.com/docs/inboxready/openapi-final/tag/Domains/

import (
	"context"

	"github.com/mailgun/errors"
	"github.com/mailgun/mailgun-go/v5/internal/types/inboxready"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

type ListMonitoredDomainsOptions = inboxready.GETV1InboxreadyDomainsParams

type MonitoredDomainsIterator struct {
	opts inboxready.GETV1InboxreadyDomainsParams
	req  *httpRequest
	err  error
}

func (mg *Client) ListMonitoredDomain(opts ListMonitoredDomainsOptions) (*MonitoredDomainsIterator, error) {
	if opts.Pagination.Limit == 0 {
		opts.Pagination.Limit = 10
	}

	req := newHTTPRequest(generateApiUrl(mg, 1, inboxreadyDomainsEndpoint))
	req.setClient(mg.HTTPClient())
	req.setBasicAuth(basicAuthUser, mg.APIKey())

	return &MonitoredDomainsIterator{
		opts: opts,
		req:  req,
	}, nil
}

func (iter *MonitoredDomainsIterator) Err() error {
	return iter.err
}

// Next retrieves the next page of items from the api. Returns false when there are
// no more pages to retrieve or if there was an error.
// Use `.Err()` to retrieve the error
func (iter *MonitoredDomainsIterator) Next(ctx context.Context, resp *mtypes.MetricsResponse) (more bool) {
	if iter.err != nil {
		return false
	}

	iter.err = iter.fetch(ctx, resp)
	if iter.err != nil {
		return false
	}

	iter.opts.Pagination.Skip += iter.opts.Pagination.Limit

	return len(resp.Items) == iter.opts.Pagination.Limit
}

func (iter *MonitoredDomainsIterator) fetch(ctx context.Context, resp *mtypes.MetricsResponse) error {
	if resp == nil {
		return errors.New("resp cannot be nil")
	}

	payload := newJSONEncodedPayload(iter.opts)

	httpResp, err := makePostRequest(ctx, iter.req, payload)
	if err != nil {
		return err
	}

	// preallocate
	resp.Items = make([]mtypes.MetricsItem, 0, iter.opts.Pagination.Limit)

	err = httpResp.parseFromJSON(resp)
	if err != nil {
		return errors.Wrap(err, "decoding response")
	}

	return nil
}

// TODO:
// AddDomainToMonitoring
// DeleteMonitoredDomain
