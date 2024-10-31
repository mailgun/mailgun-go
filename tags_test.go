package mailgun_test

import (
	"context"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/errors"
	"github.com/mailgun/mailgun-go/v4"
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
	ensure.Nil(t, msg.AddTag("newsletter"))
	ensure.Nil(t, msg.AddTag("homer"))
	ensure.Nil(t, msg.AddTag("bart"))
	ensure.NotNil(t, msg.AddTag("disco-steve"))
	ensure.NotNil(t, msg.AddTag("newsletter"))

	ctx := context.Background()
	// Create an email with some tags attached
	_, _, err := mg.Send(ctx, msg)
	ensure.Nil(t, err)

	// Wait for the tag to show up
	ensure.Nil(t, waitForTag(mg, "newsletter"))

	// Should return a list of available tags
	it := mg.ListTags(nil)
	var page []mailgun.Tag
	for it.Next(ctx, &page) {
		ensure.True(t, len(page) != 0)
	}
	ensure.Nil(t, it.Err())

	// Should return a limited list of available tags
	cursor := mg.ListTags(&mailgun.ListTagOptions{Limit: 1})

	var tags []mailgun.Tag
	for cursor.Next(ctx, &tags) {
		ensure.DeepEqual(t, len(tags), 1)
	}
	ensure.Nil(t, cursor.Err())

	err = mg.DeleteTag(ctx, "newsletter")
	ensure.Nil(t, err)

	tag, err := mg.GetTag(ctx, "homer")
	ensure.Nil(t, err)
	ensure.DeepEqual(t, tag.Value, "homer")

	_, err = mg.GetTag(ctx, "i-dont-exist")
	ensure.NotNil(t, err)
	ensure.DeepEqual(t, mailgun.GetStatusFromErr(err), 404)

}

func waitForTag(mg mailgun.Mailgun, tag string) error {
	ctx := context.Background()
	var attempts int
	for attempts <= 5 {
		_, err := mg.GetTag(ctx, tag)
		if err != nil {
			if mailgun.GetStatusFromErr(err) == 404 {
				time.Sleep(time.Second * 2)
				attempts += 1
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

	ensure.Nil(t, mg.DeleteTag(ctx, "newsletter"))
}
