package mocks

import (
	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"net/http"
	"strconv"
)

func (ms *Server) addAPIKeysRoutes(r chi.Router) {
	r.Get("/v"+strconv.Itoa(mtypes.APIKeysVersion)+"/"+mtypes.APIKeysEndpoint, ms.listAPIKeys)
	r.Post("/v"+strconv.Itoa(mtypes.APIKeysVersion)+"/"+mtypes.APIKeysEndpoint, ms.createAPIKey)
	r.Delete("/v"+strconv.Itoa(mtypes.APIKeysVersion)+"/"+mtypes.APIKeysEndpoint+"/{key_id}", ms.deleteApiKey)
}

func (ms *Server) listAPIKeys(w http.ResponseWriter, _ *http.Request) {
	resp := mtypes.GetAPIKeyListResponse{
		Items: []mtypes.APIKey{{ID: "1"}, {ID: "2"}},
	}

	toJSON(w, resp)
}

func (ms *Server) createAPIKey(w http.ResponseWriter, _ *http.Request) {
	resp := mtypes.CreateAPIKeyResponse{
		Key: mtypes.APIKey{ID: "1", Role: "basic"},
	}

	toJSON(w, resp)
}

func (ms *Server) deleteApiKey(w http.ResponseWriter, _ *http.Request) {
	resp := mtypes.DeleteAPIKeyResponse{
		Message: "success",
	}

	toJSON(w, resp)
}
