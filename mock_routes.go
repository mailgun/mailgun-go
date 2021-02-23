package mailgun

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type routeResponse struct {
	Route Route `json:"route"`
}

func (ms *MockServer) addRoutes(r *mux.Router) {
	r.HandleFunc("/routes", ms.createRoute).Methods(http.MethodPost)
	r.HandleFunc("/routes", ms.listRoutes).Methods(http.MethodGet)
	r.HandleFunc("/routes/{id}", ms.getRoute).Methods(http.MethodGet)
	r.HandleFunc("/routes/{id}", ms.updateRoute).Methods(http.MethodPut)
	r.HandleFunc("/routes/{id}", ms.deleteRoute).Methods(http.MethodDelete)

	for i := 0; i < 10; i++ {
		ms.routeList = append(ms.routeList, Route{
			Id:          randomString(10, "ID-"),
			Priority:    0,
			Description: fmt.Sprintf("Sample Route %d", i),
			Actions: []string{
				`forward("http://myhost.com/messages/")`,
				`stop()`,
			},
			Expression: `match_recipient(".*@samples.mailgun.org")`,
		})
	}
}

func (ms *MockServer) listRoutes(w http.ResponseWriter, r *http.Request) {
	skip := stringToInt(r.FormValue("skip"))
	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}

	if skip > len(ms.routeList) {
		skip = len(ms.routeList)
	}

	end := limit + skip
	if end > len(ms.routeList) {
		end = len(ms.routeList)
	}

	// If we are at the end of the list
	if skip == end {
		toJSON(w, routesListResponse{
			TotalCount: len(ms.routeList),
			Items:      []Route{},
		})
		return
	}

	toJSON(w, routesListResponse{
		TotalCount: len(ms.routeList),
		Items:      ms.routeList[skip:end],
	})
}

func (ms *MockServer) getRoute(w http.ResponseWriter, r *http.Request) {
	for _, item := range ms.routeList {
		if item.Id == mux.Vars(r)["id"] {
			toJSON(w, routeResponse{Route: item})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "route not found"})
}

func (ms *MockServer) createRoute(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("action") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "'action' parameter is required"})
		return
	}

	ms.routeList = append(ms.routeList, Route{
		CreatedAt:   RFC2822Time(time.Now().UTC()),
		Id:          randomString(10, "ID-"),
		Priority:    stringToInt(r.FormValue("priority")),
		Description: r.FormValue("description"),
		Expression:  r.FormValue("expression"),
		Actions:     r.Form["action"],
	})
	toJSON(w, createRouteResp{
		Message: "Route has been created",
		Route:   ms.routeList[len(ms.routeList)-1],
	})
}

func (ms *MockServer) updateRoute(w http.ResponseWriter, r *http.Request) {
	for i, item := range ms.routeList {
		if item.Id == mux.Vars(r)["id"] {

			if r.FormValue("action") != "" {
				ms.routeList[i].Actions = r.Form["action"]
			}
			if r.FormValue("priority") != "" {
				ms.routeList[i].Priority = stringToInt(r.FormValue("priority"))
			}
			if r.FormValue("description") != "" {
				ms.routeList[i].Description = r.FormValue("description")
			}
			if r.FormValue("expression") != "" {
				ms.routeList[i].Expression = r.FormValue("expression")
			}
			toJSON(w, ms.routeList[i])
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "route not found"})
}

func (ms *MockServer) deleteRoute(w http.ResponseWriter, r *http.Request) {
	result := ms.routeList[:0]
	for _, item := range ms.routeList {
		if item.Id == mux.Vars(r)["id"] {
			continue
		}
		result = append(result, item)
	}

	if len(result) != len(ms.domainList) {
		toJSON(w, okResp{Message: "success"})
		ms.routeList = result
		return
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "route not found"})
}
