package mailgun

import (
	"context"

	"github.com/mailgun/mailgun-go/v5/mtypes"
)

const (
	// TODO(vtopc): move to mailgun.go?
	alertsEndpoint         = "alerts"
	alertsSettingsEndpoint = alertsEndpoint + "/settings"
	alertsVersion          = 1
)

type ListAlertsEventsOptions struct{}

// ListAlertsEvents list of events that you can choose to receive alerts for.
func (mg *Client) ListAlertsEvents(ctx context.Context, _ *ListAlertsEventsOptions,
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

type ListAlertsOptions struct{}

// ListAlerts returns a list of all configured alert settings for your account.
func (mg *Client) ListAlerts(ctx context.Context, _ *ListAlertsOptions,
) (*mtypes.AlertsResponse, error) {
	r := newHTTPRequest(generateApiUrl(mg, alertsVersion, alertsSettingsEndpoint))
	r.setBasicAuth(basicAuthUser, mg.APIKey())
	r.setClient(mg.HTTPClient())

	var resp mtypes.AlertsResponse
	if err := getResponseFromJSON(ctx, r, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}
