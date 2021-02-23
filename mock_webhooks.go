package mailgun

import (
	"net/http"

	"github.com/gorilla/mux"
)

func (ms *MockServer) addWebhookRoutes(r *mux.Router) {
	sr := r.PathPrefix("/domains/{domain}/webhooks").Subrouter()
	sr.HandleFunc("", ms.listWebHooks).Methods(http.MethodGet)
	sr.HandleFunc("", ms.postWebHook).Methods(http.MethodPost)
	sr.HandleFunc("/{webhook}", ms.getWebHook).Methods(http.MethodGet)
	sr.HandleFunc("/{webhook}", ms.putWebHook).Methods(http.MethodPut)
	sr.HandleFunc("/{webhook}", ms.deleteWebHook).Methods(http.MethodDelete)

	ms.webhooks = WebHooksListResponse{
		Webhooks: map[string]UrlOrUrls{
			"new-webhook": {
				Urls: []string{"http://example.com/new"},
			},
			"legacy-webhook": {
				Url: "http://example.com/legacy",
			},
		},
	}
}

func (ms *MockServer) listWebHooks(w http.ResponseWriter, _ *http.Request) {
	toJSON(w, ms.webhooks)
}

func (ms *MockServer) getWebHook(w http.ResponseWriter, r *http.Request) {
	resp := WebHookResponse{
		Webhook: UrlOrUrls{
			Urls: ms.webhooks.Webhooks[mux.Vars(r)["webhook"]].Urls,
		},
	}
	toJSON(w, resp)
}

func (ms *MockServer) postWebHook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: err.Error()})
		return
	}

	var urls []string
	for _, url := range r.Form["url"] {
		urls = append(urls, url)
	}
	ms.webhooks.Webhooks[r.FormValue("id")] = UrlOrUrls{Urls: urls}

	toJSON(w, okResp{Message: "success"})
}

func (ms *MockServer) putWebHook(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: err.Error()})
		return
	}

	var urls []string
	for _, url := range r.Form["url"] {
		urls = append(urls, url)
	}
	ms.webhooks.Webhooks[mux.Vars(r)["webhook"]] = UrlOrUrls{Urls: urls}

	toJSON(w, okResp{Message: "success"})
}

func (ms *MockServer) deleteWebHook(w http.ResponseWriter, r *http.Request) {
	_, ok := ms.webhooks.Webhooks[mux.Vars(r)["webhook"]]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "webhook not found"})
	}

	delete(ms.webhooks.Webhooks, mux.Vars(r)["webhook"])
	toJSON(w, okResp{Message: "success"})
}
