package mailgun

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addUnsubscribesRoutes(r chi.Router) {
	r.Get("/{domain}/unsubscribes", ms.listUnsubscribes)
	r.Get("/{domain}/unsubscribes/{address}", ms.getUnsubscribe)
	r.Delete("/{domain}/unsubscribes/{address}", ms.deleteUnsubscribe)
	r.Post("/{domain}/unsubscribes", ms.createUnsubscribe)

	ms.unsubscribes = append(ms.unsubscribes, Unsubscribe{
		CreatedAt: RFC2822Time(time.Now()),
		Tags:      []string{"*"},
		ID:        "1",
		Address:   "foo@mailgun.test",
	})

	ms.unsubscribes = append(ms.unsubscribes, Unsubscribe{
		CreatedAt: RFC2822Time(time.Now()),
		Tags:      []string{"some", "tag"},
		ID:        "2",
		Address:   "alice@example.com",
	})
}

func (ms *mockServer) listUnsubscribes(w http.ResponseWriter, r *http.Request) {
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
	var results []Unsubscribe

	if start != end {
		results = ms.unsubscribes[start:end]
		nextAddress = results[len(results)-1].Address
		prevAddress = results[0].Address
	} else {
		results = []Unsubscribe{}
		nextAddress = pivot
		prevAddress = pivot
	}

	toJSON(w, unsubscribesResponse{
		Paging: Paging{
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

func (ms *mockServer) getUnsubscribe(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, unsubscribe := range ms.unsubscribes {
		if unsubscribe.Address == chi.URLParam(r, "address") {
			toJSON(w, unsubscribe)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{\"message\": \"Address not found in unsubscribers table\"}"))
}

func (ms *mockServer) createUnsubscribe(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var unsubscribes []Unsubscribe
	if r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"message\": \"Can't read request body\"}"))
			return
		}

		err = json.Unmarshal(body, &unsubscribes)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(fmt.Sprintf("{\"message\": \"Invalid json: %s\"}", err.Error())))
			return
		}
	} else {
		if err := r.ParseForm(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		tag := r.FormValue("tag")

		address := r.FormValue("address")
		if len(address) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"message\": \"Invalid format for parameter address: \"}"))
			return
		}

		unsubscribes = append(unsubscribes, Unsubscribe{Address: address, Tags: []string{tag}})
	}

	for _, unsubscribe := range unsubscribes {
		var addressExist bool
		for i, existingUnsubscribe := range ms.unsubscribes {
			if existingUnsubscribe.Address == unsubscribe.Address {
				ms.unsubscribes[i].Tags = append(ms.unsubscribes[i].Tags, unsubscribe.Tags...)
				addressExist = true
			}
		}

		if !addressExist {
			ms.unsubscribes = append(ms.unsubscribes, unsubscribe)
		}
	}

	toJSON(w, map[string]interface{}{
		"message": "Address has been added to the unsubscribes table",
		"address": fmt.Sprint(unsubscribes),
	})
}

func (ms *mockServer) deleteUnsubscribe(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var addressExist bool
	for _, unsubscribe := range ms.unsubscribes {
		if unsubscribe.Address == chi.URLParam(r, "address") {
			addressExist = true
		}
	}

	if !addressExist {
		toJSON(w, map[string]interface{}{
			"message": "Address not found in unsubscribers table",
		})
		return
	}

	tag := r.FormValue("tag")
	if len(tag) == 0 {
		for i, unsubscribe := range ms.unsubscribes {
			if unsubscribe.Address != chi.URLParam(r, "address") {
				continue
			}
			ms.unsubscribes = append(ms.unsubscribes[:i], ms.unsubscribes[i+1:len(ms.unsubscribes)]...)

			toJSON(w, map[string]interface{}{
				"message": "Unsubscribe event has been removed",
			})
			return
		}
	}

	for i, unsubscribe := range ms.unsubscribes {
		if unsubscribe.Address != chi.URLParam(r, "address") {
			continue
		}
		for j, t := range ms.unsubscribes[i].Tags {
			if t != tag {
				continue
			}
			ms.unsubscribes[i].Tags = append(ms.unsubscribes[i].Tags[:j], ms.unsubscribes[i].Tags[j+1:]...)
			toJSON(w, map[string]interface{}{
				"message": "Unsubscribe event has been removed",
			})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, map[string]string{"message": "Unsubscribe event for this tag does not exist"})
}
