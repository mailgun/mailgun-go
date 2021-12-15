package mailgun

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func (ms *mockServer) addUnsubscribesRoutes(r *mux.Router) {
	r.HandleFunc("/{domain}/unsubscribes/{address}", ms.getUnsubscribe).Methods(http.MethodGet)

	r.HandleFunc("/{domain}/unsubscribes", ms.createUnsubscribe).Methods(http.MethodPost)
	r.HandleFunc("/{domain}/unsubscribes/{address}", ms.deleteUnsubscribe).Methods(http.MethodDelete)

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

func (ms *mockServer) getUnsubscribe(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, unsubscribe := range ms.unsubscribes {
		if unsubscribe.Address == mux.Vars(r)["address"] {
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

	if r.Header.Get("Content-Type") == "application/json" {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("{\"message\": \"Add multiple unsubscribes is not yet implemented\"}"))
		return
	}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	address := r.FormValue("address")
	if len(address) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("{\"message\": \"Invalid format for parameter address: \"}"))
		return
	}

	tag := r.FormValue("tag")

	var addressExist bool
	for i, unsubscribe := range ms.unsubscribes {
		if unsubscribe.Address == address {
			ms.unsubscribes[i].Tags = append(ms.unsubscribes[i].Tags, tag)
			addressExist = true
		}
	}

	if !addressExist {
		unsubscribe := Unsubscribe{
			CreatedAt: RFC2822Time(time.Now()),
			Tags:      []string{tag},
			Address:   address,
		}
		ms.unsubscribes = append(ms.unsubscribes, unsubscribe)
	}

	toJSON(w, map[string]interface{}{
		"message": "Address has been added to the unsubscribes table",
		"address": address,
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
		if unsubscribe.Address == mux.Vars(r)["address"] {
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
			if unsubscribe.Address != mux.Vars(r)["address"] {
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
		if unsubscribe.Address != mux.Vars(r)["address"] {
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
