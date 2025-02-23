package mailgun

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v4/mtypes"
)

type DomainContainer struct {
	Domain              mtypes.Domain            `json:"domain"`
	ReceivingDNSRecords []mtypes.DNSRecord       `json:"receiving_dns_records"`
	SendingDNSRecords   []mtypes.DNSRecord       `json:"sending_dns_records"`
	Connection          *mtypes.DomainConnection `json:"connection,omitempty"`
	Tracking            *mtypes.DomainTracking   `json:"tracking,omitempty"`
	TagLimits           *mtypes.TagLimits        `json:"limits,omitempty"`
}

func (ms *mockServer) addDomainRoutes(r chi.Router) {
	ms.domainList = append(ms.domainList, DomainContainer{
		Domain: mtypes.Domain{
			CreatedAt:    mtypes.RFC2822Time(time.Now().UTC()),
			Name:         "mailgun.test",
			SMTPLogin:    "postmaster@mailgun.test",
			SMTPPassword: "4rtqo4p6rrx9",
			Wildcard:     true,
			SpamAction:   mtypes.SpamActionDisabled,
			State:        "active",
			WebScheme:    "http",
		},
		Connection: &mtypes.DomainConnection{
			RequireTLS:       true,
			SkipVerification: true,
		},
		TagLimits: &mtypes.TagLimits{
			Limit: 50000,
			Count: 5000,
		},
		Tracking: &mtypes.DomainTracking{
			Click: mtypes.TrackingStatus{Active: true},
			Open:  mtypes.TrackingStatus{Active: true},
			Unsubscribe: mtypes.TrackingStatus{
				Active:     false,
				HTMLFooter: "\n<br>\n<p><a href=\"%unsubscribe_url%\">unsubscribe</a></p>\n",
				TextFooter: "\n\nTo unsubscribe click: <%unsubscribe_url%>\n\n",
			},
		},
		ReceivingDNSRecords: []mtypes.DNSRecord{
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
		SendingDNSRecords: []mtypes.DNSRecord{
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

	r.Get("/v4/domains", ms.listDomains)
	r.Post("/v4/domains", ms.createDomain)
	r.Get("/v4/domains/{domain}", ms.getDomain)
	r.Put("/v4/domains/{domain}", ms.updateDomain)
	r.Put("/v4/domains/{domain}/verify", ms.getDomain)
	r.Delete("/v3/domains/{domain}", ms.deleteDomain)

	r.Get("/v3/domains/{domain}/connection", ms.getConnection)
	r.Put("/v3/domains/{domain}/connection", ms.updateConnection)

	r.Get("/v3/domains/{domain}/tracking", ms.getTracking)
	r.Put("/v3/domains/{domain}/tracking/click", ms.updateClickTracking)
	r.Put("/v3/domains/{domain}/tracking/open", ms.updateOpenTracking)
	r.Put("/v3/domains/{domain}/tracking/unsubscribe", ms.updateUnsubTracking)
	r.Get("/v3/domains/{domain}/limits/tag", ms.getTagLimits)

	r.Put("/v3/domains/{domain}/dkim_selector", ms.updateDKIMSelector)
}

func (ms *mockServer) listDomains(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var list []mtypes.Domain
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
		toJSON(w, mtypes.ListDomainsResponse{
			TotalCount: len(list),
			Items:      []mtypes.Domain{},
		})
		return
	}

	toJSON(w, mtypes.ListDomainsResponse{
		TotalCount: len(list),
		Items:      list[skip:end],
	})
}

func (ms *mockServer) getDomain(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

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

func (ms *mockServer) createDomain(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	ms.domainList = append(ms.domainList, DomainContainer{
		Domain: mtypes.Domain{
			CreatedAt:    mtypes.RFC2822Time(time.Now()),
			Name:         r.FormValue("name"),
			SMTPLogin:    r.FormValue("smtp_login"),
			SMTPPassword: r.FormValue("smtp_password"),
			Wildcard:     stringToBool(r.FormValue("wildcard")),
			SpamAction:   mtypes.SpamAction(r.FormValue("spam_action")),
			State:        "active",
			WebScheme:    "http",
		},
	})
	toJSON(w, okResp{Message: "Domain has been created"})
}

func (ms *mockServer) updateDomain(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, domain := range ms.domainList {
		if domain.Domain.Name == chi.URLParam(r, "domain") {
			domain.Domain.WebScheme = r.FormValue("web_scheme")
		}
	}

	toJSON(w, okResp{Message: "Domain has been updated"})
}

func (ms *mockServer) deleteDomain(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

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

func (ms *mockServer) getConnection(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			resp := mtypes.DomainConnectionResponse{
				Connection: *d.Connection,
			}
			toJSON(w, resp)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *mockServer) updateConnection(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for i, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			ms.domainList[i].Connection = &mtypes.DomainConnection{
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

func (ms *mockServer) getTracking(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			resp := mtypes.DomainTrackingResponse{
				Tracking: *d.Tracking,
			}
			toJSON(w, resp)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *mockServer) updateClickTracking(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

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

func (ms *mockServer) updateOpenTracking(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

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

func (ms *mockServer) updateUnsubTracking(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

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

func (ms *mockServer) getTagLimits(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

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

func (ms *mockServer) updateDKIMSelector(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			if r.FormValue("dkim_selector") == "" {
				toJSON(w, okResp{Message: "dkim_selector param required"})
				return
			}
			toJSON(w, okResp{Message: "updated dkim selector"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *mockServer) updateWebPrefix(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			if r.FormValue("web_prefix") == "" {
				toJSON(w, okResp{Message: "web_prefix param required"})
				return
			}
			toJSON(w, okResp{Message: "updated web prefix"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}
