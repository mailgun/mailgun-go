package mailgun

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (ms *MockServer) addIPRoutes(r *mux.Router) {
	r.HandleFunc("/ips", ms.listIPS).Methods(http.MethodGet)
	r.HandleFunc("/ips/{ip}", ms.getIPAddress).Methods(http.MethodGet)
	func(r *mux.Router) {
		r.HandleFunc("", ms.listDomainIPS).Methods(http.MethodGet)
		r.HandleFunc("/{ip}", ms.getIPAddress).Methods(http.MethodGet)
		r.HandleFunc("", ms.postDomainIPS).Methods(http.MethodPost)
		r.HandleFunc("/{ip}", ms.deleteDomainIPS).Methods(http.MethodDelete)
	}(r.PathPrefix("/domains/{domain}/ips").Subrouter())
}

func (ms *MockServer) listIPS(w http.ResponseWriter, _ *http.Request) {
	toJSON(w, ipAddressListResponse{
		TotalCount: 2,
		Items:      []string{"172.0.0.1", "192.168.1.1"},
	})
}

func (ms *MockServer) getIPAddress(w http.ResponseWriter, r *http.Request) {
	toJSON(w, IPAddress{
		IP:        mux.Vars(r)["ip"],
		RDNS:      "luna.mailgun.net",
		Dedicated: true,
	})
}

func (ms *MockServer) listDomainIPS(w http.ResponseWriter, _ *http.Request) {
	toJSON(w, ipAddressListResponse{
		TotalCount: 2,
		Items:      ms.domainIPS,
	})
}

func (ms *MockServer) postDomainIPS(w http.ResponseWriter, r *http.Request) {
	ms.domainIPS = append(ms.domainIPS, r.FormValue("ip"))
	toJSON(w, okResp{Message: "success"})
}

func (ms *MockServer) deleteDomainIPS(w http.ResponseWriter, r *http.Request) {
	result := ms.domainIPS[:0]
	for _, ip := range ms.domainIPS {
		if ip == mux.Vars(r)["ip"] {
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
