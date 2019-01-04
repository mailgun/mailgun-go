package mailgun

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
)

type MockServer struct {
	srv *httptest.Server

	domainIPS  []string
	domainList []domainResponse
	exportList []Export
}

func NewMockServer() MockServer {
	ms := MockServer{}

	// Add all our handlers
	r := chi.NewRouter()

	r.Route("/v3", func(r chi.Router) {
		ms.addIPRoutes(r)
		ms.addExportRoutes(r)
		ms.addDomainRoutes(r)
	})

	// Start the server
	ms.srv = httptest.NewServer(r)
	return ms
}

func (ms *MockServer) Stop() {
	ms.srv.Close()
}

func (ms *MockServer) URL() string {
	return ms.srv.URL + "/v3"
}

func toJSON(w http.ResponseWriter, obj interface{}) {
	if err := json.NewEncoder(w).Encode(obj); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
}
