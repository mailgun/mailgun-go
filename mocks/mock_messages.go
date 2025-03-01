package mocks

import (
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/events"
	"github.com/mailgun/mailgun-go/v5/mtypes"
)

func (ms *Server) addMessagesRoutes(r chi.Router) {
	r.Post("/{domain}/messages", ms.createMessages)

	// This path is made up; it could be anything as the storage url could change over time
	r.Get("/se.storage.url/messages/{id}", ms.getStoredMessages)
	r.Post("/se.storage.url/messages/{id}", ms.sendStoredMessages)
}

// TODO: This implementation doesn't support multiple recipients
func (ms *Server) createMessages(w http.ResponseWriter, r *http.Request) {
	to, err := mail.ParseAddress(r.FormValue("to"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "invalid 'to' address"})
		return
	}
	id := randomString(16, "ID-")

	switch to.Address {
	case "stored@mailgun.test":
		stored := new(events.Stored)
		stored.Name = events.EventStored
		stored.Timestamp = mailgun.TimeToFloat(time.Now().UTC())
		stored.ID = id
		stored.Storage.URL = ms.URL() + "/v3/se.storage.url/messages/" + id
		stored.Storage.Key = id
		stored.Message.Headers = events.MessageHeaders{
			Subject:   r.FormValue("subject"),
			From:      r.FormValue("from"),
			To:        to.Address,
			MessageID: id,
		}
		stored.Message.Recipients = []string{
			r.FormValue("to"),
		}
		stored.Message.Size = 10
		stored.Flags = events.Flags{
			IsTestMode: false,
		}
		ms.mutex.Lock()
		ms.events = append(ms.events, stored)
		ms.mutex.Unlock()
	default:
		accepted := new(events.Accepted)
		accepted.Name = events.EventAccepted
		accepted.ID = id
		accepted.Timestamp = mailgun.TimeToFloat(time.Now().UTC())
		accepted.Message.Headers.From = r.FormValue("from")
		accepted.Message.Headers.To = r.FormValue("to")
		accepted.Message.Headers.MessageID = accepted.ID
		accepted.Message.Headers.Subject = r.FormValue("subject")

		if r.MultipartForm.File != nil {
			for _, fh := range r.MultipartForm.File {
				for _, fd := range fh {
					accepted.Message.Attachments = append(accepted.Message.Attachments, events.Attachment{
						FileName:    fd.Filename,
						ContentType: fd.Header.Get("Content-Type"),
						Size:        int(fd.Size),
					})
				}
			}
		}
		accepted.Recipient = r.FormValue("to")
		accepted.RecipientDomain = strings.Split(to.Address, "@")[1]
		accepted.Flags = events.Flags{
			IsAuthenticated: true,
		}
		ms.mutex.Lock()
		ms.events = append(ms.events, accepted)
		ms.mutex.Unlock()
	}

	tags := r.Form["o:tag"]
	for _, newTag := range tags {
		var tagExists bool
		for _, existingTag := range ms.tags {
			if newTag == existingTag.Value {
				tagExists = true
				break
			}
		}

		if !tagExists {
			ms.tags = append(ms.tags, mtypes.Tag{Value: newTag})
		}
	}

	toJSON(w, okResp{ID: "<" + id + ">", Message: "Queued. Thank you."})
}

func (ms *Server) getStoredMessages(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	// Find our stored event
	var stored *events.Stored
	for _, event := range ms.events {
		if event.GetID() == id {
			stored = event.(*events.Stored)
		}
	}

	if stored == nil {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "not found"})
	}

	toJSON(w, mtypes.StoredMessage{
		Recipients: strings.Join(stored.Message.Recipients, ","),
		Sender:     stored.Message.Headers.From,
		Subject:    stored.Message.Headers.Subject,
		From:       stored.Message.Headers.From,
		MessageHeaders: [][]string{
			{"Sender", stored.Message.Headers.From},
			{"To", stored.Message.Headers.To},
			{"Subject", stored.Message.Headers.Subject},
			{"Content-Type", "text/plain"},
		},
	})
}

func (ms *Server) sendStoredMessages(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	defer ms.mutex.Unlock()
	ms.mutex.Lock()

	// Find our stored event
	var stored *events.Stored
	for _, event := range ms.events {
		if event.GetID() == id {
			stored = event.(*events.Stored)
		}
	}

	if stored == nil {
		w.WriteHeader(http.StatusNotFound)
		toJSON(w, okResp{Message: "not found"})
	}

	// DO NOTHING

	toJSON(w, okResp{ID: "<" + id + ">", Message: "Queued. Thank you."})
}
