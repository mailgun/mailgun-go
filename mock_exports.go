package mailgun

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addExportRoutes(r chi.Router) {
	r.Post("/exports", ms.postExports)
	r.Get("/exports", ms.listExports)
	r.Get("/exports/{id}", ms.getExport)
	r.Get("/exports/{id}/download_url", ms.getExportLink)
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
		if export.ID == chi.URLParam(r, "id") {
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
