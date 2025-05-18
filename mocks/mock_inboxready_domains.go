package mocks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/internal/types/inboxready"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addInboxreadyDomainsRoutes(r chi.Router) {
	r.Post(fmt.Sprintf("/v%d/%s", mtypes.InboxreadyDomainsVersion, mtypes.InboxreadyDomainsEndpoint), ms.addDomainToMonitoring)
	// TODO(vtopc): add other routes
}

func (ms *Server) addDomainToMonitoring(w http.ResponseWriter, r *http.Request) {
	resp := mtypes.AddDomainToMonitoringResponse{
		Domain: inboxready.InboxReadyGithubComMailgunInboxreadyModelDomain{
			CreatedAt: time.Now().Unix(),
			Name:      r.FormValue("domain"),
			Services: map[string]bool{
				"service1": true,
				"service2": true,
			},
			TxtRecord: "",
			Verified:  inboxready.InboxReadyGithubComMailgunInboxreadyModelVerified{},
		},
		Message: "domain added",
	}

	toJSON(w, resp)
}
