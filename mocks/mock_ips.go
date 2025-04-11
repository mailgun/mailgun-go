package mocks

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addIPRoutes(r chi.Router) {
	r.Get("/ips", ms.listIPS)
	r.Get("/ips/{ip}", ms.getIPAddress)
	r.Route("/domains/{domain}/ips", func(r chi.Router) {
		r.Get("/", ms.listDomainIPS)
		r.Get("/{ip}", ms.getIPAddress)
		r.Post("/", ms.postDomainIPS)
		r.Delete("/{ip}", ms.deleteDomainIPS)
	})
	r.Get("/ips/{ip}/domains", ms.listIPDomains)
}

func (ms *Server) listIPS(w http.ResponseWriter, _ *http.Request) {
	toJSON(w, mtypes.IPAddressListResponse{
		TotalCount: 2,
		Items:      []string{"172.0.0.1", "192.168.1.1"},
	})
}

func (ms *Server) getIPAddress(w http.ResponseWriter, r *http.Request) {
	toJSON(w, mtypes.IPAddress{
		IP:        chi.URLParam(r, "ip"),
		RDNS:      "luna.mailgun.net",
		Dedicated: true,
	})
}

func (ms *Server) listDomainIPS(w http.ResponseWriter, _ *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	toJSON(w, mtypes.IPAddressListResponse{
		TotalCount: 2,
		Items:      ms.domainIPS,
	})
}

func (ms *Server) postDomainIPS(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	ms.domainIPS = append(ms.domainIPS, r.FormValue("ip"))
	toJSON(w, okResp{Message: "success"})
}

func (ms *Server) deleteDomainIPS(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	result := ms.domainIPS[:0]
	for _, ip := range ms.domainIPS {
		if ip == chi.URLParam(r, "ip") {
			continue
		}
		result = append(result, ip)
	}

	if len(result) != len(ms.domainIPS) {
		toJSON(w, okResp{Message: "success"})
		ms.domainIPS = result
		return
	}

	// Not the actual error returned by the mailgun API
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "ip not found"})
}

func (ms *Server) listIPDomains(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var list []mtypes.DomainIPs
	for _, domain := range ms.domainList {
		list = append(list, mtypes.DomainIPs{
			Domain: domain.Domain.Name,
			IPs:    ms.domainIPS,
		})
	}

	skip := stringToInt(r.FormValue("skip"))
	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}

	if skip > len(list) {
		skip = len(list)
	}

	end := limit + skip
	if end > len(list) {
		end = len(list)
	}

	// If we are at the end of the list
	if skip == end {
		toJSON(w, mtypes.ListIPDomainsResponse{
			TotalCount: len(list),
			Items:      []mtypes.DomainIPs{},
		})
		return
	}

	toJSON(w, mtypes.ListIPDomainsResponse{
		TotalCount: len(list),
		Items:      list[skip:end],
	})
}
