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

func (ms *mockServer) addBouncesRoutes(r chi.Router) {
	r.Get("/{domain}/bounces", ms.listBounces)
	r.Get("/{domain}/bounces/{address}", ms.getBounce)
	r.Delete("/{domain}/bounces/{address}", ms.deleteBounce)
	r.Delete("/{domain}/bounces", ms.deleteBouncesList)
	r.Post("/{domain}/bounces", ms.createBounce)

	ms.bounces = append(ms.bounces, Bounce{
		CreatedAt: RFC2822Time(time.Now()),
		Error:     "invalid address",
		Code:      "INVALID",
		Address:   "foo@mailgun.test",
	})

	ms.bounces = append(ms.bounces, Bounce{
		CreatedAt: RFC2822Time(time.Now()),
		Error:     "non existing address",
		Code:      "NOT_EXIST",
		Address:   "alice@example.com",
	})
}

func (ms *mockServer) listBounces(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var idx []string
	for _, t := range ms.bounces {
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
	var results []Bounce

	if start != end {
		results = ms.bounces[start:end]
		nextAddress = results[len(results)-1].Address
		prevAddress = results[0].Address
	} else {
		results = []Bounce{}
		nextAddress = pivot
		prevAddress = pivot
	}

	toJSON(w, bouncesListResponse{
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

func (ms *mockServer) getBounce(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, bounce := range ms.bounces {
		if bounce.Address == chi.URLParam(r, "address") {
			toJSON(w, bounce)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("{\"message\": \"Address not found in bounces table\"}"))
}

func (ms *mockServer) createBounce(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var bounces []Bounce
	if r.Header.Get("Content-Type") == "application/json" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{\"message\": \"Can't read request body\"}"))
			return
		}

		err = json.Unmarshal(body, &bounces)
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

		bounceError := r.FormValue("error")
		code := r.FormValue("code")

		address := r.FormValue("address")
		if len(address) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("{\"message\": \"Invalid format for parameter address: \"}"))
			return
		}

		bounces = append(bounces, Bounce{Address: address, Code: code, Error: bounceError})
	}

	for _, bounce := range bounces {
		var addressExist bool
		for _, existingBounce := range ms.bounces {
			if existingBounce.Address == bounce.Address {
				addressExist = true
			}
		}

		if !addressExist {
			ms.bounces = append(ms.bounces, bounce)
		}
	}

	toJSON(w, map[string]interface{}{
		"message": "Address has been added to the bounces table",
		"address": fmt.Sprint(bounces),
	})
}

func (ms *mockServer) deleteBounce(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for i, bounce := range ms.bounces {
		if bounce.Address == chi.URLParam(r, "address") {
			ms.bounces = append(ms.bounces[:i], ms.bounces[i+1:len(ms.bounces)]...)

			toJSON(w, map[string]interface{}{
				"message": "Bounce has been removed",
			})
			return
		}
	}

	toJSON(w, map[string]interface{}{
		"message": "Address not found in bounces table",
	})
}

func (ms *mockServer) deleteBouncesList(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	ms.bounces = []Bounce{}

	toJSON(w, map[string]interface{}{
		"message": "All bounces has been deleted",
	})
}
