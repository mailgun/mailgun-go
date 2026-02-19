package mocks

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

type DomainContainer struct {
	Domain              mtypes.Domain            `json:"domain"`
	ReceivingDNSRecords []mtypes.DNSRecord       `json:"receiving_dns_records"`
	SendingDNSRecords   []mtypes.DNSRecord       `json:"sending_dns_records"`
	Connection          *mtypes.DomainConnection `json:"connection,omitempty"`
	Tracking            *mtypes.DomainTracking   `json:"tracking,omitempty"`
	TagLimits           *mtypes.TagLimits        `json:"limits,omitempty"`
}

func (ms *Server) addDomainRoutes(r chi.Router) {
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

	r.Put("/v4/domains/{domain}/keys/{dkim_selector}/activate", ms.activateDomainKey)
	r.Get("/v4/domains/{domain}/keys", ms.listDomainKeys)
	r.Put("/v4/domains/{domain}/keys/{dkim_selector}/deactivate", ms.deactivateDomainKey)
	r.Put("/v3/domains/{domain}/dkim_authority", ms.updateDomainDkimAuthority)
	r.Put("/v3/domains/{domain}/dkim_selector", ms.updateDomainDkimSelector)
}

func (ms *Server) listDomains(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) getDomain(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) createDomain(w http.ResponseWriter, r *http.Request) {
	const expectedContentType = multipartFormDataContentType
	if !strings.HasPrefix(r.Header.Get(contentTypeHeader), expectedContentType) {
		// NOTE: not an actual Mailgun API response, just for unit tests,
		//  see https://github.com/mailgun/mailgun-go/pull/470 for more details.
		w.WriteHeader(599)
		toJSON(w, okResp{Message: "Content-Type must be " + expectedContentType})
		return
	}

	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	webScheme := r.FormValue("web_scheme")
	if webScheme == "" {
		webScheme = "http"
	}

	domain := mtypes.Domain{
		CreatedAt:                  mtypes.RFC2822Time(time.Now()),
		Name:                       r.FormValue("name"),
		SMTPLogin:                  r.FormValue("smtp_login"),
		SMTPPassword:               r.FormValue("smtp_password"),
		Wildcard:                   stringToBool(r.FormValue("wildcard")),
		SpamAction:                 mtypes.SpamAction(r.FormValue("spam_action")),
		State:                      "active",
		WebScheme:                  webScheme,
		WebPrefix:                  r.FormValue("web_prefix"),
		RequireTLS:                 stringToBool(r.FormValue("require_tls")),
		SkipVerification:           stringToBool(r.FormValue("skip_verification")),
		UseAutomaticSenderSecurity: stringToBool(r.FormValue("use_automatic_sender_security")),
		ArchiveTo:                  r.FormValue("archive_to"),
		DKIMHost:                   r.FormValue("dkim_host_name"),
		EncryptIncomingMessage:     stringToBool(r.FormValue("encrypt_incoming_message")),
		MailFromHost:               r.FormValue("mailfrom_host"),
	}

	if messageTTL := r.FormValue("message_ttl"); messageTTL != "" {
		domain.MessageTTL = stringToInt(messageTTL)
	}

	ms.domainList = append(ms.domainList, DomainContainer{
		Domain: domain,
	})
	toJSON(w, okResp{Message: "Domain has been created"})
}

func (ms *Server) updateDomain(w http.ResponseWriter, r *http.Request) {
	const expectedContentType = multipartFormDataContentType
	if !strings.HasPrefix(r.Header.Get(contentTypeHeader), expectedContentType) {
		// NOTE: not an actual Mailgun API response, just for unit tests,
		//  see https://github.com/mailgun/mailgun-go/pull/470 for more details.
		w.WriteHeader(599)
		toJSON(w, okResp{Message: "Content-Type must be " + expectedContentType})
		return
	}

	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for i, domain := range ms.domainList {
		if domain.Domain.Name == chi.URLParam(r, "domain") {
			if webScheme := r.FormValue("web_scheme"); webScheme != "" {
				ms.domainList[i].Domain.WebScheme = webScheme
			}
			if webPrefix := r.FormValue("web_prefix"); webPrefix != "" {
				ms.domainList[i].Domain.WebPrefix = webPrefix
			}
			if requireTLS := r.FormValue("require_tls"); requireTLS != "" {
				ms.domainList[i].Domain.RequireTLS = stringToBool(requireTLS)
			}
			if skipVerification := r.FormValue("skip_verification"); skipVerification != "" {
				ms.domainList[i].Domain.SkipVerification = stringToBool(skipVerification)
			}
			if useAutoSecurity := r.FormValue("use_automatic_sender_security"); useAutoSecurity != "" {
				ms.domainList[i].Domain.UseAutomaticSenderSecurity = stringToBool(useAutoSecurity)
			}
			if archiveTo := r.FormValue("archive_to"); archiveTo != "" {
				ms.domainList[i].Domain.ArchiveTo = archiveTo
			}
			if mailFromHost := r.FormValue("mailfrom_host"); mailFromHost != "" {
				ms.domainList[i].Domain.MailFromHost = mailFromHost
			}
			if messageTTL := r.FormValue("message_ttl"); messageTTL != "" {
				ms.domainList[i].Domain.MessageTTL = stringToInt(messageTTL)
			}
			toJSON(w, okResp{Message: "Domain has been updated"})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *Server) deleteDomain(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) getConnection(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) updateConnection(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) getTracking(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) updateClickTracking(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) updateOpenTracking(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) updateUnsubTracking(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) getTagLimits(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) activateDomainKey(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	toJSON(w, nil)
}

func (ms *Server) listDomainKeys(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var list []mtypes.DomainKey
	for _, domainKey := range ms.domainKeyList {
		list = append(list, domainKey)
	}

	toJSON(w, mtypes.ListDomainKeysResponse{
		Items: list,
	})
}

func (ms *Server) deactivateDomainKey(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	toJSON(w, nil)
}

func (ms *Server) updateDomainDkimAuthority(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, d := range ms.domainList {
		if d.Domain.Name == chi.URLParam(r, "domain") {
			if r.FormValue("self") == "" {
				toJSON(w, okResp{Message: "self param required"})
				return
			}
			toJSON(w, okResp{Message: "updated dkim authority"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "domain not found"})
}

func (ms *Server) updateDomainDkimSelector(w http.ResponseWriter, r *http.Request) {
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

func (ms *Server) updateWebPrefix(w http.ResponseWriter, r *http.Request) {
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
