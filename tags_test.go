package mailgun_test

import (
	"context"
	"testing"
	"time"

	"github.com/mailgun/errors"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	fromUser       = "=?utf-8?q?Katie_Brewer=2C_CFP=C2=AE?= <joe@example.com>"
	exampleSubject = "Mailgun-go Example Subject"
	exampleText    = "Testing some Mailgun awesomeness!"
)

func TestTags(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	msg := mailgun.NewMessage(fromUser, exampleSubject, exampleText, "test@example.com")
	require.NoError(t, msg.AddTag("newsletter"))
	require.NoError(t, msg.AddTag("homer"))
	require.NoError(t, msg.AddTag("bart"))
	require.NotNil(t, msg.AddTag("disco-steve"))
	require.NotNil(t, msg.AddTag("newsletter"))

	ctx := context.Background()
	// Create an email with some tags attached
	_, _, err := mg.Send(ctx, msg)
	require.NoError(t, err)

	// Wait for the tag to show up
	require.NoError(t, waitForTag(mg, "newsletter"))

	// Should return a list of available tags
	it := mg.ListTags(nil)
	var page []mailgun.Tag
	for it.Next(ctx, &page) {
		require.True(t, len(page) != 0)
	}
	require.NoError(t, it.Err())

	// Should return a limited list of available tags
	cursor := mg.ListTags(&mailgun.ListTagOptions{Limit: 1})

	var tags []mailgun.Tag
	for cursor.Next(ctx, &tags) {
		require.Len(t, tags, 1)
	}
	require.NoError(t, cursor.Err())

	err = mg.DeleteTag(ctx, "newsletter")
	require.NoError(t, err)

	tag, err := mg.GetTag(ctx, "homer")
	require.NoError(t, err)
	assert.Equal(t, "homer", tag.Value)

	_, err = mg.GetTag(ctx, "i-dont-exist")
	require.NotNil(t, err)
	assert.Equal(t, 404, mailgun.GetStatusFromErr(err))
}

func waitForTag(mg mailgun.Mailgun, tag string) error {
	ctx := context.Background()
	var attempts int
	for attempts <= 5 {
		_, err := mg.GetTag(ctx, tag)
		if err != nil {
			if mailgun.GetStatusFromErr(err) == 404 {
				time.Sleep(time.Second * 2)
				attempts++
				continue
			}

			return err
		}

		return nil
	}

	return errors.Errorf("Waited to long for tag '%s' to show up", tag)
}

func TestDeleteTag(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	require.NoError(t, mg.DeleteTag(ctx, "newsletter"))
}
