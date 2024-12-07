package mailgun

import (
	"context"
	"strings"

	"github.com/mailgun/errors"
)

type MetricsPagination struct {
	// Colon-separated value indicating column name and sort direction e.g. 'domain:asc'.
	Sort string `json:"sort"`
	// The number of items to skip over when satisfying the request.
	// To get the first page of data set skip to zero.
	// Then increment the skip by the limit for subsequent calls.
	Skip int `json:"skip"`
	// The maximum number of items returned in the response.
	Limit int `json:"limit"`
	// The total number of items in the query result set.
	Total int `json:"total"`
}

// ListMetrics returns domain/account metrics.
//
// NOTE: Only for v1 API. To use the /v1 version define MG_URL in the environment variable
// as `https://api.mailgun.net/v1` or set `mg.SetAPIBase("https://api.mailgun.net/v1")`
//
// https://documentation.mailgun.com/docs/mailgun/api-reference/openapi-final/tag/Metrics/
func (mg *MailgunImpl) ListMetrics(opts MetricsOptions) (*MetricsIterator, error) {
	if !strings.HasSuffix(mg.APIBase(), "/v1") {
		return nil, errors.New("only v1 API is supported")
	}

	domain := mg.Domain()
	if domain != "" {
		domainFilter := MetricsFilterPredicate{
			Attribute:     "domain",
			Comparator:    "=",
			LabeledValues: []MetricsLabeledValue{{Label: domain, Value: domain}},
		}

		opts.Filter.BoolGroupAnd = append(opts.Filter.BoolGroupAnd, domainFilter)
	}

	if opts.Pagination.Limit == 0 {
		opts.Pagination.Limit = 10
	}

	req := newHTTPRequest(generatePublicApiUrl(mg, metricsEndpoint))
	req.setClient(mg.Client())
	req.setBasicAuth(basicAuthUser, mg.APIKey())

	return &MetricsIterator{
		opts: opts,
		req:  req,
	}, nil
}

type MetricsIterator struct {
	opts MetricsOptions
	req  *httpRequest
	err  error
}

func (iter *MetricsIterator) Err() error {
	return iter.err
}

// Next retrieves the next page of items from the api. Returns false when there are
// no more pages to retrieve or if there was an error.
// Use `.Err()` to retrieve the error
func (iter *MetricsIterator) Next(ctx context.Context, resp *MetricsResponse) (more bool) {
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

func (iter *MetricsIterator) fetch(ctx context.Context, resp *MetricsResponse) error {
	if resp == nil {
		return errors.New("resp cannot be nil")
	}

	payload := newJSONEncodedPayload(iter.opts)

	httpResp, err := makePostRequest(ctx, iter.req, payload)
	if err != nil {
		return err
	}

	// preallocate
	resp.Items = make([]MetricsItem, 0, iter.opts.Pagination.Limit)

	err = httpResp.parseFromJSON(resp)
	if err != nil {
		return errors.Wrap(err, "decoding response")
	}

	return nil
}
