package mailgun

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestGetBounces(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	ctx := context.Background()
	it := mg.ListBounces(nil)

	var page []Bounce
	for it.Next(ctx, &page) {
		for _, bounce := range page {
			t.Logf("Bounce: %+v\n", bounce)
		}
	}
	ensure.Nil(t, it.Err())
}

func TestGetSingleBounce(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	ctx := context.Background()
	exampleEmail := fmt.Sprintf("%s@%s", strings.ToLower(randomString(64, "")),
		os.Getenv("MG_DOMAIN"))
	_, err = mg.GetBounce(ctx, exampleEmail)
	ensure.NotNil(t, err)

	ure, ok := err.(*UnexpectedResponseError)
	ensure.True(t, ok)
	ensure.DeepEqual(t, ure.Actual, http.StatusNotFound)
}

func TestAddDelBounces(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	domain := os.Getenv("MG_DOMAIN")
	mg, err := NewMailgunFromEnv()
	ctx := context.Background()
	ensure.Nil(t, err)

	findBounce := func(address string) bool {
		it := mg.ListBounces(nil)
		var page []Bounce
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
	err = mg.AddBounce(ctx, exampleEmail, "550", "TestAddDelBounces-generated error")
	ensure.Nil(t, err)

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
