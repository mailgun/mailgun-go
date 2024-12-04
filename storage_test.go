package mailgun_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mailgun/mailgun-go/v4/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	var ctx = context.Background()

	m := mailgun.NewMessage("root@"+testDomain, "Subject", "Text Body", "stored@"+testDomain)
	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)

	t.Logf("New Email: %s Id: %s\n", msg, id)

	url, err := findStoredMessageURL(mg, strings.Trim(id, "<>"))
	require.NoError(t, err)

	resp, err := mg.GetStoredMessage(ctx, url)
	require.NoError(t, err)

	assert.Equal(t, "Subject", resp.Subject)
	assert.Equal(t, "root@"+testDomain, resp.From)
	assert.Equal(t, "stored@"+testDomain, resp.Recipients)

	_, _, err = mg.ReSend(ctx, url, "resend@"+testDomain)
	require.NoError(t, err)
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

	return "", fmt.Errorf("no stored messages found. Try changing MG_EMAIL_TO to an address that stores messages and try again.")
}
