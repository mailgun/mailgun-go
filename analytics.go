package mailgun

import (
	"context"

	"github.com/pkg/errors"
)

// ListMetrics returns account metrics.
//
// NOTE: Only for v1 API. To use the /v1 version define MG_URL in the environment variable
// as `https://api.mailgun.net/v1` or set `v.SetAPIBase("https://api.mailgun.net/v1")`
//
// https://documentation.mailgun.com/docs/mailgun/api-reference/openapi-final/tag/Metrics/
func (mg *MailgunImpl) ListMetrics(ctx context.Context, opts MetricsOptions) (*MetricsResponse, error) {
	payload := newJSONEncodedPayload(opts)
	req := newHTTPRequest(generatePublicApiUrl(mg, metricsEndpoint))
	req.setClient(mg.Client())
	req.setBasicAuth(basicAuthUser, mg.APIKey())

	resp, err := makePostRequest(ctx, req, payload)
	if err != nil {
		return nil, errors.Errorf("POST %s failed: %s", metricsEndpoint, err)
	}

	var ret MetricsResponse
	err = resp.parseFromJSON(&ret)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}

	return &ret, nil
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
