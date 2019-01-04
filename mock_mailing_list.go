package mailgun

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-chi/chi"
)

type mailingListContainer struct {
	MailingList MailingList
	Members     []Member
}

type paging struct {
	First    string `json:"first"`
	Last     string `json:"last"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
}

type mailingListResponse struct {
	Items  []MailingList `json:"items"`
	Paging paging        `json:"paging"`
}

func (ms *MockServer) addMailingListRoutes(r chi.Router) {
	r.Get("/lists/pages", ms.listMailingLists)
	/*r.Get("/lists/{address}", ms.getMailingList)
	r.Post("/lists", ms.createMailingList)
	r.Put("/lists/{address}", ms.updateMailingList)
	r.Delete("/lists/{address}", ms.deleteMailingList)

	r.Get("/lists/{address}/members/pages", ms.listMembers)
	r.Get("/lists/{address}/members/{address}", ms.getMember)
	r.Post("/lists/{address}/members", ms.createMember)
	r.Post("/lists/{address}/members", ms.createMember)
	r.Put("/lists/{address}/members/{address}", ms.updateMember)
	r.Delete("/lists/{address}/members/{address}", ms.deleteMember)
	r.Post("/lists/{address}/members.json", ms.bulkCreate)*/

	ms.mailingList = append(ms.mailingList, mailingListContainer{
		MailingList: MailingList{
			AccessLevel:  "everyone",
			Address:      "foo@mailgun.test",
			CreatedAt:    "Tue, 06 Mar 2012 05:44:45 GMT",
			Description:  "Mailgun developers list",
			MembersCount: 1,
			Name:         "",
		},
		Members: []Member{
			{
				Address: "dev@samples.mailgun.org",
				Name:    "Developer",
			},
		},
	})

	for i := 0; i < 20; i++ {
		ms.mailingList = append(ms.mailingList, mailingListContainer{
			MailingList: MailingList{
				AccessLevel:  "everyone",
				Address:      fmt.Sprintf("%0d@mailgun.test", i),
				CreatedAt:    "Tue, 06 Mar 2012 05:44:45 GMT",
				Description:  "Mailgun developers list",
				MembersCount: 1,
				Name:         "",
			},
			Members: []Member{
				{
					Address: "dev@samples.mailgun.org",
					Name:    "Developer",
				},
			},
		})
	}
}

func (ms *MockServer) listMailingLists(w http.ResponseWriter, r *http.Request) {
	var list []MailingList
	var idx []string

	for _, item := range ms.mailingList {
		list = append(list, item.MailingList)
		idx = append(idx, item.MailingList.Address)
	}

	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}
	start, end := pageOffsets(idx, r.FormValue("page"), r.FormValue("address"), limit)
	results := list[start:end]

	if len(results) == 0 {
		toJSON(w, mailingListResponse{})
		return
	}

	resp := mailingListResponse{
		Paging: paging{
			First: getPageURL(r, url.Values{
				"page": []string{"first"},
			}),
			Last: getPageURL(r, url.Values{
				"page": []string{"last"},
			}),
			Next: getPageURL(r, url.Values{
				"page":    []string{"next"},
				"address": []string{results[len(results)-1].Address},
			}),
			Previous: getPageURL(r, url.Values{
				"page":    []string{"prev"},
				"address": []string{results[0].Address},
			}),
		},
		Items: results,
	}
	toJSON(w, resp)
}

func getPageURL(r *http.Request, params url.Values) string {
	params.Add("limit", r.FormValue("limit"))
	return "http://" + r.Host + r.URL.EscapedPath() + "?" + params.Encode()
}
