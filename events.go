package mailgun

import (
	"github.com/mbanzon/simplehttp"
	"time"
	"fmt"
)

// Events are open-ended, loosely-defined JSON documents.
// They will always have an event and a timestamp field, however.
type Event map[string]interface{}

// Links encapsulates navigation opportunities to find more information
// about things.
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
	Begin, End	time.Time
	ForceAscending, ForceDescending, Compact	bool
	Limit	int
	Filter map[string]string
}

func (mg *mailgunImpl) GetEvents(opts GetEventsOptions) ([]Event, Links, error) {
	if opts.ForceAscending && opts.ForceDescending {
		return nil, nil, fmt.Errorf("collation cannot at once be both ascending and descending")
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

	url, err := generateParameterizedUrl(mg, eventsEndpoint, payload)
	if err != nil {
		return nil, nil, err
	}
	r := simplehttp.NewHTTPRequest(url)
	r.SetBasicAuth(basicAuthUser, mg.ApiKey())
	var response map[string]interface{}
	err = r.GetResponseFromJSON(&response)
	if err != nil {
		return nil, nil, err
	}

	items := response["items"].([]interface{})
	events := make([]Event, len(items))
	for i, item := range items {
		events[i] = item.(map[string]interface{})
	}

	pagings := response["paging"].(map[string]interface{})
	var links = make(Links, len(pagings))
	for key, page := range pagings {
		links[key] = page.(string)
	}

	return events, links, err
}
