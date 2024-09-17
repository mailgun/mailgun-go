package mailgun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/pkg/errors"
)

// ListMetrics returns account metrics.
// Must be /v1 API.
// https://documentation.mailgun.com/docs/mailgun/api-reference/openapi-final/tag/Metrics/
func (c *Client) ListMetrics(ctx context.Context, opts MetricsOptions) (*MetricsResponse, error) {
	url := fmt.Sprintf("%s", metricsEndpoint)

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(opts)
	if err != nil {
		return nil, errors.Wrap(err, "while marshalling analytics metrics request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, errors.Wrap(err, "while creating analytics metrics request")
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "while making analytics metrics request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.Errorf("POST on '%s' returned %s", url, string(body))
	}

	var ret MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, errors.Wrap(err, "while decoding analytics metrics response")
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
