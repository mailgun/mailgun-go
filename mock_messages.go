package mailgun

import (
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/mailgun/mailgun-go/events"
)

func (ms *MockServer) addMessagesRoutes(r chi.Router) {
	r.Post("/{domain}/messages", ms.createMessages)
}

// TODO: This implementation doesn't support multiple recipients
func (ms *MockServer) createMessages(w http.ResponseWriter, r *http.Request) {
	to, err := mail.ParseAddress(r.FormValue("to"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		toJSON(w, okResp{Message: "invalid 'to' address"})
		return
	}

	accepted := new(events.Accepted)
	accepted.ID = randomString(16, "ID-")
	accepted.Name = events.EventAccepted
	accepted.Timestamp = TimeToFloat(time.Now().UTC())
	accepted.Message.Headers.From = r.FormValue("from")
	accepted.Message.Headers.To = r.FormValue("to")
	accepted.Message.Headers.MessageID = accepted.ID
	accepted.Message.Headers.Subject = r.FormValue("subject")

	accepted.Recipient = r.FormValue("to")
	accepted.RecipientDomain = strings.Split(to.Address, "@")[1]
	accepted.Flags = events.Flags{
		IsAuthenticated: true,
	}
	ms.events = append(ms.events, accepted)

	toJSON(w, okResp{ID: "<" + accepted.ID + ">", Message: "Queued. Thank you."})
}
