package mailgun_test

import (
	"context"
	"fmt"
	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/mailgun/mailgun-go/v3/events"
	"strings"
	"testing"
)

func TestStorage(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	var ctx = context.Background()

	m := mg.NewMessage("root@"+testDomain, "Subject", "Text Body", "stored@"+testDomain)
	msg, id, err := mg.Send(ctx, m)
	ensure.Nil(t, err)

	t.Logf("New Email: %s Id: %s\n", msg, id)

	url, err := findStoredMessageURL(mg, strings.Trim(id, "<>"))

	resp, err := mg.GetStoredMessage(ctx, url)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, "Subject", resp.Subject)
	ensure.DeepEqual(t, "root@"+testDomain, resp.From)
	ensure.DeepEqual(t, "stored@"+testDomain, resp.Recipients)

	_, _, err = mg.ReSend(ctx, url, "resend@"+testDomain)
	ensure.Nil(t, err)
}

// Tries to locate the first stored event type, returning the associated stored message key.
func findStoredMessageURL(mg mailgun.Mailgun, id string) (string, error) {
	it := mg.ListEvents(nil)

	var page []mailgun.Event
	for it.Next(context.Background(), &page) {
		for _, event := range page {
			if event.GetName() == events.EventStored && event.GetID() == id {
				return event.(*events.Stored).Storage.URL, nil
			}
		}
	}
	if it.Err() != nil {
		return "", it.Err()
	}
	return "", fmt.Errorf("No stored messages found.  Try changing MG_EMAIL_TO to an address that stores messages and try again.")
}
