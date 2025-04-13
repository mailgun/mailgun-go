package mailgun

import (
	"context"

	"github.com/mailgun/mailgun-go/v5/mtypes"
)

const (
	// TODO(vtopc): move to mailgun.go?
	alertsEndpoint = "alerts"
	alertsVersion  = 1
)

type ListAlertEventsOptions struct{}

func (mg *Client) ListAlertsEvents(ctx context.Context, _ *ListAlertEventsOptions,
) (*mtypes.AlertsEventsResponse, error) {
	r := newHTTPRequest(generateApiUrl(mg, alertsVersion, alertsEndpoint))
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.setClient(mg.HTTPClient())

	var resp mtypes.AlertsEventsResponse
	if err := getResponseFromJSON(ctx, r, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
