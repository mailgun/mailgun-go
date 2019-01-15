package mailgun

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

func (ms *MockServer) addExportRoutes(r chi.Router) {
	r.Post("/exports", ms.postExports)
	r.Get("/exports", ms.listExports)
	r.Get("/exports/{id}", ms.getExport)
	r.Get("/exports/{id}/download_url", ms.getExportLink)
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
		if export.ID == chi.URLParam(r, "id") {
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
