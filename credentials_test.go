package mailgun

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestGetCredentials(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	ctx := context.Background()
	it := mg.ListCredentials(nil)

	var page []Credential
	for it.Next(ctx, &page) {
		t.Logf("Login\tCreated At\t\n")
		for _, c := range page {
			t.Logf("%s\t%s\t\n", c.Login, c.CreatedAt)
		}
	}
	ensure.Nil(t, it.Err())
}

func TestCreateDeleteCredentials(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	domain := os.Getenv("MG_DOMAIN")
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	randomPassword := randomString(16, "pw")
	randomID := strings.ToLower(randomString(16, "usr"))
	randomLogin := fmt.Sprintf("%s@%s", randomID, domain)

	ctx := context.Background()
	ensure.Nil(t, mg.CreateCredential(ctx, randomLogin, randomPassword))
	ensure.Nil(t, mg.ChangeCredentialPassword(ctx, randomID, randomString(16, "pw2")))
	ensure.Nil(t, mg.DeleteCredential(ctx, randomID))
}
