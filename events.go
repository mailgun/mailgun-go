package mailgun

import (
	"context"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/mailgun/mailgun-go/v4/events"
)

// ListEventOptions{} modifies the behavior of ListEvents()
type ListEventOptions struct {
	// Limits the results to a specific start and end time
	Begin, End time.Time
	// ForceAscending and ForceDescending are used to force Mailgun to use a given
	// traversal order of the events. If both ForceAscending and ForceDescending are
	// true, an error will result. If none, the default will be inferred from the Begin
	// and End parameters.
	ForceAscending, ForceDescending bool
	// Compact, if true, compacts the returned JSON to minimize transmission bandwidth.
	Compact bool
	// Limit caps the number of results returned.  If left unspecified, MailGun assumes 100.
	Limit int
	// Filter allows the caller to provide more specialized filters on the query.
	// Consult the Mailgun documentation for more details.
	Filter       map[string]string
	PollInterval time.Duration
}

// EventIterator maintains the state necessary for paging though small parcels of a larger set of events.
type EventIterator struct {
	events.Response
	mg  Mailgun
	err error
}

// Create an new iterator to fetch a page of events from the events api with a specific domain
func (mg *MailgunImpl) ListEventsWithDomain(opts *ListEventOptions, domain string) *EventIterator {
	url := generateApiUrlWithDomain(mg, eventsEndpoint, domain)
	return mg.listEvents(url, opts)
}

// Create an new iterator to fetch a page of events from the events api
func (mg *MailgunImpl) ListEvents(opts *ListEventOptions) *EventIterator {
	url := generateApiUrl(mg, eventsEndpoint)
	return mg.listEvents(url, opts)
}

func (mg *MailgunImpl) listEvents(url string, opts *ListEventOptions) *EventIterator {
	req := newHTTPRequest(url)
	if opts != nil {
		if opts.Limit > 0 {
			req.addParameter("limit", fmt.Sprintf("%d", opts.Limit))
		}
		if opts.Compact {
			req.addParameter("pretty", "no")
		}
		if opts.ForceAscending {
			req.addParameter("ascending", "yes")
		} else if opts.ForceDescending {
			req.addParameter("ascending", "no")
		}
		if !opts.Begin.IsZero() {
			req.addParameter("begin", formatMailgunTime(opts.Begin))
		}
		if !opts.End.IsZero() {
			req.addParameter("end", formatMailgunTime(opts.End))
		}
		if opts.Filter != nil {
			for k, v := range opts.Filter {
				req.addParameter(k, v)
			}
		}
	}
	url, err := req.generateUrlWithParameters()
	return &EventIterator{
		mg:       mg,
		Response: events.Response{Paging: events.Paging{Next: url, First: url}},
		err:      err,
	}
}

// If an error occurred during iteration `Err()` will return non nil
func (ei *EventIterator) Err() error {
	return ei.err
}

// Next retrieves the next page of events from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error
func (ei *EventIterator) Next(ctx context.Context, events *[]Event) bool {
	if ei.err != nil {
		return false
	}
	ei.err = ei.fetch(ctx, ei.Paging.Next)
	if ei.err != nil {
		return false
	}
	*events, ei.err = ParseEvents(ei.Items)
	if ei.err != nil {
		return false
	}
	if len(ei.Items) == 0 {
		return false
	}
	return true
}

// First retrieves the first page of events from the api. Returns false if there
// was an error. It also sets the iterator object to the first page.
// Use `.Err()` to retrieve the error.
func (ei *EventIterator) First(ctx context.Context, events *[]Event) bool {
	if ei.err != nil {
		return false
	}
	ei.err = ei.fetch(ctx, ei.Paging.First)
	if ei.err != nil {
		return false
	}
	*events, ei.err = ParseEvents(ei.Items)
	return true
}

// Last retrieves the last page of events from the api.
// Calling Last() is invalid unless you first call First() or Next()
// Returns false if there was an error. It also sets the iterator object
// to the last page. Use `.Err()` to retrieve the error.
func (ei *EventIterator) Last(ctx context.Context, events *[]Event) bool {
	if ei.err != nil {
		return false
	}
	ei.err = ei.fetch(ctx, ei.Paging.Last)
	if ei.err != nil {
		return false
	}
	*events, ei.err = ParseEvents(ei.Items)
	return true
}

// Previous retrieves the previous page of events from the api. Returns false when there
// no more pages to retrieve or if there was an error. Use `.Err()` to retrieve
// the error if any
func (ei *EventIterator) Previous(ctx context.Context, events *[]Event) bool {
	if ei.err != nil {
		return false
	}
	if ei.Paging.Previous == "" {
		return false
	}
	ei.err = ei.fetch(ctx, ei.Paging.Previous)
	if ei.err != nil {
		return false
	}
	*events, ei.err = ParseEvents(ei.Items)
	if len(ei.Items) == 0 {
		return false
	}
	return true
}

func (ei *EventIterator) fetch(ctx context.Context, url string) error {
	ei.Items = nil
	r := newHTTPRequest(url)
	r.setClient(ei.mg.Client())
	r.setBasicAuth(basicAuthUser, ei.mg.APIKey())

	resp, err := makeRequest(ctx, r, "GET", nil)
	if err != nil {
		return err
	}

	if err := jsoniter.Unmarshal(resp.Data, &ei.Response); err != nil {
		return fmt.Errorf("failed to un-marshall event.Response: %s", err)
	}
	return nil
}

// EventPoller maintains the state necessary for polling events
type EventPoller struct {
	it            *EventIterator
	opts          ListEventOptions
	thresholdTime time.Time
	beginTime     time.Time
	sleepUntil    time.Time
	mg            Mailgun
	err           error
}

// Poll the events api and return new events as they occur
//  it = mg.PollEvents(&ListEventOptions{
//    // Only events with a timestamp after this date/time will be returned
//    Begin:        time.Now().Add(time.Second * -3),
//    // How often we poll the api for new events
//    PollInterval: time.Second * 4
//  })
//
//  var events []Event
//  ctx, cancel := context.WithCancel(context.Background())
//
//  // Blocks until new events appear or context is cancelled
//  for it.Poll(ctx, &events) {
//    for _, event := range(events) {
//      fmt.Printf("Event %+v\n", event)
//    }
//  }
//  if it.Err() != nil {
//    log.Fatal(it.Err())
//  }
func (mg *MailgunImpl) PollEvents(opts *ListEventOptions) *EventPoller {
	now := time.Now()
	// ForceAscending must be set
	opts.ForceAscending = true

	// Default begin time is 30 minutes ago
	if opts.Begin.IsZero() {
		opts.Begin = now.Add(time.Minute * -30)
	}

	// Set a 15 second poll interval if none set
	if opts.PollInterval.Nanoseconds() == 0 {
		opts.PollInterval = time.Duration(time.Second * 15)
	}

	return &EventPoller{
		it:   mg.ListEvents(opts),
		opts: *opts,
		mg:   mg,
	}
}

// If an error occurred during polling `Err()` will return non nil
func (ep *EventPoller) Err() error {
	return ep.err
}

func (ep *EventPoller) Poll(ctx context.Context, events *[]Event) bool {
	var currentPage string
	var results []Event

	if ep.opts.Begin.IsZero() {
		ep.beginTime = time.Now().UTC()
	}

	for {
		// Remember our current page url
		currentPage = ep.it.Paging.Next

		// Attempt to get a page of events
		var page []Event
		if ep.it.Next(ctx, &page) == false {
			if ep.it.Err() == nil && len(page) == 0 {
				// No events, sleep for our poll interval
				goto SLEEP
			}
			ep.err = ep.it.Err()
			return false
		}

		for _, e := range page {
			// If any events on the page are older than our being time
			if e.GetTimestamp().After(ep.beginTime) {
				results = append(results, e)
			}
		}

		// If we have events to return
		if len(results) != 0 {
			*events = results
			results = nil
			return true
		}

	SLEEP:
		// Since we didn't find an event older than our
		// threshold, fetch this same page again
		ep.it.Paging.Next = currentPage

		// Sleep the rest of our duration
		tick := time.NewTicker(ep.opts.PollInterval)
		select {
		case <-ctx.Done():
			return false
		case <-tick.C:
			tick.Stop()
		}
	}

}

// Given time.Time{} return a float64 as given in mailgun event timestamps
func TimeToFloat(t time.Time) float64 {
	return float64(t.Unix()) + (float64(t.Nanosecond()/int(time.Microsecond)) / float64(1000000))
}
