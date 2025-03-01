package mocks

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addWebhookRoutes(r chi.Router) {
	r.Route("/domains/{domain}/webhooks", func(sr chi.Router) {
		sr.Get("/", ms.listWebHooks)
		sr.Post("/", ms.postWebHook)
		sr.Get("/{webhook}", ms.getWebHook)
		sr.Put("/{webhook}", ms.putWebHook)
		sr.Delete("/{webhook}", ms.deleteWebHook)
	})

	ms.webhooks = mtypes.WebHooksListResponse{
		Webhooks: map[string]mtypes.UrlOrUrls{
			"new-webhook": {
				Urls: []string{"http://example.com/new"},
			},
			"legacy-webhook": {
				Url: "http://example.com/legacy",
			},
		},
	}
}

func (ms *Server) listWebHooks(w http.ResponseWriter, _ *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	toJSON(w, ms.webhooks)
}

func (ms *Server) getWebHook(w http.ResponseWriter, r *http.Request) {
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	resp := mtypes.WebHookResponse{
		Webhook: mtypes.UrlOrUrls{
			Urls: ms.webhooks.Webhooks[chi.URLParam(r, "webhook")].Urls,
		},
	}
	toJSON(w, resp)
}

func (ms *Server) postWebHook(w http.ResponseWriter, r *http.Request) {
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
	ms.webhooks.Webhooks[r.FormValue("id")] = mtypes.UrlOrUrls{Urls: urls}

	toJSON(w, okResp{Message: "success"})
}

func (ms *Server) putWebHook(w http.ResponseWriter, r *http.Request) {
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
	ms.webhooks.Webhooks[chi.URLParam(r, "webhook")] = mtypes.UrlOrUrls{Urls: urls}

	toJSON(w, okResp{Message: "success"})
}

func (ms *Server) deleteWebHook(w http.ResponseWriter, r *http.Request) {
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
