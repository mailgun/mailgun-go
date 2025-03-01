package mocks

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addExportRoutes(r chi.Router) {
	r.Post("/exports", ms.postExports)
	r.Get("/exports", ms.listExports)
	r.Get("/exports/{id}", ms.getExport)
	r.Get("/exports/{id}/download_url", ms.getExportLink)
}

func (ms *Server) postExports(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	e := mtypes.Export{
		ID:     strconv.Itoa(len(ms.exportList)),
		URL:    r.FormValue("url"),
		Status: "complete",
	}

	ms.exportList = append(ms.exportList, e)
	toJSON(w, okResp{Message: "success"})
}

func (ms *Server) listExports(w http.ResponseWriter, _ *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	toJSON(w, mtypes.ExportList{
		Items: ms.exportList,
	})
}

func (ms *Server) getExport(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) getExportLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/some/s3/url")
	w.WriteHeader(http.StatusFound)
}
