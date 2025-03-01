package mocks

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

type routeResponse struct {
	Route mtypes.Route `json:"route"`
}

func (ms *Server) addRoutes(r chi.Router) {
	r.Post("/routes", ms.createRoute)
	r.Get("/routes", ms.listRoutes)
	r.Get("/routes/{id}", ms.getRoute)
	r.Put("/routes/{id}", ms.updateRoute)
	r.Delete("/routes/{id}", ms.deleteRoute)

	for i := 0; i < 10; i++ {
		ms.routeList = append(ms.routeList, mtypes.Route{
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

func (ms *Server) listRoutes(w http.ResponseWriter, r *http.Request) {
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
		toJSON(w, mtypes.RoutesListResponse{
			TotalCount: len(ms.routeList),
			Items:      []mtypes.Route{},
		})
		return
	}

	toJSON(w, mtypes.RoutesListResponse{
		TotalCount: len(ms.routeList),
		Items:      ms.routeList[skip:end],
	})
}

func (ms *Server) getRoute(w http.ResponseWriter, r *http.Request) {
	for _, item := range ms.routeList {
		if item.Id == chi.URLParam(r, "id") {
			toJSON(w, routeResponse{Route: item})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "route not found"})
}

func (ms *Server) createRoute(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	if r.FormValue("action") == "" {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "'action' parameter is required"})
		return
	}

	ms.routeList = append(ms.routeList, mtypes.Route{
		CreatedAt:   mtypes.RFC2822Time(time.Now().UTC()),
		Id:          randomString(10, "ID-"),
		Priority:    stringToInt(r.FormValue("priority")),
		Description: r.FormValue("description"),
		Expression:  r.FormValue("expression"),
		Actions:     r.Form["action"],
	})
	toJSON(w, mtypes.CreateRouteResp{
		Message: "Route has been created",
		Route:   ms.routeList[len(ms.routeList)-1],
	})
}

func (ms *Server) updateRoute(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for i, item := range ms.routeList {
		if item.Id == chi.URLParam(r, "id") {

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

func (ms *Server) deleteRoute(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	result := ms.routeList[:0]
	for _, item := range ms.routeList {
		if item.Id == chi.URLParam(r, "id") {
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
