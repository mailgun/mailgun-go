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

	"github.com/facebookgo/ensure"
)

func TestGetBounces(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	it := mg.ListBounces(nil)

	var page []mailgun.Bounce
	for it.Next(ctx, &page) {
		for _, bounce := range page {
			t.Logf("Bounce: %+v\n", bounce)
		}
	}
	ensure.Nil(t, it.Err())
}

func TestGetSingleBounce(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	exampleEmail := fmt.Sprintf("%s@%s", strings.ToLower(randomString(64, "")),
		os.Getenv("MG_DOMAIN"))
	_, err := mg.GetBounce(ctx, exampleEmail)
	ensure.NotNil(t, err)

	ure, ok := err.(*mailgun.UnexpectedResponseError)
	ensure.True(t, ok)
	ensure.DeepEqual(t, ure.Actual, http.StatusNotFound)
}

func TestAddDelBounces(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	findBounce := func(address string) bool {
		it := mg.ListBounces(nil)
		var page []mailgun.Bounce
		for it.Next(ctx, &page) {
			ensure.True(t, len(page) != 0)
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
	err := mg.AddBounce(ctx, exampleEmail, "550", "TestAddDelBounces-generated error")
	ensure.Nil(t, err)

	// Give API some time to refresh cache
	time.Sleep(time.Second)

	// We should now have one bounce listed when we query the API.
	if !findBounce(exampleEmail) {
		t.Fatalf("Expected bounce for address %s in list of bounces", exampleEmail)
	}

	bounce, err := mg.GetBounce(ctx, exampleEmail)
	ensure.Nil(t, err)
	if bounce.Address != exampleEmail {
		t.Fatalf("Expected at least one bounce for %s", exampleEmail)
	}
	t.Logf("Bounce Created At: %s", bounce.CreatedAt)

	// Delete it.  This should put us back the way we were.
	err = mg.DeleteBounce(ctx, exampleEmail)
	ensure.Nil(t, err)

	// Make sure we're back to the way we were.
	if findBounce(exampleEmail) {
		t.Fatalf("Un-expected bounce for address %s in list of bounces", exampleEmail)
	}

	_, err = mg.GetBounce(ctx, exampleEmail)
	ensure.NotNil(t, err)
}

func TestAddDelBounceList(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	findBounce := func(address string) bool {
		it := mg.ListBounces(nil)
		var page []mailgun.Bounce
		for it.Next(ctx, &page) {
			ensure.True(t, len(page) != 0)
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

	createdAt, err := mailgun.NewRFC2822Time("Thu, 13 Oct 2011 18:02:00 +0000")
	if err != nil {
		t.Fatalf("invalid time")
	}

	// Generate a list of bounces
	bounces := []mailgun.Bounce{
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
	err = mg.AddBounces(ctx, bounces)
	ensure.Nil(t, err)

	for _, expect := range bounces {
		if !findBounce(expect.Address) {
			t.Fatalf("Expected bounce for address %s in list of bounces", expect.Address)
		}

		bounce, err := mg.GetBounce(ctx, expect.Address)
		ensure.Nil(t, err)
		if bounce.Address != expect.Address {
			t.Fatalf("Expected at least one bounce for %s", expect.Address)
		}
		t.Logf("Bounce Created At: %s", bounce.CreatedAt)
		if !expect.CreatedAt.IsZero() && !time.Time(bounce.CreatedAt).Equal(time.Time(expect.CreatedAt)) {
			t.Fatalf("Expected bounce createdAt to be %s, got %s", expect.CreatedAt, bounce.CreatedAt)
		}
	}

	// Delete the bounce list.  This should put us back the way we were.
	err = mg.DeleteBounceList(ctx)
	ensure.Nil(t, err)

	it := mg.ListBounces(nil)
	var page []mailgun.Bounce
	if it.Next(ctx, &page) {
		t.Fatalf("Expected no item in the bounce list")
	}
}
