package mailgun

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
)

type MailingListContainer struct {
	MailingList MailingList
	Members     []Member
}

func (ms *mockServer) addMailingListRoutes(r *mux.Router) {
	r.HandleFunc("/lists/pages", ms.listMailingLists).Methods(http.MethodGet)
	r.HandleFunc("/lists/{address}", ms.getMailingList).Methods(http.MethodGet)
	r.HandleFunc("/lists", ms.createMailingList).Methods(http.MethodPost)
	r.HandleFunc("/lists/{address}", ms.updateMailingList).Methods(http.MethodPut)
	r.HandleFunc("/lists/{address}", ms.deleteMailingList).Methods(http.MethodDelete)

	r.HandleFunc("/lists/{address}/members/pages", ms.listMembers).Methods(http.MethodGet)
	r.HandleFunc("/lists/{address}/members/{member}", ms.getMember).Methods(http.MethodGet)
	r.HandleFunc("/lists/{address}/members", ms.createMember).Methods(http.MethodPost)
	r.HandleFunc("/lists/{address}/members/{member}", ms.updateMember).Methods(http.MethodPut)
	r.HandleFunc("/lists/{address}/members/{member}", ms.deleteMember).Methods(http.MethodDelete)
	r.HandleFunc("/lists/{address}/members.json", ms.bulkCreate).Methods(http.MethodPost)

	ms.mailingList = append(ms.mailingList, MailingListContainer{
		MailingList: MailingList{
			ReplyPreference: "list",
			AccessLevel:     "everyone",
			Address:         "foo@mailgun.test",
			CreatedAt:       RFC2822Time(time.Now().UTC()),
			Description:     "Mailgun developers list",
			MembersCount:    1,
			Name:            "",
		},
		Members: []Member{
			{
				Address: "dev@samples.mailgun.org",
				Name:    "Developer",
			},
		},
	})
}

func (ms *mockServer) listMailingLists(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var list []MailingList
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
		toJSON(w, listsResponse{})
		return
	}

	resp := listsResponse{
		Paging: Paging{
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

func (ms *mockServer) getMailingList(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, ml := range ms.mailingList {
		if ml.MailingList.Address == mux.Vars(r)["address"] {
			toJSON(w, mailingListResponse{MailingList: ml.MailingList})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "mailing list not found"})
}

func (ms *mockServer) deleteMailingList(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	result := ms.mailingList[:0]
	for _, ml := range ms.mailingList {
		if ml.MailingList.Address == mux.Vars(r)["address"] {
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

func (ms *mockServer) updateMailingList(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for i, d := range ms.mailingList {
		if d.MailingList.Address == mux.Vars(r)["address"] {
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
				ms.mailingList[i].MailingList.AccessLevel = AccessLevel(r.FormValue("access_level"))
			}
			if r.FormValue("reply_preference") != "" {
				ms.mailingList[i].MailingList.ReplyPreference = ReplyPreference(r.FormValue("reply_preference"))
			}
			toJSON(w, okResp{Message: "Mailing list member has been updated"})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "mailing list not found"})
}

func (ms *mockServer) createMailingList(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	ms.mailingList = append(ms.mailingList, MailingListContainer{
		MailingList: MailingList{
			CreatedAt:       RFC2822Time(time.Now().UTC()),
			Name:            r.FormValue("name"),
			Address:         r.FormValue("address"),
			Description:     r.FormValue("description"),
			AccessLevel:     AccessLevel(r.FormValue("access_level")),
			ReplyPreference: ReplyPreference(r.FormValue("reply_preference")),
		},
	})
	toJSON(w, okResp{Message: "Mailing list has been created"})
}

func (ms *mockServer) listMembers(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var list []Member
	var idx []string
	var found bool

	for _, ml := range ms.mailingList {
		if ml.MailingList.Address == mux.Vars(r)["address"] {
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
		toJSON(w, memberListResponse{})
		return
	}

	resp := memberListResponse{
		Paging: Paging{
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

func (ms *mockServer) getMember(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var found bool
	for _, ml := range ms.mailingList {
		if ml.MailingList.Address == mux.Vars(r)["address"] {
			found = true
			for _, member := range ml.Members {
				if member.Address == mux.Vars(r)["member"] {
					toJSON(w, memberResponse{Member: member})
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

func (ms *mockServer) deleteMember(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	idx := -1
	for i, ml := range ms.mailingList {
		if ml.MailingList.Address == mux.Vars(r)["address"] {
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
		if m.Address == mux.Vars(r)["member"] {
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

func (ms *mockServer) updateMember(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	idx := -1
	for i, ml := range ms.mailingList {
		if ml.MailingList.Address == mux.Vars(r)["address"] {
			idx = i
		}
	}

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "mailing list not found"})
		return
	}

	for i, m := range ms.mailingList[idx].Members {
		if m.Address == mux.Vars(r)["member"] {
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

func (ms *mockServer) createMember(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	idx := -1
	for i, ml := range ms.mailingList {
		if ml.MailingList.Address == mux.Vars(r)["address"] {
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

	ms.mailingList[idx].Members = append(ms.mailingList[idx].Members, Member{
		Name:       r.FormValue("name"),
		Address:    parseAddress(r.FormValue("address")),
		Vars:       stringToMap(r.FormValue("vars")),
		Subscribed: &sub,
	})
	toJSON(w, okResp{Message: "Mailing list member has been created"})
}

func (ms *mockServer) bulkCreate(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	idx := -1
	for i, ml := range ms.mailingList {
		if ml.MailingList.Address == mux.Vars(r)["address"] {
			idx = i
		}
	}

	if idx == -1 {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "mailing list not found"})
		return
	}

	var bulkList []Member
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
