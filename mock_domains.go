package mailgun

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type domainContainer struct {
	Domain              Domain            `json:"domain"`
	ReceivingDNSRecords []DNSRecord       `json:"receiving_dns_records"`
	SendingDNSRecords   []DNSRecord       `json:"sending_dns_records"`
	Connection          *DomainConnection `json:"connection,omitempty"`
	Tracking            *DomainTracking   `json:"tracking,omitempty"`
	TagLimits           *TagLimits        `json:"limits,omitempty"`
}

func (ms *MockServer) addDomainRoutes(r chi.Router) {

	ms.domainList = append(ms.domainList, domainContainer{
		Domain: Domain{
			CreatedAt:    RFC2822Time(time.Now().UTC()),
			Name:         "mailgun.test",
			SMTPLogin:    "postmaster@mailgun.test",
			SMTPPassword: "4rtqo4p6rrx9",
			Wildcard:     true,
			SpamAction:   SpamActionDisabled,
			State:        "active",
		},
		Connection: &DomainConnection{
			RequireTLS:       true,
			SkipVerification: true,
		},
		TagLimits: &TagLimits{
			Limit: 50000,
			Count: 5000,
		},
		Tracking: &DomainTracking{
			Click: TrackingStatus{Active: true},
			Open:  TrackingStatus{Active: true},
			Unsubscribe: TrackingStatus{
				Active:     false,
				HTMLFooter: "\n<br>\n<p><a href=\"%unsubscribe_url%\">unsubscribe</a></p>\n",
				TextFooter: "\n\nTo unsubscribe click: <%unsubscribe_url%>\n\n",
			},
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
	r.Put("/domains/{domain}/verify", ms.getDomain)
	r.Delete("/domains/{domain}", ms.deleteDomain)
	//r.Get("/domains/{domain}/credentials", ms.getCredentials)
	//r.Post("/domains/{domain}/credentials", ms.createCredentials)
	//r.Put("/domains/{domain}/credentials/{login}", ms.updateCredentials)
	//r.Delete("/domains/{domain}/credentials/{login}", ms.deleteCredentials)
	r.Get("/domains/{domain}/connection", ms.getConnection)
	r.Put("/domains/{domain}/connection", ms.updateConnection)
	r.Get("/domains/{domain}/tracking", ms.getTracking)
	r.Put("/domains/{domain}/tracking/click", ms.updateClickTracking)
	r.Put("/domains/{domain}/tracking/open", ms.updateOpenTracking)
	r.Put("/domains/{domain}/tracking/unsubscribe", ms.updateUnsubTracking)
	r.Get("/domains/{domain}/limits/tag", ms.getTagLimits)
}

func (ms *MockServer) listDomains(w http.ResponseWriter, r *http.Request) {
	var list []Domain
	for _, domain := range ms.domainList {
		list = append(list, domain.Domain)
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
		toJSON(w, domainsListResponse{
			TotalCount: len(list),
			Items:      []Domain{},
		})
		return
	}

	toJSON(w, domainsListResponse{
		TotalCount: len(list),
		Items:      list[skip:end],
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
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *MockServer) createDomain(w http.ResponseWriter, r *http.Request) {
	ms.domainList = append(ms.domainList, domainContainer{
		Domain: Domain{
			CreatedAt:    RFC2822Time(time.Now()),
			Name:         r.FormValue("name"),
			SMTPLogin:    r.FormValue("smtp_login"),
			SMTPPassword: r.FormValue("smtp_password"),
			Wildcard:     stringToBool(r.FormValue("wildcard")),
			SpamAction:   SpamAction(r.FormValue("spam_action")),
			State:        "active",
		},
	})
	toJSON(w, okResp{Message: "Domain has been created"})
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
		toJSON(w, okResp{Message: "success"})
		ms.domainList = result
		return
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *MockServer) getConnection(w http.ResponseWriter, r *http.Request) {
	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			resp := domainConnectionResponse{
				Connection: *d.Connection,
			}
			toJSON(w, resp)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *MockServer) updateConnection(w http.ResponseWriter, r *http.Request) {
	for i, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			ms.domainList[i].Connection = &DomainConnection{
				RequireTLS:       stringToBool(r.FormValue("require_tls")),
				SkipVerification: stringToBool(r.FormValue("skip_verification")),
			}
			toJSON(w, okResp{Message: "Domain connection settings have been updated, may take 10 minutes to fully propagate"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *MockServer) getTracking(w http.ResponseWriter, r *http.Request) {
	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			resp := domainTrackingResponse{
				Tracking: *d.Tracking,
			}
			toJSON(w, resp)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *MockServer) updateClickTracking(w http.ResponseWriter, r *http.Request) {
	for i, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			ms.domainList[i].Tracking.Click.Active = stringToBool(r.FormValue("active"))
			toJSON(w, okResp{Message: "Domain tracking settings have been updated"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *MockServer) updateOpenTracking(w http.ResponseWriter, r *http.Request) {
	for i, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			ms.domainList[i].Tracking.Open.Active = stringToBool(r.FormValue("active"))
			toJSON(w, okResp{Message: "Domain tracking settings have been updated"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *MockServer) updateUnsubTracking(w http.ResponseWriter, r *http.Request) {
	for i, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			ms.domainList[i].Tracking.Unsubscribe.Active = stringToBool(r.FormValue("active"))
			if len(r.FormValue("html_footer")) != 0 {
				ms.domainList[i].Tracking.Unsubscribe.HTMLFooter = r.FormValue("html_footer")
			}
			if len(r.FormValue("text_footer")) != 0 {
				ms.domainList[i].Tracking.Unsubscribe.TextFooter = r.FormValue("text_footer")
			}
			toJSON(w, okResp{Message: "Domain tracking settings have been updated"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *MockServer) getTagLimits(w http.ResponseWriter, r *http.Request) {
	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			if d.TagLimits == nil {
				w.WriteHeader(http.StatusNotFound)
				toJSON(w, okResp{Message: "no limits defined for domain"})
				return
			}
			toJSON(w, d.TagLimits)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}
