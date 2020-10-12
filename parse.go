package mailgun

import (
	"fmt"
	"reflect"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/mailgun/mailgun-go/v4/events"
)

// All events returned by the EventIterator conform to this interface
type Event interface {
	GetName() string
	SetName(name string)
	GetTimestamp() time.Time
	SetTimestamp(time.Time)
	GetID() string
	SetID(id string)
}

// A list of all JSON event types returned by the /events API
var EventNames = map[string]func() Event{
	"accepted":                 new_(events.Accepted{}),
	"clicked":                  new_(events.Clicked{}),
	"complained":               new_(events.Complained{}),
	"delivered":                new_(events.Delivered{}),
	"failed":                   new_(events.Failed{}),
	"opened":                   new_(events.Opened{}),
	"rejected":                 new_(events.Rejected{}),
	"stored":                   new_(events.Stored{}),
	"unsubscribed":             new_(events.Unsubscribed{}),
	"list_member_uploaded":     new_(events.ListMemberUploaded{}),
	"list_member_upload_error": new_(events.ListMemberUploadError{}),
	"list_uploaded":            new_(events.ListUploaded{}),
}

// new_ is a universal event "constructor".
func new_(e interface{}) func() Event {
	typ := reflect.TypeOf(e)
	return func() Event {
		return reflect.New(typ).Interface().(Event)
	}
}

func parseResponse(raw []byte) ([]Event, error) {
	var resp events.Response
	if err := jsoniter.Unmarshal(raw, &resp); err != nil {
		return nil, fmt.Errorf("failed to un-marshall event.Response: %s", err)
	}

	var result []Event
	for _, value := range resp.Items {
		event, err := ParseEvent(value)
		if err != nil {
			return nil, fmt.Errorf("while parsing event: %s", err)
		}
		result = append(result, event)
	}
	return result, nil
}

// Given a slice of events.RawJSON events return a slice of Event for each parsed event
func ParseEvents(raw []events.RawJSON) ([]Event, error) {
	var result []Event
	for _, value := range raw {
		event, err := ParseEvent(value)
		if err != nil {
			return nil, fmt.Errorf("while parsing event: %s", err)
		}
		result = append(result, event)
	}
	return result, nil
}

// Parse converts raw bytes data into an event struct. Can accept events.RawJSON as input
func ParseEvent(raw []byte) (Event, error) {
	// Try to recognize the event first.
	var e events.EventName
	if err := jsoniter.Unmarshal(raw, &e); err != nil {
		return nil, fmt.Errorf("failed to recognize event: %v", err)
	}

	// Get the event "constructor" from the map.
	newEvent, ok := EventNames[e.GetName()]
	if !ok {
		return nil, fmt.Errorf("unsupported event: '%s'", e.GetName())
	}
	event := newEvent()

	// Parse the known event.
	if err := jsoniter.Unmarshal(raw, event); err != nil {
		return nil, fmt.Errorf("failed to parse event '%s': %v", e.GetName(), err)
	}

	return event, nil
}
