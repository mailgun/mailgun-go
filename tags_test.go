package mailgun

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/facebookgo/ensure"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var _ = Describe("/v3/{domain}/tags", func() {
	log := log.New(GinkgoWriter, "tags_test - ", 0)
	var t GinkgoTInterface
	var mg Mailgun
	var err error

	BeforeSuite(func() {
		mg, err = NewMailgunFromEnv()
		msg := mg.NewMessage(fromUser, exampleSubject, exampleText, reqEnv(t, "MG_EMAIL_TO"))
		msg.AddTag("newsletter")
		msg.AddTag("homer")
		msg.AddTag("bart")
		msg.AddTag("disco-steve")
		msg.AddTag("newsletter")

		ctx := context.Background()
		// Create an email with some tags attached
		_, _, err := mg.Send(ctx, msg)
		if err != nil {
			Fail(fmt.Sprintf("Mesage send: '%s'", err.Error()))
		}
		// Wait for the tag to show up
		if err := waitForTag(mg, "newsletter"); err != nil {
			Fail(fmt.Sprintf("While waiting for message: '%s'", err.Error()))
		}
	})

	BeforeEach(func() {
		t = GinkgoT()
		mg, err = NewMailgunFromEnv()
		ensure.Nil(t, err)
	})

	Describe("ListTags()", func() {
		Context("When a limit parameter of -1 is supplied", func() {
			It("Should return a list of available tags", func() {
				ctx := context.Background()
				it := mg.ListTags(nil)
				var page TagsPage
				for it.Next(ctx, &page) {
					Expect(len(page.Items)).NotTo(Equal(0))
					log.Printf("Tags: %+v\n", page)
				}
				ensure.Nil(t, it.Err())
			})
		})
		Context("When limit parameter is supplied", func() {
			It("Should return a limited list of available tags", func() {
				ctx := context.Background()
				cursor := mg.ListTags(&TagOptions{Limit: 1})

				var tags TagsPage
				for cursor.Next(ctx, &tags) {
					ensure.DeepEqual(t, len(tags.Items), 1)
					log.Printf("Tags: %+v\n", tags.Items)
				}
				ensure.Nil(t, cursor.Err())
			})
		})
	})

	Describe("DeleteTag()", func() {
		Context("When deleting an existing tag", func() {
			It("Should not error", func() {
				ctx := context.Background()
				err = mg.DeleteTag(ctx, "newsletter")
				ensure.Nil(t, err)
			})
		})
	})

	Describe("GetTag()", func() {
		Context("When requesting an existing tag", func() {
			It("Should not error", func() {
				ctx := context.Background()
				tag, err := mg.GetTag(ctx, "homer")
				ensure.Nil(t, err)
				ensure.DeepEqual(t, tag.Value, "homer")
			})
		})
		Context("When requesting an non-existant tag", func() {
			It("Should return error", func() {
				ctx := context.Background()
				_, err := mg.GetTag(ctx, "i-dont-exist")
				ensure.NotNil(t, err)
				ensure.DeepEqual(t, GetStatusFromErr(err), 404)
			})
		})
	})
})

func waitForTag(mg Mailgun, tag string) error {
	ctx := context.Background()
	var attempts int
	for attempts <= 5 {
		_, err := mg.GetTag(ctx, tag)
		if err != nil {
			if GetStatusFromErr(err) == 404 {
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
