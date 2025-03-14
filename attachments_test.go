package mailgun_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createAttachment(t *testing.T) string {
	t.Helper()
	name := "/tmp/" + randomString(10, "attachment-")
	f, err := os.Create(name)
	require.NoError(t, err)

	_, err = f.WriteString(randomString(100, ""))
	require.NoError(t, err)
	require.Nil(t, f.Close())
	return name
}

func TestMultipleAttachments(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	var ctx = context.Background()

	m := mailgun.NewMessage(testDomain, "root@"+testDomain, "Subject", "Text Body", "attachment@"+testDomain)

	// Add 2 attachments
	m.AddAttachment(createAttachment(t))
	m.AddAttachment(createAttachment(t))

	resp, err := mg.Send(ctx, m)
	require.NoError(t, err)

	id := strings.Trim(resp.ID, "<>")
	t.Logf("New Email: %s ID: %s\n", resp.Message, id)

	e, err := findAcceptedMessage(mg, id)
	require.NoError(t, err)
	require.NotNil(t, e)

	assert.Equal(t, e.ID, id)
	assert.Len(t, e.Message.Attachments, 2)
	for _, f := range e.Message.Attachments {
		t.Logf("attachment: %v\n", f)
		assert.Equal(t, 100, f.Size)
	}
}

func findAcceptedMessage(mg mailgun.Mailgun, id string) (*events.Accepted, error) {
	it := mg.ListEvents(testDomain, nil)

	var page []events.Event
	for it.Next(context.Background(), &page) {
		for _, event := range page {
			if event.GetName() == events.EventAccepted && event.GetID() == id {
				e, ok := event.(*events.Accepted)
				if !ok {
					return nil, fmt.Errorf("unexpected event type: %T", event)
				}

				return e, nil
			}
		}
	}
	if it.Err() != nil {
		return nil, it.Err()
	}
	return nil, fmt.Errorf("no accepted messages found for '%s'", id)
}
