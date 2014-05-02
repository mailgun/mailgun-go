package mailgun

import (
	"fmt"
	"github.com/mbanzon/simplehttp"
	"time"
)

// TODO(sfalvo):
// Abstract Paging/Links into an interface type or something which lets you page through
// data.  Borrow from Gophercloud's data interpreters.

// Events are open-ended, loosely-defined JSON documents.
// They will always have an event and a timestamp field, however.
type Event map[string]interface{}

// Links encapsulates navigation opportunities to find more information
// about things.
// TODO(sfalvo): Rename to Paging
type Links map[string]string

// noTime always equals an uninitialized Time structure.
// It's used to detect when a time parameter is provided.
var noTime time.Time

// GetEventsOptions lets the caller of GetEvents() specify how the results are to be returned.
// Begin and End time-box the results returned.
// ForceAscending and ForceDescending are used to force Mailgun to use a given traversal order of the events.
// If both ForceAscending and ForceDescending are true, an error will result.
// If none, the default will be inferred from the Begin and End parameters.
// Limit caps the number of results returned.  If left unspecified, Mailgun assumes 100.
// Compact, if true, compacts the returned JSON to minimize transmission bandwidth.
// Otherwise, the JSON is spaced appropriately for human consumption.
// Filter allows the caller to provide more specialized filters on the query.
// Consult the Mailgun documentation for more details.
type GetEventsOptions struct {
	Begin, End                               time.Time
	ForceAscending, ForceDescending, Compact bool
	Limit                                    int
	Filter                                   map[string]string
}

type EventIterator struct {
	events []Event
	nextURL, prevURL string
	mg Mailgun
}

// GetEvents provides the caller with a list of log entries.
// See the GetEventsOptions structure for information on how to customize the list returned.
// Note that the API responds with events with open definitions;
// that is, no specific standard structure exists for them.
// Thus, you'll need to provide your own accessors to the information of interest.
//
// DEPRECATED.  Use GetEventsIterator() instead.
func (mg *MailgunImpl) GetEvents(opts GetEventsOptions) ([]Event, Links, error) {
	ei, err := mg.GetEventsIterator(opts)
	if err != nil {
		return nil, nil, err
	}
	return ei.Events(), nil, nil
}

func (mg *MailgunImpl) GetEventsIterator(opts GetEventsOptions) (*EventIterator, error) {
	ei := mg.NewEventIterator()
	err := ei.GetFirstPage(opts)
	return ei, err
}

func (mg *MailgunImpl) NewEventIterator() *EventIterator {
	return &EventIterator{mg: mg}
}

func (ei *EventIterator) Events() []Event {
	return ei.events
}

func (ei *EventIterator) IsAtBeginning() bool {
	return ei.prevURL == ""
}

func (ei *EventIterator) IsAtEnd() bool {
	return ei.nextURL == ""
}

func (ei *EventIterator) GetFirstPage(opts GetEventsOptions) error {
	if opts.ForceAscending && opts.ForceDescending {
		return fmt.Errorf("collation cannot at once be both ascending and descending")
	}

	payload := simplehttp.NewUrlEncodedPayload()
	if opts.Limit != 0 {
		payload.AddValue("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Compact {
		payload.AddValue("pretty", "no")
	}
	if opts.ForceAscending {
		payload.AddValue("ascending", "yes")
	}
	if opts.ForceDescending {
		payload.AddValue("ascending", "no")
	}
	if opts.Begin != noTime {
		payload.AddValue("begin", formatMailgunTime(&opts.Begin))
	}
	if opts.End != noTime {
		payload.AddValue("end", formatMailgunTime(&opts.End))
	}
	if opts.Filter != nil {
		for k, v := range opts.Filter {
			payload.AddValue(k, v)
		}
	}

	url, err := generateParameterizedUrl(ei.mg, eventsEndpoint, payload)
	if err != nil {
		return err
	}
	r := simplehttp.NewHTTPRequest(url)
	r.SetBasicAuth(basicAuthUser, ei.mg.ApiKey())
	var response map[string]interface{}
	err = getResponseFromJSON(r, &response)
	if err != nil {
		return err
	}

	items := response["items"].([]interface{})
	ei.events = make([]Event, len(items))
	for i, item := range items {
		ei.events[i] = item.(map[string]interface{})
	}

	pagings := response["paging"].(map[string]interface{})
	links := make(map[string]string, len(pagings))
	for key, page := range pagings {
		links[key] = page.(string)
	}
	ei.nextURL = links["next"]
	ei.prevURL = links["previous"]
	return err
}
