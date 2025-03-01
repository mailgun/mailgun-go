package mailgun_test

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWebhook(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	list, err := mg.ListWebhooks(ctx, testDomain)
	require.NoError(t, err)
	require.Len(t, list, 2)

	urls, err := mg.GetWebhook(ctx, testDomain, "new-webhook")
	require.NoError(t, err)

	assert.Equal(t, []string{"http://example.com/new"}, urls)
}

func TestWebhookCRUD(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	list, err := mg.ListWebhooks(ctx, testDomain)
	require.NoError(t, err)
	require.Len(t, list, 2)

	var countHooks = func() int {
		hooks, err := mg.ListWebhooks(ctx, testDomain)
		require.NoError(t, err)
		return len(hooks)
	}
	hookCount := countHooks()

	webHookURLs := []string{"http://api.mailgun.net/webhook"}
	require.NoError(t, mg.CreateWebhook(ctx, testDomain, "deliver", webHookURLs))

	defer func() {
		require.NoError(t, mg.DeleteWebhook(ctx, testDomain, "deliver"))
		newCount := countHooks()
		require.Equal(t, hookCount, newCount)
	}()

	newCount := countHooks()
	require.False(t, newCount <= hookCount)

	urls, err := mg.GetWebhook(ctx, testDomain, "deliver")
	require.NoError(t, err)
	require.Equal(t, webHookURLs, urls)

	updatedWebHookURL := []string{"http://api.mailgun.net/messages"}
	require.NoError(t, mg.UpdateWebhook(ctx, testDomain, "deliver", updatedWebHookURL))

	hooks, err := mg.ListWebhooks(ctx, testDomain)
	require.NoError(t, err)
	require.Equal(t, updatedWebHookURL, hooks["deliver"])
}

var signedTests = []bool{
	true,
	false,
}

func TestVerifyWebhookSignature(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetWebhookSigningKey(testWebhookSigningKey)

	for _, v := range signedTests {
		fields := getSignatureFields(mg.WebhookSigningKey(), v)
		sig := mtypes.Signature{
			TimeStamp: fields["timestamp"],
			Token:     fields["token"],
			Signature: fields["signature"],
		}

		verified, err := mg.VerifyWebhookSignature(sig)
		require.NoError(t, err)

		if v != verified {
			t.Errorf("VerifyWebhookSignature should return '%v' but got '%v'", v, verified)
		}
	}
}

func getSignatureFields(key string, signed bool) map[string]string {
	badSignature := hex.EncodeToString([]byte("badsignature"))

	fields := map[string]string{
		"token":     "token",
		"timestamp": "123456789",
		"signature": badSignature,
	}

	if signed {
		h := hmac.New(sha256.New, []byte(key))
		_, _ = io.WriteString(h, fields["timestamp"])
		_, _ = io.WriteString(h, fields["token"])
		hash := h.Sum(nil)

		fields["signature"] = hex.EncodeToString(hash)
	}

	return fields
}
