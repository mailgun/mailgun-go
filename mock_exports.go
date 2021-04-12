package mailgun

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func (ms *mockServer) addExportRoutes(r *mux.Router) {
	r.HandleFunc("/exports", ms.postExports).Methods(http.MethodPost)
	r.HandleFunc("/exports", ms.listExports).Methods(http.MethodGet)
	r.HandleFunc("/exports/{id}", ms.getExport).Methods(http.MethodGet)
	r.HandleFunc("/exports/{id}/download_url", ms.getExportLink).Methods(http.MethodGet)
}

func (ms *mockServer) postExports(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	e := Export{
		ID:     strconv.Itoa(len(ms.exportList)),
		URL:    r.FormValue("url"),
		Status: "complete",
	}

	ms.exportList = append(ms.exportList, e)
	toJSON(w, okResp{Message: "success"})
}

func (ms *mockServer) listExports(w http.ResponseWriter, _ *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	toJSON(w, ExportList{
		Items: ms.exportList,
	})
}

func (ms *mockServer) getExport(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	for _, export := range ms.exportList {
		if export.ID == mux.Vars(r)["id"] {
			toJSON(w, export)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "export not found"})
}

func (ms *mockServer) getExportLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/some/s3/url")
	w.WriteHeader(http.StatusFound)
}
