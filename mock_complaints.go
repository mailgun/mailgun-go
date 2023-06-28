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

func (ms *mockServer) addComplaintsRoutes(r chi.Router) {
	r.Get("/{domain}/complaints", ms.listComplaints)
	r.Get("/{domain}/complaints/{address}", ms.getComplaint)
	r.Delete("/{domain}/complaints/{address}", ms.deleteComplaint)
	r.Post("/{domain}/complaints", ms.createComplaint)

	ms.complaints = append(ms.complaints, Complaint{
		CreatedAt: RFC2822Time(time.Now()),
		Address:   "foo@mailgun.test",
	})

	ms.complaints = append(ms.complaints, Complaint{
		CreatedAt: RFC2822Time(time.Now()),
		Address:   "alice@example.com",
	})
}

func (ms *mockServer) listComplaints(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var idx []string
	for _, t := range ms.complaints {
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
	var results []Complaint

	if start != end {
		results = ms.complaints[start:end]
		nextAddress = results[len(results)-1].Address
		prevAddress = results[0].Address
	} else {
		results = []Complaint{}
		nextAddress = pivot
		prevAddress = pivot
	}

	toJSON(w, complaintsResponse{
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

func (ms *mockServer) getComplaint(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, complaint := range ms.complaints {
		if complaint.Address == chi.URLParam(r, "address") {
			toJSON(w, complaint)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{\"message\": \"Address not found in complaints table\"}"))
}

func (ms *mockServer) createComplaint(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var complaints []Complaint
	if r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"message\": \"Can't read request body\"}"))
			return
		}

		err = json.Unmarshal(body, &complaints)
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

		address := r.FormValue("address")
		if len(address) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"message\": \"Invalid format for parameter address: \"}"))
			return
		}

		complaints = append(complaints, Complaint{Address: address, CreatedAt: RFC2822Time(time.Now())})
	}

	for _, complaint := range complaints {
		var addressExist bool
		for _, existingComplaint := range ms.complaints {
			if existingComplaint.Address == complaint.Address {
				addressExist = true
			}
		}

		if !addressExist {
			complaint.CreatedAt = RFC2822Time(time.Now())
			ms.complaints = append(ms.complaints, complaint)
		}
	}

	toJSON(w, map[string]interface{}{
		"message": "Address has been added to the complaints table",
		"address": fmt.Sprint(complaints),
	})
}

func (ms *mockServer) deleteComplaint(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for i, complaint := range ms.complaints {
		if complaint.Address == chi.URLParam(r, "address") {
			ms.complaints = append(ms.complaints[:i], ms.complaints[i+1:len(ms.complaints)]...)

			toJSON(w, map[string]interface{}{
				"message": "Complaint has been removed",
			})
			return
		}
	}

	toJSON(w, map[string]interface{}{
		"message": "Address not found in complaints table",
	})
	return
}
