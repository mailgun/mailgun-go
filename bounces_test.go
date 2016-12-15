package mailgun

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetBounces(t *testing.T) {
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}
	n, bounces, err := mg.GetBounces(-1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(bounces) {
		t.Fatalf("Expected length of bounces %d to equal returned length %d", len(bounces), n)
	}
}

func TestGetSingleBounce(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}
	exampleEmail := fmt.Sprintf("%s@%s", strings.ToLower(randomString(64, "")), domain)
	_, err = mg.GetSingleBounce(exampleEmail)
	if err == nil {
		t.Fatal("Did not expect a bounce to exist")
	}
	ure, ok := err.(*UnexpectedResponseError)
	if !ok {
		t.Fatal("Expected UnexpectedResponseError")
	}
	if ure.Actual != 404 {
		t.Fatalf("Expected 404 response code; got %d", ure.Actual)
	}
}

func TestAddDelBounces(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}

	// Compute an e-mail address for our domain.
	exampleEmail := fmt.Sprintf("%s@%s", strings.ToLower(randomString(8, "bounce")), domain)

	// First, basic sanity check.
	// Fail early if we have bounces for a fictitious e-mail address.

	n, _, err := mg.GetBounces(-1, -1)
	if err != nil {
		t.Fatal(err)
	}
	// Add the bounce for our address.

	err = mg.AddBounce(exampleEmail, "550", "TestAddDelBounces-generated error")
	if err != nil {
		t.Fatal(err)
	}

	// We should now have one bounce listed when we query the API.

	n, bounces, err := mg.GetBounces(-1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if n == 0 {
		t.Fatal("Expected at least one bounce for this domain.")
	}

	found := 0
	for _, bounce := range bounces {
		t.Logf("Bounce Address: %s\n", bounce.Address)
		if bounce.Address == exampleEmail {
			found++
		}
	}

	if found == 0 {
		t.Fatalf("Expected bounce for address %s in list of bounces", exampleEmail)
	}

	bounce, err := mg.GetSingleBounce(exampleEmail)
	if err != nil {
		t.Fatal(err)
	}
	if bounce.CreatedAt == "" {
		t.Fatalf("Expected at least one bounce for %s", exampleEmail)
	}

	// Delete it.  This should put us back the way we were.

	err = mg.DeleteBounce(exampleEmail)
	if err != nil {
		t.Fatal(err)
	}

	// Make sure we're back to the way we were.

	n, bounces, err = mg.GetBounces(-1, -1)
	if err != nil {
		t.Fatal(err)
	}

	found = 0
	for _, bounce := range bounces {
		t.Logf("Bounce Address: %s\n", bounce.Address)
		if bounce.Address == exampleEmail {
			found++
		}
	}

	if found != 0 {
		t.Fatalf("Expected no bounce for address %s in list of bounces", exampleEmail)
	}

	_, err = mg.GetSingleBounce(exampleEmail)
	if err == nil {
		t.Fatalf("Expected no bounces for %s", exampleEmail)
	}
}
