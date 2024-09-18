package mailgun

import (
	"context"
	"strings"

	"github.com/pkg/errors"
)

// ListMetrics returns domain/account metrics.
//
// NOTE: Only for v1 API. To use the /v1 version define MG_URL in the environment variable
// as `https://api.mailgun.net/v1` or set `v.SetAPIBase("https://api.mailgun.net/v1")`
//
// https://documentation.mailgun.com/docs/mailgun/api-reference/openapi-final/tag/Metrics/
func (mg *MailgunImpl) ListMetrics(ctx context.Context, opts MetricsOptions) (*MetricsResponse, error) {
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
