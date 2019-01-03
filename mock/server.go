package mock

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
)

type MailgunServer struct {
	srv *httptest.Server

	domainIPS []string
}

func NewServer() MailgunServer {
	ms := MailgunServer{}

	// Add all our handlers
	r := chi.NewRouter()

	r.Route("/v3", func(r chi.Router) {
		addIPRoutes(r)
	})

	// Start the server
	ms.srv = httptest.NewServer(r)
	return ms
}

func (ms *MailgunServer) Stop() {
	ms.srv.Close()
}

func (ms *MailgunServer) URL() string {
	return ms.srv.URL + "/v3"
}

func toJSON(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}
