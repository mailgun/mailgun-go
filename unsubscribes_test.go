package mailgun

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestCreateUnsubscriber(t *testing.T) {
	email := randomEmail("unsubcribe", reqEnv(t, "MG_DOMAIN"))
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	// Create unsubscription record
	ensure.Nil(t, mg.Unsubscribe(ctx, email, "*"))
}

func TestGetUnsubscribes(t *testing.T) {
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	n, us, err := mg.ListUnsubscribes(ctx, DefaultLimit, DefaultSkip)
	ensure.Nil(t, err)

	t.Logf("Received %d out of %d unsubscribe records.\n", len(us), n)
	if len(us) > 0 {
		t.Log("ID\tAddress\tCreated At\tTags\t")
		for _, u := range us {
			t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tags)
		}
	}
}

func TestGetUnsubscriptionByAddress(t *testing.T) {
	email := randomEmail("unsubcribe", reqEnv(t, "MG_DOMAIN"))
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	// Create unsubscription record
	ensure.Nil(t, mg.Unsubscribe(ctx, email, "*"))

	n, us, err := mg.GetUnsubscribes(ctx, email)
	ensure.Nil(t, err)

	t.Logf("Received %d out of %d unsubscribe records.\n", len(us), n)
	if len(us) > 0 {
		t.Log("ID\tAddress\tCreated At\tTags\t")
		for _, u := range us {
			t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tags)
		}
	}
	// Destroy the unsubscription record
	ensure.Nil(t, mg.RemoveUnsubscribe(ctx, email))
}

func TestCreateDestroyUnsubscription(t *testing.T) {
	email := randomEmail("unsubcribe", reqEnv(t, "MG_DOMAIN"))
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	ctx := context.Background()

	// Create unsubscription record
	ensure.Nil(t, mg.Unsubscribe(ctx, email, "*"))

	n, us, err := mg.GetUnsubscribes(ctx, email)
	ensure.Nil(t, err)
	t.Logf("Received %d out of %d unsubscribe records.\n", len(us), n)

	// Destroy the unsubscription record
	ensure.Nil(t, mg.RemoveUnsubscribe(ctx, email))
}
