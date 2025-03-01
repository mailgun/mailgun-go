package mocks

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v4/mtypes"
)

func (ms *Server) addSubaccountRoutes(r chi.Router) {
	ms.subaccountList = append(ms.subaccountList, mtypes.Subaccount{
		ID:     "enabled.subaccount",
		Name:   "mailgun.test",
		Status: "enabled",
	}, mtypes.Subaccount{
		ID:     "disabled.subaccount",
		Name:   "mailgun.test",
		Status: "disabled",
	})

	r.Get("/accounts/subaccounts", ms.listSubaccounts)
	r.Post("/accounts/subaccounts", ms.createSubaccount)

	r.Get("/accounts/subaccounts/{subaccountID}", ms.getSubaccount)
	r.Post("/accounts/subaccounts/{subaccountID}/enable", ms.enableSubaccount)
	r.Post("/accounts/subaccounts/{subaccountID}/disable", ms.disableSubaccount)
}

func (ms *Server) listSubaccounts(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	var list mtypes.ListSubaccountsResponse
	for _, subaccount := range ms.subaccountList {
		list.Items = append(list.Items, subaccount)
	}

	skip := stringToInt(r.FormValue("skip"))
	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}

	if skip > len(list.Items) {
		skip = len(list.Items)
	}

	end := limit + skip
	if end > len(list.Items) {
		end = len(list.Items)
	}

	// If we are at the end of the list
	if skip == end {
		toJSON(w, mtypes.ListSubaccountsResponse{
			Total: len(list.Items),
			Items: []mtypes.Subaccount{},
		})
		return
	}

	toJSON(w, mtypes.ListSubaccountsResponse{
		Total: len(list.Items),
		Items: list.Items[skip:end],
	})
}

func (ms *Server) getSubaccount(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, s := range ms.subaccountList {
		if s.ID == chi.URLParam(r, "subaccountID") {
			toJSON(w, mtypes.SubaccountResponse{Item: s})
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	toJSON(w, okResp{Message: "Not Found"})
}

func (ms *Server) createSubaccount(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	ms.subaccountList = append(ms.subaccountList, mtypes.Subaccount{
		ID:     "test",
		Name:   r.FormValue("name"),
		Status: "active",
	})
	toJSON(w, okResp{Message: "Subaccount has been created"})
}

func (ms *Server) enableSubaccount(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, subaccount := range ms.subaccountList {
		if subaccount.ID == chi.URLParam(r, "subaccountID") && subaccount.Status == "disabled" {
			subaccount.Status = "enabled"
			toJSON(w, mtypes.SubaccountResponse{Item: subaccount})
			return
		}
		if subaccount.ID == chi.URLParam(r, "subaccountID") && subaccount.Status == "enabled" {
			toJSON(w, okResp{Message: "subaccount is already enabled"})
			return
		}
	}
	toJSON(w, okResp{Message: "Not Found"})
}

func (ms *Server) disableSubaccount(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	for _, subaccount := range ms.subaccountList {
		if subaccount.ID == chi.URLParam(r, "subaccountID") && subaccount.Status == "enabled" {
			subaccount.Status = "disabled"
			toJSON(w, mtypes.SubaccountResponse{Item: subaccount})
			return
		}
		if subaccount.ID == chi.URLParam(r, "subaccountID") && subaccount.Status == "disabled" {
			toJSON(w, okResp{Message: "subaccount is already disabled"})
			return
		}
	}
	toJSON(w, okResp{Message: "Not Found"})
}
