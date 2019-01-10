package mailgun_test

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/pkg/errors"
)

const (
	fromUser       = "=?utf-8?q?Katie_Brewer=2C_CFP=C2=AE?= <joe@example.com>"
	exampleSubject = "Mailgun-go Example Subject"
	exampleText    = "Testing some Mailgun awesomeness!"
)

func TestTags(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}
	mg, err := mailgun.NewMailgunFromEnv()
	msg := mg.NewMessage(fromUser, exampleSubject, exampleText, os.Getenv("MG_EMAIL_TO"))
	msg.AddTag("newsletter")
	msg.AddTag("homer")
	msg.AddTag("bart")
	msg.AddTag("disco-steve")
	msg.AddTag("newsletter")

	ctx := context.Background()
	// Create an email with some tags attached
	_, _, err = mg.Send(ctx, msg)
	ensure.Nil(t, err)

	// Wait for the tag to show up
	ensure.Nil(t,waitForTag(mg, "newsletter"))

	// Should return a list of available tags
	it := mg.ListTags( nil)
	var page []mailgun.TagItem
	for it.Next(ctx, &page) {
		ensure.True(t, len(page)!= 0)
		log.Printf("Tags: %+v\n", page)
	}
	ensure.Nil(t, it.Err())

	// Should return a limited list of available tags
	cursor := mg.ListTags(&mailgun.ListTagOptions{Limit: 1})

	var tags []mailgun.TagItem
	for cursor.Next(ctx, &tags) {
		ensure.DeepEqual(t, len(tags), 1)
		log.Printf("Tags: %+v\n", tags)
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
