package mtypes

import (
	"github.com/mailgun/mailgun-go/v5/internal/types/inboxready"
)

const (
	AlertsEndpoint         = "alerts"
	AlertsSettingsEndpoint = AlertsEndpoint + "/settings"
	AlertsVersion          = 1
)

// channels
const (
	EmailChannel   = "email"
	WebhookChannel = "webhook"
	SlackChannel   = "slack"
)

type AlertsEventsResponse = inboxready.GithubComMailgunAlertsInternalAPIListAlertsEventsResp

type AlertsEventSettingRequest = inboxready.POSTV1AlertsSettingsEventsJSONRequestBody
type AlertsSettingsResponse = inboxready.GithubComMailgunAlertsInternalAPISettingsResponse
type AlertsEventSettingResponse = inboxready.GithubComMailgunAlertsInternalSettingsEventSettings
type AlertsChannelSettings = inboxready.GithubComMailgunAlertsInternalSettingsChannelSettings
