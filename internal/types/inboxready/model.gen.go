// Package inboxready provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package inboxready

import (
	"time"
)

const (
	BasicAuthScopes = "basicAuth.Scopes"
)

// GithubComMailgunAlertsInternalAPIListAlertsEventsResp defines model for github.com-mailgun-alerts-internal-api-ListAlertsEventsResp.
type GithubComMailgunAlertsInternalAPIListAlertsEventsResp struct {
	Events []string `json:"events"`
}

// GithubComMailgunAlertsInternalAPIListSlackChannelsResp defines model for github.com-mailgun-alerts-internal-api-ListSlackChannelsResp.
type GithubComMailgunAlertsInternalAPIListSlackChannelsResp struct {
	Items []struct {
		ID         string `json:"id"`
		IsArchived bool   `json:"is_archived"`
		Name       string `json:"name"`
	} `json:"items"`
}

// GithubComMailgunAlertsInternalAPIResetWebhooksSigningKeyResp defines model for github.com-mailgun-alerts-internal-api-ResetWebhooksSigningKeyResp.
type GithubComMailgunAlertsInternalAPIResetWebhooksSigningKeyResp struct {
	SigningKey string `json:"signing_key"`
}

// GithubComMailgunAlertsInternalAPISettingsResponse defines model for github.com-mailgun-alerts-internal-api-SettingsResponse.
type GithubComMailgunAlertsInternalAPISettingsResponse struct {
	Events []struct {
		// Channel The delivery channel for the alert.
		Channel string `json:"channel"`

		// DisabledAt Read only. When present
		DisabledAt *time.Time `json:"disabled_at,omitempty"`

		// EventType The event type that is alerted on.
		EventType string `json:"event_type"`

		// ID The unique identifier for the alert settings record.
		ID *[]int32 `json:"id,omitempty"`

		// Settings This object contains channel-specific settings.
		Settings struct {
			// ChannelIds For slack channel
			ChannelIds *[]string `json:"channel_ids,omitempty"`

			// Emails For email channel
			Emails *[]string `json:"emails,omitempty"`

			// URL For webhook channel.
			URL *string `json:"url,omitempty"`
		} `json:"settings"`
	} `json:"events"`
	Slack    GithubComMailgunAlertsInternalSettingsSlack    `json:"slack"`
	Webhooks GithubComMailgunAlertsInternalSettingsWebhooks `json:"webhooks"`
}

// GithubComMailgunAlertsInternalAPISlackChannel defines model for github.com-mailgun-alerts-internal-api-SlackChannel.
type GithubComMailgunAlertsInternalAPISlackChannel struct {
	ID         string `json:"id"`
	IsArchived bool   `json:"is_archived"`
	Name       string `json:"name"`
}

// GithubComMailgunAlertsInternalAPITestEmailRequest defines model for github.com-mailgun-alerts-internal-api-TestEmailRequest.
type GithubComMailgunAlertsInternalAPITestEmailRequest struct {
	Emails    []string `json:"emails"`
	EventType string   `json:"event_type"`
}

// GithubComMailgunAlertsInternalAPITestSlackRequest defines model for github.com-mailgun-alerts-internal-api-TestSlackRequest.
type GithubComMailgunAlertsInternalAPITestSlackRequest struct {
	// ChannelIds If omitted would be taken from the event settings
	ChannelIds *[]string `json:"channel_ids,omitempty"`
	EventType  string    `json:"event_type"`
}

// GithubComMailgunAlertsInternalAPITestWebhookRequest defines model for github.com-mailgun-alerts-internal-api-TestWebhookRequest.
type GithubComMailgunAlertsInternalAPITestWebhookRequest struct {
	EventType string `json:"event_type"`
	URL       string `json:"url"`
}

// GithubComMailgunAlertsInternalAPIMutateEventSettingReq defines model for github.com-mailgun-alerts-internal-api-mutateEventSettingReq.
type GithubComMailgunAlertsInternalAPIMutateEventSettingReq struct {
	// Channel The delivery method for the alert. Supported values include 'webhook'
	Channel string `json:"channel"`

	// EventType The type of event for which you would like to receive alerts.
	EventType string `json:"event_type"`

	// Settings The details pertaining to the specified channel. Please note that the contents of this object differ per channel type.
	Settings struct {
		// ChannelIds For slack channel
		ChannelIds *[]string `json:"channel_ids,omitempty"`

		// Emails For email channel
		Emails *[]string `json:"emails,omitempty"`

		// URL For webhook channel.
		URL *string `json:"url,omitempty"`
	} `json:"settings"`
}

// GithubComMailgunAlertsInternalSettingsEventSettings defines model for github.com-mailgun-alerts-internal-settings-EventSettings.
type GithubComMailgunAlertsInternalSettingsEventSettings struct {
	// Channel The delivery channel for the alert.
	Channel string `json:"channel"`

	// DisabledAt Read only. When present
	DisabledAt *time.Time `json:"disabled_at,omitempty"`

	// EventType The event type that is alerted on.
	EventType string `json:"event_type"`

	// ID The unique identifier for the alert settings record.
	ID *[]int32 `json:"id,omitempty"`

	// Settings This object contains channel-specific settings.
	Settings struct {
		// ChannelIds For slack channel
		ChannelIds *[]string `json:"channel_ids,omitempty"`

		// Emails For email channel
		Emails *[]string `json:"emails,omitempty"`

		// URL For webhook channel.
		URL *string `json:"url,omitempty"`
	} `json:"settings"`
}

// GithubComMailgunAlertsInternalSettingsSlack defines model for github.com-mailgun-alerts-internal-settings-Slack.
type GithubComMailgunAlertsInternalSettingsSlack struct {
	Scope    string `json:"scope"`
	TeamID   string `json:"team_id"`
	TeamName string `json:"team_name"`
	Token    string `json:"token"`
}

// GithubComMailgunAlertsInternalSettingsWebhooks defines model for github.com-mailgun-alerts-internal-settings-Webhooks.
type GithubComMailgunAlertsInternalSettingsWebhooks struct {
	SigningKey string `json:"signing_key"`
}

// GithubComMailgunScaffoldHttpapiGenericResponse defines model for github.com-mailgun-scaffold-httpapi-GenericResponse.
type GithubComMailgunScaffoldHttpapiGenericResponse struct {
	Message string `json:"message"`
}

// POSTV1AlertsEmailTestJSONRequestBody defines body for POSTV1AlertsEmailTest for application/json ContentType.
type POSTV1AlertsEmailTestJSONRequestBody = GithubComMailgunAlertsInternalAPITestEmailRequest

// POSTV1AlertsSettingsEventsJSONRequestBody defines body for POSTV1AlertsSettingsEvents for application/json ContentType.
type POSTV1AlertsSettingsEventsJSONRequestBody = GithubComMailgunAlertsInternalAPIMutateEventSettingReq

// PUTV1AlertsSettingsEventsIDJSONRequestBody defines body for PUTV1AlertsSettingsEventsID for application/json ContentType.
type PUTV1AlertsSettingsEventsIDJSONRequestBody = GithubComMailgunAlertsInternalAPIMutateEventSettingReq

// PUTV1AlertsSettingsSlackJSONRequestBody defines body for PUTV1AlertsSettingsSlack for application/json ContentType.
type PUTV1AlertsSettingsSlackJSONRequestBody = GithubComMailgunAlertsInternalSettingsSlack

// POSTV1AlertsSlackTestJSONRequestBody defines body for POSTV1AlertsSlackTest for application/json ContentType.
type POSTV1AlertsSlackTestJSONRequestBody = GithubComMailgunAlertsInternalAPITestSlackRequest

// POSTV1AlertsWebhooksTestJSONRequestBody defines body for POSTV1AlertsWebhooksTest for application/json ContentType.
type POSTV1AlertsWebhooksTestJSONRequestBody = GithubComMailgunAlertsInternalAPITestWebhookRequest
