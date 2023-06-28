package mailgun

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addIPRoutes(r chi.Router) {
	r.Get("/ips", ms.listIPS)
	r.Get("/ips/{ip}", ms.getIPAddress)
	r.Route("/domains/{domain}/ips", func(r chi.Router) {
		r.Get("/", ms.listDomainIPS)
		r.Get("/{ip}", ms.getIPAddress)
		r.Post("/", ms.postDomainIPS)
		r.Delete("/{ip}", ms.deleteDomainIPS)
	})
}

func (ms *mockServer) listIPS(w http.ResponseWriter, _ *http.Request) {
	toJSON(w, ipAddressListResponse{
		TotalCount: 2,
		Items:      []string{"172.0.0.1", "192.168.1.1"},
	})
}

func (ms *mockServer) getIPAddress(w http.ResponseWriter, r *http.Request) {
	toJSON(w, IPAddress{
		IP:        chi.URLParam(r, "ip"),
		RDNS:      "luna.mailgun.net",
		Dedicated: true,
	})
}

func (ms *mockServer) listDomainIPS(w http.ResponseWriter, _ *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	toJSON(w, ipAddressListResponse{
		TotalCount: 2,
		Items:      ms.domainIPS,
	})
}

func (ms *mockServer) postDomainIPS(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()
	ms.domainIPS = append(ms.domainIPS, r.FormValue("ip"))
	toJSON(w, okResp{Message: "success"})
}

func (ms *mockServer) deleteDomainIPS(w http.ResponseWriter, r *http.Request) {
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
