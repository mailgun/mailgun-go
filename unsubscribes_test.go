package mailgun_test

import (
	"context"
	"os"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

func TestCreateUnsubscriber(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	email := randomEmail("unsubcribe", os.Getenv("MG_DOMAIN"))
	ctx := context.Background()

	// Create unsubscription record
	require.NoError(t, mg.CreateUnsubscribe(ctx, email, "*"))
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
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	// Create unsubscription records
	require.NoError(t, mg.CreateUnsubscribes(ctx, unsubscribes))
}

func TestListUnsubscribes(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	it := mg.ListUnsubscribes(nil)
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
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	email := randomEmail("unsubcribe", os.Getenv("MG_DOMAIN"))

	ctx := context.Background()

	// Create unsubscription record
	require.NoError(t, mg.CreateUnsubscribe(ctx, email, "*"))

	u, err := mg.GetUnsubscribe(ctx, email)
	require.NoError(t, err)
	t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tags)

	// Destroy the unsubscription record
	require.NoError(t, mg.DeleteUnsubscribe(ctx, email))
}

func TestCreateDestroyUnsubscription(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	email := randomEmail("unsubcribe", os.Getenv("MG_DOMAIN"))

	ctx := context.Background()

	// Create unsubscription record
	require.NoError(t, mg.CreateUnsubscribe(ctx, email, "*"))

	_, err := mg.GetUnsubscribe(ctx, email)
	require.NoError(t, err)
	/*t.Logf("Received %d out of %d unsubscribe records.\n", len(us), n)*/

	// Destroy the unsubscription record
	require.NoError(t, mg.DeleteUnsubscribe(ctx, email))
}
