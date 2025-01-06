package mailgun

import (
	"net/http"

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

	var results EmailVerification
	results.Address = r.FormValue("address")
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
