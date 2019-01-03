package mock

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/mailgun/mailgun-go/schema"
)

var domainIPS []string

func addIPRoutes(r chi.Router) {
	r.Get("/ips", listIPS)
	r.Get("/ips/{ip}", getIPAddress)
	r.Route("/domains/{domain}/ips", func(r chi.Router) {
		r.Get("/", listDomainIPS)
		r.Get("/{ip}", getIPAddress)
		r.Post("/", postDomainIPS)
		r.Delete("/{ip}", deleteDomainIPS)
	})
}

func listIPS(w http.ResponseWriter, _ *http.Request) {
	toJSON(w, schema.IPAddressList{
		TotalCount: 2,
		Items:      []string{"172.0.0.1", "192.168.1.1"},
	})
}

func getIPAddress(w http.ResponseWriter, r *http.Request) {
	toJSON(w, schema.IPAddress{
		IP:        chi.URLParam(r, "ip"),
		RDNS:      "luna.mailgun.net",
		Dedicated: true,
	})
}

func listDomainIPS(w http.ResponseWriter, _ *http.Request) {
	toJSON(w, schema.IPAddressList{
		TotalCount: 2,
		Items:      domainIPS,
	})
}

func postDomainIPS(w http.ResponseWriter, r *http.Request) {
	domainIPS = append(domainIPS, r.FormValue("ip"))
	toJSON(w, schema.OK{Message: "success"})
}

func deleteDomainIPS(w http.ResponseWriter, r *http.Request) {
	result := domainIPS[:0]
	for _, ip := range domainIPS {
		if ip == chi.URLParam(r, "ip") {
			continue
		}
		result = append(result, ip)
	}

	if len(result) != len(domainIPS) {
		toJSON(w, schema.OK{Message: "success"})
		domainIPS = result
		return
	}

	// Not the actual error returned by the mailgun API
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, schema.OK{Message: "ip not deleted"})
}
