package mailgun

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (ms *MockServer) addExportRoutes(r *mux.Router) {
	r.HandleFunc("/exports", ms.postExports).Methods(http.MethodPost)
	r.HandleFunc("/exports", ms.listExports).Methods(http.MethodGet)
	r.HandleFunc("/exports/{id}", ms.getExport).Methods(http.MethodGet)
	r.HandleFunc("/exports/{id}/download_url", ms.getExportLink).Methods(http.MethodGet)
}

func (ms *MockServer) postExports(w http.ResponseWriter, r *http.Request) {
	e := Export{
		ID:     strconv.Itoa(len(ms.exportList)),
		URL:    r.FormValue("url"),
		Status: "complete",
	}

	ms.exportList = append(ms.exportList, e)
	toJSON(w, okResp{Message: "success"})
}

func (ms *MockServer) listExports(w http.ResponseWriter, _ *http.Request) {
	toJSON(w, ExportList{
		Items: ms.exportList,
	})
}

func (ms *MockServer) getExport(w http.ResponseWriter, r *http.Request) {
	for _, export := range ms.exportList {
		if export.ID == mux.Vars(r)["id"] {
			toJSON(w, export)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "export not found"})
}

func (ms *MockServer) getExportLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/some/s3/url")
	w.WriteHeader(http.StatusFound)
}
