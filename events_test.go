package mailgun_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go"
	"github.com/mailgun/mailgun-go/events"
)

func TestEventIteratorGetNext(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	it := mg.ListEvents(&mailgun.EventsOptions{Limit: 5})

	var firstPage, secondPage, previousPage []mailgun.Event
	var ctx = context.Background()

	ensure.True(t, it.Next(ctx, &firstPage))
	ensure.True(t, it.Paging.Next != "")
	ensure.True(t, len(firstPage) != 0)
	firstIterator := *it

	ensure.True(t, it.Next(ctx, &secondPage))
	ensure.True(t, len(secondPage) != 0)

	// Pages should be different
	ensure.NotDeepEqual(t, firstPage, secondPage)
	ensure.True(t, firstIterator.Paging.Next != it.Paging.Next)
	ensure.True(t, firstIterator.Paging.Previous != it.Paging.Previous)
	ensure.Nil(t, it.Err())

	// Previous()
	ensure.True(t, it.First(ctx, &firstPage))
	ensure.True(t, it.Next(ctx, &secondPage))

	ensure.True(t, it.Previous(ctx, &previousPage))
	ensure.True(t, len(previousPage) != 0)
	ensure.DeepEqual(t, previousPage[0].GetID(), firstPage[0].GetID())

	// First()
	ensure.True(t, it.First(ctx, &firstPage))
	ensure.True(t, len(firstPage) != 0)

	// Calling first resets the iterator to the first page
	ensure.True(t, it.Next(ctx, &secondPage))
	ensure.NotDeepEqual(t, firstPage, secondPage)

	// Last()
	var lastPage []mailgun.Event
	ensure.True(t, it.Next(ctx, &firstPage))
	ensure.True(t, len(firstPage) != 0)

	// Calling Last() is invalid unless you first use First() or Next()
	ensure.True(t, it.Last(ctx, &lastPage))
	ensure.True(t, len(lastPage) != 0)
}

func TestEventPoller(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	begin := time.Now().Add(time.Second * -3)

	// Very short poll interval
	it := mg.PollEvents(&mailgun.EventsOptions{
		// Only events with a timestamp after this date/time will be returned
		Begin: &begin,
		// How often we poll the api for new events
		PollInterval: time.Second * 4})

	eventChan := make(chan mailgun.Event, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		// Poll until our email event arrives
		var page []mailgun.Event
		for it.Poll(ctx, &page) {
			for _, e := range page {
				eventChan <- e
			}
		}
		close(eventChan)
	}()

	// Send an email
	m := mg.NewMessage("root@"+testDomain, "Subject", "Text Body", "user@"+testDomain)
	msg, id, err := mg.Send(m)
	ensure.Nil(t, err)

	t.Logf("New Email: %s Id: %s\n", msg, id)

	var accepted *events.Accepted
	for e := range eventChan {
		switch event := e.(type) {
		case *events.Accepted:
			t.Logf("Accepted Event: %s - %v", event.Message.Headers.MessageID, event.GetTimestamp())
			// If we find our accepted email event
			if id == ("<" + event.Message.Headers.MessageID + ">") {
				accepted = event
				cancel()
			}
		}
	}
	// Ensure we found our email
	ensure.Nil(t, it.Err())
	ensure.True(t, accepted != nil)
	ensure.DeepEqual(t, accepted.Recipient, "user@"+testDomain)
}

func ExampleListEvents() {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	it := mg.ListEvents(&mailgun.EventsOptions{Limit: 100})

	var page []mailgun.Event

	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// For each page of 100 events
	for it.Next(ctx, &page) {
		for _, e := range page {
			// You can access some fields via the interface
			//fmt.Printf("Event: '%s' TimeStamp: '%s'\n", e.GetName(), e.GetTimestamp())

			// and you can act upon each event by type
			switch event := e.(type) {
			case *events.Accepted:
				fmt.Printf("Accepted: auth: %t\n", event.Flags.IsAuthenticated)
			case *events.Delivered:
				fmt.Printf("Delivered transport: %s\n", event.Envelope.Transport)
			case *events.Failed:
				fmt.Printf("Failed reason: %s\n", event.Reason)
			case *events.Clicked:
				fmt.Printf("Clicked GeoLocation: %s\n", event.GeoLocation.Country)
			case *events.Opened:
				fmt.Printf("Opened GeoLocation: %s\n", event.GeoLocation.Country)
			case *events.Rejected:
				fmt.Printf("Rejected reason: %s\n", event.Reject.Reason)
			case *events.Stored:
				fmt.Printf("Stored URL: %s\n", event.Storage.URL)
			case *events.Unsubscribed:
				fmt.Printf("Unsubscribed client OS: %s\n", event.ClientInfo.ClientOS)
			}
		}
	}
	// Output:
	// Accepted: auth: false
	// Accepted: auth: true
	// Delivered transport: smtp
	// Delivered transport: http
	// Stored URL: http://mailgun.text/some/url
	// Clicked GeoLocation: US
	// Clicked GeoLocation: US
	// Clicked GeoLocation: US
	// Opened GeoLocation: US
	// Opened GeoLocation: US
	// Opened GeoLocation: US
	// Unsubscribed client OS: OS X
	// Unsubscribed client OS: OS X
}
