package mailgun_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mailgun/mailgun-go/v4/events"
	"github.com/stretchr/testify/require"
)

func createAttachment(t *testing.T) string {
	t.Helper()
	name := "/tmp/" + randomString(10, "attachment-")
	f, err := os.Create(name)
	require.NoError(t, err)

	_, err = f.Write([]byte(randomString(100, "")))
	require.NoError(t, err)
	require.Nil(t, f.Close())
	return name
}

func TestMultipleAttachments(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	var ctx = context.Background()

	m := mailgun.NewMessage("root@"+testDomain, "Subject", "Text Body", "attachment@"+testDomain)

	// Add 2 attachments
	m.AddAttachment(createAttachment(t))
	m.AddAttachment(createAttachment(t))

	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)

	id = strings.Trim(id, "<>")
	t.Logf("New Email: %s Id: %s\n", msg, id)

	e, err := findAcceptedMessage(mg, id)
	require.NotNil(t, e)

	require.Equal(t, e.ID, id)
	require.Len(t, e.Message.Attachments, 2)
	for _, f := range e.Message.Attachments {
		t.Logf("attachment: %v\n", f)
		require.Equal(t, 100, f.Size)
	}
}

func findAcceptedMessage(mg mailgun.Mailgun, id string) (*events.Accepted, error) {
	it := mg.ListEvents(nil)

	var page []mailgun.Event
	for it.Next(context.Background(), &page) {
		for _, event := range page {
			if event.GetName() == events.EventAccepted && event.GetID() == id {
				return event.(*events.Accepted), nil
			}
		}
	}
	if it.Err() != nil {
		return nil, it.Err()
	}
	return nil, fmt.Errorf("no accepted messages found for '%s'", id)
}
