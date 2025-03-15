package mailgun_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStorage(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	var ctx = context.Background()

	m := mailgun.NewMessage(testDomain, "root@"+testDomain, "Subject", "Text Body", "stored@"+testDomain)
	resp, err := mg.Send(ctx, m)
	require.NoError(t, err)

	t.Logf("New Email: %s ID: %s\n", resp.Message, resp.ID)

	url, err := findStoredMessageURL(mg, strings.Trim(resp.ID, "<>"))
	require.NoError(t, err)

	stored, err := mg.GetStoredMessage(ctx, url)
	require.NoError(t, err)

	assert.Equal(t, "Subject", stored.Subject)
	assert.Equal(t, "root@"+testDomain, stored.From)
	assert.Equal(t, "stored@"+testDomain, stored.Recipients)

	_, err = mg.ReSend(ctx, url, "resend@"+testDomain)
	require.NoError(t, err)
}

// Tries to locate the first stored event type, returning the associated stored message key.
func findStoredMessageURL(mg mailgun.Mailgun, id string) (string, error) {
	it := mg.ListEvents(testDomain, nil)

	var page []events.Event
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

	return "", fmt.Errorf("no stored messages found; try changing MG_EMAIL_TO to an address that stores messages and try again")
}
