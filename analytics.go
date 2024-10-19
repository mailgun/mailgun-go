package mailgun

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

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

type MetricsPagination struct {
	// Colon-separated value indicating column name and sort direction e.g. 'domain:asc'.
	Sort string `json:"sort"`
	// The number of items to skip over when satisfying the request. To get the first page of data set skip to zero. Then increment the skip by the limit for subsequent calls.
	Skip int64 `json:"skip"`
	// The maximum number of items returned in the response.
	Limit int64 `json:"limit"`
	// The total number of items in the query result set.
	Total int64 `json:"total"`
}

type MetricsIterator struct {
	opts MetricsOptions
	req  *httpRequest
	resp MetricsResponse
	err  error
}

func (iter *MetricsIterator) Err() error {
	return iter.err
}

// Next retrieves the next page of items from the api. Returns false when there are
// no more pages to retrieve or if there was an error.
// Use `.Err()` to retrieve the error
func (iter *MetricsIterator) Next(ctx context.Context, resp *MetricsResponse) bool {
	if iter.err != nil {
		return false
	}

	iter.err = iter.fetch(ctx)
	if iter.err != nil {
		return false
	}

	*resp = iter.resp
	if len(iter.resp.Items) == 0 {
		return false
	}
	iter.opts.Pagination.Skip = iter.opts.Pagination.Skip + iter.opts.Pagination.Limit

	return true
}

func (iter *MetricsIterator) fetch(ctx context.Context) error {
	payload := newJSONEncodedPayload(iter.opts)

	httpResp, err := makePostRequest(ctx, iter.req, payload)
	if err != nil {
		return err
	}

	var resp MetricsResponse
	err = httpResp.parseFromJSON(&resp)
	if err != nil {
		return errors.Wrap(err, "decoding response")
	}

	iter.resp = resp

	return nil
}