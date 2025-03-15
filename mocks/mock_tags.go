package mocks

import (
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addTagsRoutes(r chi.Router) {
	r.Get("/{domain}/tags", ms.listTags)
	r.Get("/{domain}/tags/{tag}", ms.getTags)
	r.Delete("/{domain}/tags/{tag}", ms.deleteTags)
	r.Put("/{domain}/tags/{tag}", ms.createUpdateTags)

	tenMinutesBefore := time.Now().Add(-10 * time.Minute)
	now := time.Now()
	ms.tags = append(ms.tags, mtypes.Tag{
		Value:       "test",
		Description: "test description",
		FirstSeen:   &tenMinutesBefore,
		LastSeen:    &now,
	})

	ms.tags = append(ms.tags, mtypes.Tag{
		Value:       "test2",
		Description: "test2 description",
		FirstSeen:   &tenMinutesBefore,
		LastSeen:    &now,
	})
}

func (ms *Server) listTags(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var idx []string
	for _, t := range ms.unsubscribes {
		idx = append(idx, t.Address)
	}

	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}

	page := r.FormValue("page")
	var pivot string
	if len(page) != 0 {
		pivot = r.FormValue("p")
		if pivot == "" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"message\": \"Invalid parameter: pivot \"}"))
			return
		}
	}
	start, end := pageOffsets(idx, page, pivot, limit)
	var nextAddress, prevAddress string
	var results []mtypes.Tag

	if start != end {
		results = ms.tags[start:end]
		nextAddress = results[len(results)-1].Value
		prevAddress = results[0].Value
	} else {
		results = []mtypes.Tag{}
		nextAddress = pivot
		prevAddress = pivot
	}

	toJSON(w, mtypes.TagsResponse{
		Paging: mtypes.Paging{
			First: getPageURL(r, url.Values{
				"page": []string{"first"},
			}),
			Last: getPageURL(r, url.Values{
				"page": []string{"last"},
			}),
			Next: getPageURL(r, url.Values{
				"page": []string{"next"},
				"p":    []string{nextAddress},
			}),
			Previous: getPageURL(r, url.Values{
				"page": []string{"prev"},
				"p":    []string{prevAddress},
			}),
		},
		Items: results,
	})
}

func (ms *Server) getTags(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	for _, tag := range ms.tags {
		if tag.Value == chi.URLParam(r, "tag") {
			toJSON(w, tag)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{\"message\": \"Tag not found\"}"))
}

func (ms *Server) createUpdateTags(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tag := chi.URLParam(r, "tag")
	description := r.FormValue("description")

	var tagExists bool
	for i, existingTag := range ms.tags {
		if tag == existingTag.Value {
			ms.tags[i].Description = description
			tagExists = true
		}
	}

	if !tagExists {
		ms.tags = append(ms.tags, mtypes.Tag{Value: tag, Description: description})
	}

	toJSON(w, map[string]any{
		"message": "Tag updated",
	})
}

func (ms *Server) deleteTags(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for i, existingTag := range ms.tags {
		if existingTag.Value == chi.URLParam(r, "tag") {
			ms.tags = append(ms.tags[:i], ms.tags[i+1:len(ms.tags)]...)
		}
	}

	toJSON(w, map[string]any{
		"message": "Tag deleted",
	})
}
