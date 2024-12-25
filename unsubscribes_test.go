package mailgun_test

import (
	"context"
	"os"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

func TestCreateUnsubscriber(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL3())

	email := randomEmail("unsubcribe", os.Getenv("MG_DOMAIN"))
	ctx := context.Background()

	// Create unsubscription record
	require.NoError(t, mg.CreateUnsubscribe(ctx, testDomain, email, "*"))
}

func TestCreateUnsubscribes(t *testing.T) {
	unsubscribes := []mailgun.Unsubscribe{
		{
			Address: randomEmail("unsubcribe", os.Getenv("MG_DOMAIN")),
		},
		{
			Address: randomEmail("unsubcribe", os.Getenv("MG_DOMAIN")),
			Tags:    []string{"tag1"},
		},
	}
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL3())
	ctx := context.Background()

	// Create unsubscription records
	require.NoError(t, mg.CreateUnsubscribes(ctx, testDomain, unsubscribes))
}

func TestListUnsubscribes(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL3())
	ctx := context.Background()

	it := mg.ListUnsubscribes(testDomain, nil)
	var page []mailgun.Unsubscribe
	for it.Next(ctx, &page) {
		t.Logf("Received %d unsubscribe records.\n", len(page))
		if len(page) > 0 {
			t.Log("ID\tAddress\tCreated At\tTags\t")
			for _, u := range page {
				t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tags)
			}
		}
	}
	require.NoError(t, it.Err())
}

func TestGetUnsubscribe(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL3())

	email := randomEmail("unsubcribe", os.Getenv("MG_DOMAIN"))

	ctx := context.Background()

	// Create unsubscription record
	require.NoError(t, mg.CreateUnsubscribe(ctx, testDomain, email, "*"))

	u, err := mg.GetUnsubscribe(ctx, testDomain, email)
	require.NoError(t, err)
	t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tags)

	// Destroy the unsubscription record
	require.NoError(t, mg.DeleteUnsubscribe(ctx, testDomain, email))
}

func TestCreateDestroyUnsubscription(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL3())

	email := randomEmail("unsubcribe", os.Getenv("MG_DOMAIN"))

	ctx := context.Background()

	// Create unsubscription record
	require.NoError(t, mg.CreateUnsubscribe(ctx, testDomain, email, "*"))

	_, err := mg.GetUnsubscribe(ctx, testDomain, email)
	require.NoError(t, err)

	// Destroy the unsubscription record
	require.NoError(t, mg.DeleteUnsubscribe(ctx, testDomain, email))
}
