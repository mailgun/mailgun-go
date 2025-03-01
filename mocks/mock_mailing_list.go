package mocks

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v4/mtypes"
)

type MailingListContainer struct {
	MailingList mtypes.MailingList
	Members     []mtypes.Member
}

func (ms *Server) addMailingListRoutes(r chi.Router) {
	r.Get("/lists/pages", ms.listMailingLists)
	r.Get("/lists/{address}", ms.getMailingList)
	r.Post("/lists", ms.createMailingList)
	r.Put("/lists/{address}", ms.updateMailingList)
	r.Delete("/lists/{address}", ms.deleteMailingList)

	r.Get("/lists/{address}/members/pages", ms.listMembers)
	r.Get("/lists/{address}/members/{member}", ms.getMember)
	r.Post("/lists/{address}/members", ms.createMember)
	r.Put("/lists/{address}/members/{member}", ms.updateMember)
	r.Delete("/lists/{address}/members/{member}", ms.deleteMember)
	r.Post("/lists/{address}/members.json", ms.bulkCreate)

	ms.mailingList = append(ms.mailingList, MailingListContainer{
		MailingList: mtypes.MailingList{
			ReplyPreference: "list",
			AccessLevel:     "everyone",
			Address:         "foo@mailgun.test",
			CreatedAt:       mtypes.RFC2822Time(time.Now().UTC()),
			Description:     "Mailgun developers list",
			MembersCount:    1,
			Name:            "",
		},
		Members: []mtypes.Member{
			{
				Address: "dev@samples.mailgun.org",
				Name:    "Developer",
			},
		},
	})
}

func (ms *Server) listMailingLists(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var list []mtypes.MailingList
	var idx []string

	for _, ml := range ms.mailingList {
		list = append(list, ml.MailingList)
		idx = append(idx, ml.MailingList.Address)
	}

	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}
	start, end := pageOffsets(idx, r.FormValue("page"), r.FormValue("address"), limit)
	results := list[start:end]

	if len(results) == 0 {
		toJSON(w, mtypes.ListMailingListsResponse{})
		return
	}

	resp := mtypes.ListMailingListsResponse{
		Paging: mtypes.Paging{
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

func (ms *Server) getMailingList(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, ml := range ms.mailingList {
		if ml.MailingList.Address == chi.URLParam(r, "address") {
			toJSON(w, mtypes.GetMailingListResponse{MailingList: ml.MailingList})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "mailing list not found"})
}

func (ms *Server) deleteMailingList(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	result := ms.mailingList[:0]
	for _, ml := range ms.mailingList {
		if ml.MailingList.Address == chi.URLParam(r, "address") {
			continue
		}
		result = append(result, ml)
	}

	if len(result) != len(ms.mailingList) {
		toJSON(w, okResp{Message: "success"})
		ms.mailingList = result
		return
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "mailing list not found"})
}

func (ms *Server) updateMailingList(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for i, d := range ms.mailingList {
		if d.MailingList.Address == chi.URLParam(r, "address") {
			if r.FormValue("address") != "" {
				ms.mailingList[i].MailingList.Address = r.FormValue("address")
			}
			if r.FormValue("name") != "" {
				ms.mailingList[i].MailingList.Name = r.FormValue("name")
			}
			if r.FormValue("description") != "" {
				ms.mailingList[i].MailingList.Description = r.FormValue("description")
			}
			if r.FormValue("access_level") != "" {
				ms.mailingList[i].MailingList.AccessLevel = mtypes.AccessLevel(r.FormValue("access_level"))
			}
			if r.FormValue("reply_preference") != "" {
				ms.mailingList[i].MailingList.ReplyPreference = mtypes.ReplyPreference(r.FormValue("reply_preference"))
			}
			toJSON(w, okResp{Message: "Mailing list member has been updated"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "mailing list not found"})
}

func (ms *Server) createMailingList(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	ms.mailingList = append(ms.mailingList, MailingListContainer{
		MailingList: mtypes.MailingList{
			CreatedAt:       mtypes.RFC2822Time(time.Now().UTC()),
			Name:            r.FormValue("name"),
			Address:         r.FormValue("address"),
			Description:     r.FormValue("description"),
			AccessLevel:     mtypes.AccessLevel(r.FormValue("access_level")),
			ReplyPreference: mtypes.ReplyPreference(r.FormValue("reply_preference")),
		},
	})
	toJSON(w, okResp{Message: "Mailing list has been created"})
}

func (ms *Server) listMembers(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var list []mtypes.Member
	var idx []string
	var found bool

	for _, ml := range ms.mailingList {
		if ml.MailingList.Address == chi.URLParam(r, "address") {
			found = true
			for _, member := range ml.Members {
				list = append(list, member)
				idx = append(idx, member.Address)
			}
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "mailing list not found"})
		return
	}

	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}
	start, end := pageOffsets(idx, r.FormValue("page"), r.FormValue("address"), limit)
	results := list[start:end]

	if len(results) == 0 {
		toJSON(w, mtypes.MemberListResponse{})
		return
	}

	resp := mtypes.MemberListResponse{
		Paging: mtypes.Paging{
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
		Lists: results,
	}
	toJSON(w, resp)
}

func (ms *Server) getMember(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var found bool
	for _, ml := range ms.mailingList {
		if ml.MailingList.Address == chi.URLParam(r, "address") {
			found = true
			for _, member := range ml.Members {
				if member.Address == chi.URLParam(r, "member") {
					toJSON(w, mtypes.MemberResponse{Member: member})
					return
				}
			}
		}
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "mailing list not found"})
		return
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "member not found"})
}

func (ms *Server) deleteMember(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	idx := -1
	for i, ml := range ms.mailingList {
		if ml.MailingList.Address == chi.URLParam(r, "address") {
			idx = i
		}
	}

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "mailing list not found"})
		return
	}

	result := ms.mailingList[idx].Members[:0]
	for _, m := range ms.mailingList[idx].Members {
		if m.Address == chi.URLParam(r, "member") {
			continue
		}
		result = append(result, m)
	}

	if len(result) != len(ms.mailingList[idx].Members) {
		toJSON(w, okResp{Message: "Mailing list member has been deleted"})
		ms.mailingList[idx].Members = result
		return
	}

	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "member not found"})
}

func (ms *Server) updateMember(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	idx := -1
	for i, ml := range ms.mailingList {
		if ml.MailingList.Address == chi.URLParam(r, "address") {
			idx = i
		}
	}

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "mailing list not found"})
		return
	}

	for i, m := range ms.mailingList[idx].Members {
		if m.Address == chi.URLParam(r, "member") {
			if r.FormValue("address") != "" {
				ms.mailingList[idx].Members[i].Address = parseAddress(r.FormValue("address"))
			}
			if r.FormValue("name") != "" {
				ms.mailingList[idx].Members[i].Name = r.FormValue("name")
			}
			if r.FormValue("vars") != "" {
				ms.mailingList[idx].Members[i].Vars = stringToMap(r.FormValue("vars"))
			}
			if r.FormValue("subscribed") != "" {
				sub := stringToBool(r.FormValue("subscribed"))
				ms.mailingList[idx].Members[i].Subscribed = &sub
			}
			toJSON(w, okResp{Message: "Mailing list member has been updated"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "member not found"})
}

func (ms *Server) createMember(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	idx := -1
	for i, ml := range ms.mailingList {
		if ml.MailingList.Address == chi.URLParam(r, "address") {
			idx = i
		}
	}

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "mailing list not found"})
		return
	}

	sub := stringToBool(r.FormValue("subscribed"))

	if len(ms.mailingList[idx].Members) != 0 {
		for i, m := range ms.mailingList[idx].Members {
			if m.Address == r.FormValue("address") {
				if !stringToBool(r.FormValue("upsert")) {
					w.WriteHeader(http.StatusConflict)
					toJSON(w, okResp{Message: "member already exists"})
					return
				}

				ms.mailingList[idx].Members[i].Address = parseAddress(r.FormValue("address"))
				ms.mailingList[idx].Members[i].Name = r.FormValue("name")
				ms.mailingList[idx].Members[i].Vars = stringToMap(r.FormValue("vars"))
				ms.mailingList[idx].Members[i].Subscribed = &sub
				break
			}
		}
	}

	ms.mailingList[idx].Members = append(ms.mailingList[idx].Members, mtypes.Member{
		Name:       r.FormValue("name"),
		Address:    parseAddress(r.FormValue("address")),
		Vars:       stringToMap(r.FormValue("vars")),
		Subscribed: &sub,
	})
	toJSON(w, okResp{Message: "Mailing list member has been created"})
}

func (ms *Server) bulkCreate(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	idx := -1
	for i, ml := range ms.mailingList {
		if ml.MailingList.Address == chi.URLParam(r, "address") {
			idx = i
		}
	}

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "mailing list not found"})
		return
	}

	var bulkList []mtypes.Member
	if err := json.Unmarshal([]byte(r.FormValue("members")), &bulkList); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		toJSON(w, okResp{Message: "while un-marshalling 'members' param - " + err.Error()})
		return
	}

BULK:
	for _, member := range bulkList {
		member.Address = parseAddress(member.Address)
		if len(ms.mailingList[idx].Members) != 0 {
			for i, m := range ms.mailingList[idx].Members {
				if m.Address == member.Address {
					if !stringToBool(r.FormValue("upsert")) {
						w.WriteHeader(http.StatusConflict)
						toJSON(w, okResp{Message: "member already exists"})
						return
					}
					ms.mailingList[idx].Members[i] = member
					continue BULK
				}
			}
		}
		ms.mailingList[idx].Members = append(ms.mailingList[idx].Members, member)
	}
	toJSON(w, okResp{Message: "Mailing list has been updated"})
}
