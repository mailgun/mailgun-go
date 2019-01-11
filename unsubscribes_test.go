package mailgun

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestCreateUnsubscriber(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	email := randomEmail("unsubcribe", reqEnv(t, "MG_DOMAIN"))
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	// Create unsubscription record
	ensure.Nil(t, mg.CreateUnsubscribe(ctx, email, "*"))
}

func TestGetUnsubscribes(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	it := mg.ListUnsubscribes(nil)
	var page []Unsubscribe
	for it.Next(ctx, &page) {
		t.Logf("Received %d unsubscribe records.\n", len(page))
		if len(page) > 0 {
			t.Log("ID\tAddress\tCreated At\tTags\t")
			for _, u := range page {
				t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tags)
			}
		}
	}
	ensure.Nil(t, it.Err())
}

func TestGetUnsubscribe(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	email := randomEmail("unsubcribe", reqEnv(t, "MG_DOMAIN"))
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	// Create unsubscription record
	ensure.Nil(t, mg.CreateUnsubscribe(ctx, email, "*"))

	u, err := mg.GetUnsubscribe(ctx, email)
	ensure.Nil(t, err)
	t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tags)

	/*t.Logf("Received %d unsubscribe records.\n", len(us))
	if len(us) > 0 {
		t.Log("ID\tAddress\tCreated At\tTags\t")
		for _, u := range us {
			t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tags)
		}
	}
	*/
	// Destroy the unsubscription record
	ensure.Nil(t, mg.DeleteUnsubscribe(ctx, email))
}

func TestCreateDestroyUnsubscription(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	email := randomEmail("unsubcribe", reqEnv(t, "MG_DOMAIN"))
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	ctx := context.Background()

	// Create unsubscription record
	ensure.Nil(t, mg.CreateUnsubscribe(ctx, email, "*"))

	_, err = mg.GetUnsubscribe(ctx, email)
	ensure.Nil(t, err)
	/*t.Logf("Received %d out of %d unsubscribe records.\n", len(us), n)*/

	// Destroy the unsubscription record
	ensure.Nil(t, mg.DeleteUnsubscribe(ctx, email))
}
