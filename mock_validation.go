package mailgun

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addValidationRoutes(r chi.Router) {
	r.Get("/v4/address/validate", ms.validateEmailV4)
}

func (ms *mockServer) validateEmailV4(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("address") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "'address' parameter is required"})
		return
	}

	var results v4EmailValidationResp
	parts, err := mail.ParseAddress(r.FormValue("address"))
	if err == nil {
		results.IsValid = true
		results.Parts.Domain = strings.Split(parts.Address, "@")[1]
		results.Parts.LocalPart = strings.Split(parts.Address, "@")[0]
		results.Parts.DisplayName = parts.Name
	}
	results.Reason = []string{"no-reason"}
	results.Risk = "unknown"
	results.Result = "deliverable"
	results.Engagement = &EngagementData{
		Engaging: false,
		Behavior: "disengaged",
		IsBot:    false,
	}
	toJSON(w, results)
}
