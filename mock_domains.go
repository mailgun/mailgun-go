package mailgun

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func (ms *MockServer) addDomainRoutes(r chi.Router) {

	ms.domainList = append(ms.domainList, SingleDomainResponse{
		Domain: Domain{
			CreatedAt:    "Wed, 10 Jul 2013 19:26:52 GMT",
			Name:         "samples.mailgun.org",
			SMTPLogin:    "postmaster@samples.mailgun.org",
			SMTPPassword: "4rtqo4p6rrx9",
			Wildcard:     true,
			SpamAction:   "disabled",
			State:        "active",
		},
		Connection: &DomainConnection{
			RequireTLS:       true,
			SkipVerification: true,
		},
		ReceivingDNSRecords: []DNSRecord{
			{
				Priority:   "10",
				RecordType: "MX",
				Valid:      "valid",
				Value:      "mxa.mailgun.org",
			},
			{
				Priority:   "10",
				RecordType: "MX",
				Valid:      "valid",
				Value:      "mxb.mailgun.org",
			},
		},
		SendingDNSRecords: []DNSRecord{
			{
				RecordType: "TXT",
				Valid:      "valid",
				Name:       "domain.com",
				Value:      "v=spf1 include:mailgun.org ~all",
			},
			{
				RecordType: "TXT",
				Valid:      "valid",
				Name:       "domain.com",
				Value:      "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUA....",
			},
			{
				RecordType: "CNAME",
				Valid:      "valid",
				Name:       "email.domain.com",
				Value:      "mailgun.org",
			},
		},
	})

	r.Get("/domains", ms.listDomains)
	r.Post("/domains", ms.createDomain)
	r.Get("/domains/{domain}", ms.getDomain)
	//r.Put("/domains/{domain}/verify", ms.verifyDomain)
	r.Delete("/domains/{domain}", ms.deleteDomain)
	//r.Get("/domains/{domain}/credentials", ms.getCredentials)
	//r.Post("/domains/{domain}/credentials", ms.createCredentials)
	//r.Put("/domains/{domain}/credentials/{login}", ms.updateCredentials)
	//r.Delete("/domains/{domain}/credentials/{login}", ms.deleteCredentials)
	r.Get("/domains/{domain}/connection", ms.getConnection)
	r.Put("/domains/{domain}/connection", ms.updateConnection)
}

func (ms *MockServer) listDomains(w http.ResponseWriter, _ *http.Request) {
	var list []Domain
	for _, domain := range ms.domainList {
		list = append(list, domain.Domain)
	}

	toJSON(w, DomainList{
		TotalCount: len(list),
		Items:      list,
	})
}

func (ms *MockServer) getDomain(w http.ResponseWriter, r *http.Request) {
	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			d.Connection = nil
			toJSON(w, d)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, OK{Message: "domain not found"})
}

func (ms *MockServer) createDomain(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	ms.domainList = append(ms.domainList, SingleDomainResponse{
		Domain: Domain{
			CreatedAt:    formatMailgunTime(&now),
			Name:         r.FormValue("name"),
			SMTPLogin:    r.FormValue("smtp_login"),
			SMTPPassword: r.FormValue("smtp_password"),
			Wildcard:     stringToBool(r.FormValue("wildcard")),
			SpamAction:   r.FormValue("spam_action"),
			State:        "active",
		},
	})
	toJSON(w, OK{Message: "Domain has been created"})
}

func (ms *MockServer) deleteDomain(w http.ResponseWriter, r *http.Request) {
	result := ms.domainList[:0]
	for _, domain := range ms.domainList {
		if domain.Domain.Name == chi.URLParam(r, "domain") {
			continue
		}
		result = append(result, domain)
	}

	if len(result) != len(ms.domainList) {
		toJSON(w, OK{Message: "success"})
		ms.domainList = result
		return
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, OK{Message: "domain not found"})
}

func (ms *MockServer) getConnection(w http.ResponseWriter, r *http.Request) {
	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			resp := DomainConnectionResponse{
				Connection: *d.Connection,
			}
			toJSON(w, resp)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, OK{Message: "domain not found"})
}

func (ms *MockServer) updateConnection(w http.ResponseWriter, r *http.Request) {
	for i, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			ms.domainList[i].Connection = &DomainConnection{
				RequireTLS:       stringToBool(r.FormValue("require_tls")),
				SkipVerification: stringToBool(r.FormValue("skip_verification")),
			}
			toJSON(w, OK{Message: "Domain connection settings have been updated, may take 10 minutes to fully propagate"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, OK{Message: "domain not found"})
}

func stringToBool(b string) bool {
	result, err := strconv.ParseBool(b)
	if err != nil {
		panic(err)
	}
	return result
}
