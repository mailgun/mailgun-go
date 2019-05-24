package mailgun

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (ms *MockServer) addWebhookRoutes(r chi.Router) {
	r.Route("/domains/{domain}/webhooks", func(r chi.Router) {
		r.Get("/", ms.listWebHooks)
		r.Post("/", ms.postWebHook)
		r.Get("/{webhook}", ms.getWebHook)
		r.Put("/{webhook}", ms.putWebHook)
		r.Delete("/{webhook}", ms.deleteWebHook)
	})
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
			Urls: ms.webhooks.Webhooks[chi.URLParam(r, "webhook")].Urls,
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
	ms.webhooks.Webhooks[chi.URLParam(r, "webhook")] = UrlOrUrls{Urls: urls}

	toJSON(w, okResp{Message: "success"})
}

func (ms *MockServer) deleteWebHook(w http.ResponseWriter, r *http.Request) {
	_, ok := ms.webhooks.Webhooks[chi.URLParam(r, "webhook")]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "webhook not found"})
	}

	delete(ms.webhooks.Webhooks, chi.URLParam(r, "webhook"))
	toJSON(w, okResp{Message: "success"})
}
