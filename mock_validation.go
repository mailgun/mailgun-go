package mailgun

import (
	"net/http"
	"net/mail"
	"strings"

	"github.com/gorilla/mux"
)

func (ms *MockServer) addValidationRoutes(r *mux.Router) {
	r.HandleFunc("/v3/address/validate", ms.validateEmail).Methods(http.MethodGet)
	r.HandleFunc("/v3/address/parse", ms.parseEmail).Methods(http.MethodGet)
	r.HandleFunc("/v3/address/private/validate", ms.validateEmail).Methods(http.MethodGet)
	r.HandleFunc("/v3/address/private/parse", ms.parseEmail).Methods(http.MethodGet)
	r.HandleFunc("/v4/address/validate", ms.validateEmailV4).Methods(http.MethodGet)
}

func (ms *MockServer) validateEmailV4(w http.ResponseWriter, r *http.Request) {
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
	toJSON(w, results)
}

func (ms *MockServer) validateEmail(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("address") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "'address' parameter is required"})
		return
	}

	var results EmailVerification
	parts, err := mail.ParseAddress(r.FormValue("address"))
	if err == nil {
		results.IsValid = true
		results.Parts.Domain = strings.Split(parts.Address, "@")[1]
		results.Parts.LocalPart = strings.Split(parts.Address, "@")[0]
		results.Parts.DisplayName = parts.Name
	}
	results.Reason = "no-reason"
	results.Risk = "unknown"
	toJSON(w, results)
}

func (ms *MockServer) parseEmail(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("addresses") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "'addresses' parameter is required"})
		return
	}

	addresses := strings.Split(r.FormValue("addresses"), ",")

	var results addressParseResult
	for _, address := range addresses {
		_, err := mail.ParseAddress(address)
		if err != nil {
			results.Unparseable = append(results.Unparseable, address)
		} else {
			results.Parsed = append(results.Parsed, address)
		}
	}
	toJSON(w, results)
}
