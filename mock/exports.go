package mock

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/mailgun/mailgun-go/schema"
)

var exportList []schema.Export

func addExportRoutes(r chi.Router) {
	r.Post("/exports", postExports)
	r.Get("/exports", listExports)
	r.Get("/exports/{id}", getExport)
	r.Get("/exports/{id}/download_url", getExportLink)
}

func postExports(w http.ResponseWriter, r *http.Request) {
	e := schema.Export{
		ID:     strconv.Itoa(len(exportList)),
		URL:    r.FormValue("url"),
		Status: "complete",
	}

	exportList = append(exportList, e)
	toJSON(w, schema.OK{Message: "success"})
}

func listExports(w http.ResponseWriter, _ *http.Request) {
	toJSON(w, schema.ExportList{
		Items: exportList,
	})
}

func getExport(w http.ResponseWriter, r *http.Request) {
	for _, export := range exportList {
		if export.ID == chi.URLParam(r, "id") {
			toJSON(w, export)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, schema.OK{Message: "export not found"})
}

func getExportLink(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/some/s3/url")
	w.WriteHeader(http.StatusFound)
}
