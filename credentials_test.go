package mailgun

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestGetCredentials(t *testing.T) {
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	ctx := context.Background()
	n, cs, err := mg.ListCredentials(ctx, DefaultLimit, DefaultSkip)
	ensure.Nil(t, err)

	t.Logf("Login\tCreated At\t\n")
	for _, c := range cs {
		t.Logf("%s\t%s\t\n", c.Login, c.CreatedAt)
	}
	t.Logf("%d credentials listed out of %d\n", len(cs), n)
}

func TestCreateDeleteCredentials(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
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
