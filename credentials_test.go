package mailgun_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/facebookgo/ensure"
)

func TestGetCredentials(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	it := mg.ListCredentials(nil)

	var page []mailgun.Credential
	for it.Next(ctx, &page) {
		t.Logf("Login\tCreated At\t\n")
		for _, c := range page {
			t.Logf("%s\t%s\t\n", c.Login, c.CreatedAt)
		}
	}
	ensure.Nil(t, it.Err())
}

func TestCreateDeleteCredentials(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	randomPassword := randomString(16, "pw")
	randomID := strings.ToLower(randomString(16, "usr"))
	randomLogin := fmt.Sprintf("%s@%s", randomID, testDomain)

	ctx := context.Background()
	ensure.Nil(t, mg.CreateCredential(ctx, randomLogin, randomPassword))
	ensure.Nil(t, mg.ChangeCredentialPassword(ctx, randomID, randomString(16, "pw2")))
	ensure.Nil(t, mg.DeleteCredential(ctx, randomID))
}
