package mailgun

import (
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/mux"
	"github.com/mailgun/mailgun-go/v4/events"
)

func (ms *MockServer) addEventRoutes(r *mux.Router) {
	r.HandleFunc("/{domain}/events", ms.listEvents).Methods(http.MethodGet)

	var (
		tags            = []string{"tag1", "tag2"}
		recipients      = []string{"one@mailgun.test", "two@mailgun.test"}
		recipientDomain = "mailgun.test"
		timeStamp       = TimeToFloat(time.Now().UTC())
		ipAddress       = "192.168.1.1"
		message         = events.Message{Headers: events.MessageHeaders{MessageID: "1234"}}
		clientInfo      = events.ClientInfo{
			AcceptLanguage: "EN",
			ClientName:     "Firefox",
			ClientOS:       "OS X",
			ClientType:     "browser",
			DeviceType:     "desktop",
			IP:             "8.8.8.8",
			UserAgent:      "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.12; rv:54.0) Gecko/20100101 Firefox/54.0",
		}
		geoLocation = events.GeoLocation{
			City:    "San Antonio",
			Country: "US",
			Region:  "TX",
		}
	)

	// AcceptedNoAuth
	accepted := new(events.Accepted)
	accepted.ID = randomString(16, "ID-")
	accepted.Message.Headers.MessageID = accepted.ID
	accepted.Name = events.EventAccepted
	accepted.Tags = tags
	accepted.Timestamp = timeStamp
	accepted.Recipient = recipients[0]
	accepted.RecipientDomain = recipientDomain
	accepted.Flags = events.Flags{
		IsAuthenticated: false,
	}
	ms.events = append(ms.events, accepted)

	// AcceptedAuth
	accepted = new(events.Accepted)
	accepted.ID = randomString(16, "ID-")
	accepted.Message.Headers.MessageID = accepted.ID
	accepted.Name = events.EventAccepted
	accepted.Tags = tags
	accepted.Timestamp = timeStamp
	accepted.Recipient = recipients[0]
	accepted.RecipientDomain = recipientDomain
	accepted.Campaigns = []events.Campaign{
		{ID: "test-id", Name: "test"},
	}
	accepted.Flags = events.Flags{
		IsAuthenticated: true,
	}
	ms.events = append(ms.events, accepted)

	// DeliveredSMTP
	delivered := new(events.Delivered)
	delivered.ID = randomString(16, "ID-")
	delivered.Message.Headers.MessageID = delivered.ID
	delivered.Name = events.EventDelivered
	delivered.Tags = tags
	delivered.Timestamp = timeStamp
	delivered.Recipient = recipients[0]
	delivered.RecipientDomain = recipientDomain
	delivered.DeliveryStatus.Message = "We sent an email Yo"
	delivered.Envelope = events.Envelope{
		Transport: "smtp",
		SendingIP: ipAddress,
	}
	delivered.Flags = events.Flags{
		IsAuthenticated: true,
	}
	ms.events = append(ms.events, delivered)

	// DeliveredHTTP
	delivered = new(events.Delivered)
	delivered.ID = randomString(16, "ID-")
	delivered.Message.Headers.MessageID = delivered.ID
	delivered.Name = events.EventDelivered
	delivered.Tags = tags
	delivered.Timestamp = timeStamp
	delivered.Recipient = recipients[0]
	delivered.RecipientDomain = recipientDomain
	delivered.DeliveryStatus.Message = "We sent an email Yo"
	delivered.Envelope = events.Envelope{
		Transport: "http",
		SendingIP: ipAddress,
	}
	delivered.Flags = events.Flags{
		IsAuthenticated: true,
	}
	ms.events = append(ms.events, delivered)

	// Stored
	stored := new(events.Stored)
	stored.ID = randomString(16, "ID-")
	stored.Name = events.EventStored
	stored.Tags = tags
	stored.Timestamp = timeStamp
	stored.Storage.URL = "http://mailgun.text/some/url"
	ms.events = append(ms.events, stored)

	// Clicked
	for _, recipient := range recipients {
		clicked := new(events.Clicked)
		clicked.ID = randomString(16, "ID-")
		clicked.Name = events.EventClicked
		clicked.Message = message
		clicked.Tags = tags
		clicked.Recipient = recipient
		clicked.ClientInfo = clientInfo
		clicked.GeoLocation = geoLocation
		clicked.Timestamp = timeStamp
		ms.events = append(ms.events, clicked)
	}

	clicked := new(events.Clicked)
	clicked.ID = randomString(16, "ID-")
	clicked.Name = events.EventClicked
	clicked.Message = message
	clicked.Tags = tags
	clicked.Recipient = recipients[0]
	clicked.ClientInfo = clientInfo
	clicked.GeoLocation = geoLocation
	clicked.Timestamp = timeStamp
	ms.events = append(ms.events, clicked)

	// Opened
	for _, recipient := range recipients {
		opened := new(events.Opened)
		opened.ID = randomString(16, "ID-")
		opened.Name = events.EventOpened
		opened.Message = message
		opened.Tags = tags
		opened.Recipient = recipient
		opened.ClientInfo = clientInfo
		opened.GeoLocation = geoLocation
		opened.Timestamp = timeStamp
		ms.events = append(ms.events, opened)
	}

	opened := new(events.Opened)
	opened.ID = randomString(16, "ID-")
	opened.Name = events.EventOpened
	opened.Message = message
	opened.Tags = tags
	opened.Recipient = recipients[0]
	opened.ClientInfo = clientInfo
	opened.GeoLocation = geoLocation
	opened.Timestamp = timeStamp
	ms.events = append(ms.events, opened)

	// Unsubscribed
	for _, recipient := range recipients {
		unsub := new(events.Unsubscribed)
		unsub.ID = randomString(16, "ID-")
		unsub.Name = events.EventUnsubscribed
		unsub.Tags = tags
		unsub.Recipient = recipient
		unsub.ClientInfo = clientInfo
		unsub.GeoLocation = geoLocation
		unsub.Timestamp = timeStamp
		ms.events = append(ms.events, unsub)
	}

	// Complained
	for _, recipient := range recipients {
		complained := new(events.Complained)
		complained.ID = randomString(16, "ID-")
		complained.Name = events.EventComplained
		complained.Tags = tags
		complained.Recipient = recipient
		complained.Timestamp = timeStamp
		ms.events = append(ms.events, complained)
	}
}

type eventsResponse struct {
	Items  []Event `json:"items"`
	Paging Paging  `json:"paging"`
}

func (ms *MockServer) listEvents(w http.ResponseWriter, r *http.Request) {
	var idx []string

	for _, e := range ms.events {
		idx = append(idx, e.GetID())
	}

	limit := stringToInt(r.FormValue("limit"))
	if limit == 0 {
		limit = 100
	}
	start, end := pageOffsets(idx, r.FormValue("page"), r.FormValue("address"), limit)

	var nextAddress, prevAddress string
	var results []Event

	if start != end {
		results = ms.events[start:end]
		nextAddress = results[len(results)-1].GetID()
		prevAddress = results[0].GetID()
	} else {
		results = []Event{}
		nextAddress = r.FormValue("address")
		prevAddress = r.FormValue("address")
	}

	resp := eventsResponse{
		Paging: Paging{
			First: getPageURL(r, url.Values{
				"page": []string{"first"},
			}),
			Last: getPageURL(r, url.Values{
				"page": []string{"last"},
			}),
			Next: getPageURL(r, url.Values{
				"page":    []string{"next"},
				"address": []string{nextAddress},
			}),
			Previous: getPageURL(r, url.Values{
				"page":    []string{"prev"},
				"address": []string{prevAddress},
			}),
		},
		Items: results,
	}
	toJSON(w, resp)
}
