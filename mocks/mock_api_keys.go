package mocks

import (
	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"net/http"
)

func (ms *Server) addAPIKeysRoutes(r chi.Router) {
	r.Get("/v1/"+mtypes.APIKeysEndpoint, ms.listAPIKeys)
}

func (ms *Server) listAPIKeys(w http.ResponseWriter, _ *http.Request) {
	resp := mtypes.APIKeyList{
		Items: []mtypes.APIKey{},
	}

	toJSON(w, resp)
}
