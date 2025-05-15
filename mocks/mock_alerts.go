package mocks

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v5/internal/types/inboxready"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

var alertID = uuid.MustParse("12345678-1234-5678-1234-123456789012")

const uuidRE = `{id:[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12}}`

func (ms *Server) addAlertsRoutes(r chi.Router) {
	r.Get("/v1/"+mtypes.AlertsSettingsEndpoint, ms.listAlerts)
	r.Post("/v1/"+mtypes.AlertsSettingsEndpoint, ms.addAlert)
	r.Delete("/v1/"+mtypes.AlertsSettingsEndpoint+"/"+uuidRE, ms.ok)
}

func (ms *Server) listAlerts(w http.ResponseWriter, _ *http.Request) {
	resp := mtypes.AlertsSettingsResponse{
		Events: []mtypes.AlertsEventSettingResponse{
			{
				Channel:   mtypes.AlertsEmailChannel,
				EventType: "ip_listed",
				ID:        ptr(uuid.New()),
				Settings: mtypes.AlertsChannelSettings{
					Emails: []string{"mail1@example.com", "mail2@example.com"},
				},
			},
			{
				Channel:   mtypes.AlertsWebhookChannel,
				EventType: "ip_delisted",
				ID:        ptr(uuid.New()),
				Settings: mtypes.AlertsChannelSettings{
					URL: ptr("https://example.com/hook"),
				},
			},
		},
		Webhooks: inboxready.GithubComMailgunAlertsInternalSettingsWebhooks{
			SigningKey: "secret-key",
		},
	}

	toJSON(w, resp)
}

func (ms *Server) addAlert(w http.ResponseWriter, r *http.Request) {
	var v mtypes.AlertsEventSettingResponse

	err := json.NewDecoder(r.Body).Decode(&v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	v.ID = ptr(alertID)
	toJSON(w, v)
}
