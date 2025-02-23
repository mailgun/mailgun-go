package mailgun_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mailgun/mailgun-go/v4/mtypes"
	"github.com/stretchr/testify/require"
)

func TestGetBounces(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	it := mg.ListBounces(testDomain, nil)

	var page []mtypes.Bounce
	for it.Next(ctx, &page) {
		for _, bounce := range page {
			t.Logf("Bounce: %+v\n", bounce)
		}
	}
	require.NoError(t, it.Err())
}

func TestGetSingleBounce(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	exampleEmail := fmt.Sprintf("%s@%s", strings.ToLower(randomString(64, "")),
		os.Getenv("MG_DOMAIN"))
	_, err = mg.GetBounce(ctx, testDomain, exampleEmail)
	require.NotNil(t, err)

	var ure *mailgun.UnexpectedResponseError
	require.ErrorAs(t, err, &ure)
	require.Equal(t, http.StatusNotFound, ure.Actual)
}

func TestAddDelBounces(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	findBounce := func(address string) bool {
		it := mg.ListBounces(testDomain, nil)
		var page []mtypes.Bounce
		for it.Next(ctx, &page) {
			require.True(t, len(page) != 0)
			for _, bounce := range page {
				t.Logf("Bounce Address: %s\n", bounce.Address)
				if bounce.Address == address {
					return true
				}
			}
		}
		if it.Err() != nil {
			t.Logf("BounceIterator err: %s", it.Err())
		}
		return false
	}

	// Compute an e-mail address for our Bounce.
	exampleEmail := fmt.Sprintf("%s@%s", strings.ToLower(randomString(8, "bounce")), domain)

	// Add the bounce for our address.
	err = mg.AddBounce(ctx, testDomain, exampleEmail, "550", "TestAddDelBounces-generated error")
	require.NoError(t, err)

	// Give API some time to refresh cache
	time.Sleep(time.Second)

	// We should now have one bounce listed when we query the API.
	if !findBounce(exampleEmail) {
		t.Fatalf("Expected bounce for address %s in list of bounces", exampleEmail)
	}

	bounce, err := mg.GetBounce(ctx, testDomain, exampleEmail)
	require.NoError(t, err)
	if bounce.Address != exampleEmail {
		t.Fatalf("Expected at least one bounce for %s", exampleEmail)
	}
	t.Logf("Bounce Created At: %s", bounce.CreatedAt)

	// Delete it.  This should put us back the way we were.
	err = mg.DeleteBounce(ctx, testDomain, exampleEmail)
	require.NoError(t, err)

	// Make sure we're back to the way we were.
	if findBounce(exampleEmail) {
		t.Fatalf("Un-expected bounce for address %s in list of bounces", exampleEmail)
	}

	_, err = mg.GetBounce(ctx, testDomain, exampleEmail)
	require.NotNil(t, err)
}

func TestAddDelBounceList(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	findBounce := func(address string) bool {
		it := mg.ListBounces(testDomain, nil)
		var page []mtypes.Bounce
		for it.Next(ctx, &page) {
			require.True(t, len(page) != 0)
			for _, bounce := range page {
				t.Logf("Bounce Address: %s\n", bounce.Address)
				if bounce.Address == address {
					return true
				}
			}
		}
		if it.Err() != nil {
			t.Logf("BounceIterator err: %s", it.Err())
		}
		return false
	}

	createdAt, err := mtypes.NewRFC2822Time("Thu, 13 Oct 2011 18:02:00 +0000")
	if err != nil {
		t.Fatalf("invalid time")
	}

	// Generate a list of bounces
	bounces := []mtypes.Bounce{
		{
			Code:    "550",
			Address: fmt.Sprintf("%s@%s", strings.ToLower(randomString(8, "bounce")), domain),
			Error:   "TestAddDelBounces-generated error",
		},
		{
			Code:      "550",
			Address:   fmt.Sprintf("%s@%s", strings.ToLower(randomString(8, "bounce")), domain),
			Error:     "TestAddDelBounces-generated error",
			CreatedAt: createdAt,
		},
	}

	// Add the bounce for our address.
	err = mg.AddBounces(ctx, testDomain, bounces)
	require.NoError(t, err)

	for _, expect := range bounces {
		if !findBounce(expect.Address) {
			t.Fatalf("Expected bounce for address %s in list of bounces", expect.Address)
		}

		bounce, err := mg.GetBounce(ctx, testDomain, expect.Address)
		require.NoError(t, err)
		if bounce.Address != expect.Address {
			t.Fatalf("Expected at least one bounce for %s", expect.Address)
		}
		t.Logf("Bounce Created At: %s", bounce.CreatedAt)
		if !expect.CreatedAt.IsZero() && !time.Time(bounce.CreatedAt).Equal(time.Time(expect.CreatedAt)) {
			t.Fatalf("Expected bounce createdAt to be %s, got %s", expect.CreatedAt, bounce.CreatedAt)
		}
	}

	// Delete the bounce list.  This should put us back the way we were.
	err = mg.DeleteBounceList(ctx, testDomain)
	require.NoError(t, err)

	it := mg.ListBounces(testDomain, nil)
	var page []mtypes.Bounce
	if it.Next(ctx, &page) {
		t.Fatalf("Expected no item in the bounce list")
	}
}
