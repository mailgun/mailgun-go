package mailgun

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (ms *mockServer) addWebhookRoutes(r chi.Router) {
	r.Route("/domains/{domain}/webhooks", func(sr chi.Router) {
		sr.Get("/", ms.listWebHooks)
		sr.Post("/", ms.postWebHook)
		sr.Get("/{webhook}", ms.getWebHook)
		sr.Put("/{webhook}", ms.putWebHook)
		sr.Delete("/{webhook}", ms.deleteWebHook)
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

func (ms *mockServer) listWebHooks(w http.ResponseWriter, _ *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	toJSON(w, ms.webhooks)
}

func (ms *mockServer) getWebHook(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	resp := WebHookResponse{
		Webhook: UrlOrUrls{
			Urls: ms.webhooks.Webhooks[chi.URLParam(r, "webhook")].Urls,
		},
	}
	toJSON(w, resp)
}

func (ms *mockServer) postWebHook(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

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

func (ms *mockServer) putWebHook(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

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

func (ms *mockServer) deleteWebHook(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	_, ok := ms.webhooks.Webhooks[chi.URLParam(r, "webhook")]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "webhook not found"})
	}

	delete(ms.webhooks.Webhooks, chi.URLParam(r, "webhook"))
	toJSON(w, okResp{Message: "success"})
}
