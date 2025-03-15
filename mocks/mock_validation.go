package mocks

import (
	"net/http"
	"net/mail"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addValidationRoutes(r chi.Router) {
	r.Get("/v4/address/validate", ms.validateEmailV4)
}

func (ms *Server) validateEmailV4(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("address") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "'address' parameter is required"})
		return
	}

	var results mtypes.ValidateEmailResponse
	results.Risk = "unknown"
	_, err := mail.ParseAddress(r.FormValue("address"))
	if err == nil {
		results.Risk = "low"
	}
	results.Reason = []string{"no-reason"}
	results.Result = "deliverable"
	results.Engagement = &mtypes.EngagementData{
		Engaging: false,
		Behavior: "disengaged",
		IsBot:    false,
	}
	toJSON(w, results)
}
