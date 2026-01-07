package mocks

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addDomainKeysRoutes(r chi.Router) {
	ms.domainKeyList = append(
		ms.domainKeyList,
		mtypes.DomainKey{
			SigningDomain: "mailgun.test",
			Selector:      "pic",
			DNSRecord: mtypes.DNSRecord{
				Active:     true,
				Cached:     make([]string, 0),
				Name:       "mailgun.test",
				Priority:   "10",
				RecordType: "TXT",
				Valid:      "VALID",
				Value:      "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUA....",
			},
		},
		mtypes.DomainKey{
			SigningDomain: "mailgun.test",
			Selector:      "pic2",
			DNSRecord: mtypes.DNSRecord{
				Active:     true,
				Cached:     make([]string, 0),
				Name:       "mailgun.test",
				Priority:   "10",
				RecordType: "TXT",
				Valid:      "VALID",
				Value:      "k=rsa; p=MIGfMA0GCSqGSIb3DQEBAQUA....",
			},
		},
	)

	r.Get("/dkim/keys", ms.listAllDomainsKeys)
	r.Post("/dkim/keys", ms.createDomainKey)
	r.Delete("/dkim/keys", ms.deleteDomainKey)
}

func (ms *Server) listAllDomainsKeys(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}

	var list []mtypes.DomainKey
	for index, domainKey := range ms.domainKeyList {
		if index >= limit {
			break
		}

		list = append(list, domainKey)
	}

	toJSON(w, mtypes.ListAllDomainsKeysResponse{
		TotalCount: len(list),
		Items:      list,
	})
}

func (ms *Server) createDomainKey(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	toJSON(w, ms.domainKeyList[0])
}

func (ms *Server) deleteDomainKey(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	toJSON(w, nil)
}
