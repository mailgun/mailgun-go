package mocks

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addInboxPlacementRoutes(r chi.Router) {
	r.Post(fmt.Sprintf("/v%d/%s", mtypes.InboxPlacementVersion, mtypes.InboxPlacementTestsEndpoint), ms.createInboxPlacementTest)
}

func (ms *Server) createInboxPlacementTest(w http.ResponseWriter, _ *http.Request) {
	resp := mtypes.CreateInboxPlacementTestResponse{
		Links: mtypes.CreateInboxPlacementTestResponseLinks{
			Results: "https://api.mailgun.net/v4/inbox/results/result-id",
		},
		MailingList: "ibp-123@mailgun.net,seed@domain.tld",
		ResultID:    "result-id",
	}

	toJSON(w, resp)
}
