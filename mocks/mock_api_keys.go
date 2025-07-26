package mocks

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addAPIKeysRoutes(r chi.Router) {
	apiVersion := strconv.Itoa(mtypes.APIKeysVersion)
	r.Get("/v"+apiVersion+"/"+mtypes.APIKeysEndpoint, ms.listAPIKeys)
	r.Post("/v"+apiVersion+"/"+mtypes.APIKeysEndpoint, ms.createAPIKey)
	r.Delete("/v"+apiVersion+"/"+mtypes.APIKeysEndpoint+"/{key_id}", ms.deleteApiKey)
	r.Post("/v"+apiVersion+"/"+mtypes.APIKeysRegenerateEndpoint, ms.regeneratePublicAPIKey)
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

func (ms *Server) regeneratePublicAPIKey(w http.ResponseWriter, _ *http.Request) {
	resp := mtypes.RegeneratePublicAPIKeyResponse{
		Key:     "public-1",
		Message: "success",
	}

	toJSON(w, resp)
}
