package mailgun_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetWebhook(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	list, err := mg.ListWebhooks(ctx)
	require.NoError(t, err)
	require.Len(t, list, 2)

	urls, err := mg.GetWebhook(ctx, "new-webhook")
	require.NoError(t, err)

	assert.Equal(t, []string{"http://example.com/new"}, urls)
}

func TestWebhookCRUD(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	list, err := mg.ListWebhooks(ctx)
	require.NoError(t, err)
	require.Len(t, list, 2)

	var countHooks = func() int {
		hooks, err := mg.ListWebhooks(ctx)
		require.NoError(t, err)
		return len(hooks)
	}
	hookCount := countHooks()

	webHookURLs := []string{"http://api.mailgun.net/webhook"}
	require.NoError(t, mg.CreateWebhook(ctx, "deliver", webHookURLs))

	defer func() {
		require.NoError(t, mg.DeleteWebhook(ctx, "deliver"))
		newCount := countHooks()
		require.Equal(t, hookCount, newCount)
	}()

	newCount := countHooks()
	require.False(t, newCount <= hookCount)

	urls, err := mg.GetWebhook(ctx, "deliver")
	require.NoError(t, err)
	require.Equal(t, webHookURLs, urls)

	updatedWebHookURL := []string{"http://api.mailgun.net/messages"}
	require.NoError(t, mg.UpdateWebhook(ctx, "deliver", updatedWebHookURL))

	hooks, err := mg.ListWebhooks(ctx)
	require.NoError(t, err)
	require.Equal(t, updatedWebHookURL, hooks["deliver"])
}

var signedTests = []bool{
	true,
	false,
}

func TestVerifyWebhookSignature(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetWebhookSigningKey(testWebhookSigningKey)

	for _, v := range signedTests {
		fields := getSignatureFields(mg.WebhookSigningKey(), v)
		sig := mailgun.Signature{
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

func TestVerifyWebhookRequest_Form(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetWebhookSigningKey(testWebhookSigningKey)

	for _, v := range signedTests {
		fields := getSignatureFields(mg.WebhookSigningKey(), v)
		req := buildFormRequest(fields)

		verified, err := mg.VerifyWebhookRequest(req)
		require.NoError(t, err)

		if v != verified {
			t.Errorf("VerifyWebhookRequest should return '%v' but got '%v'", v, verified)
		}
	}
}

func TestVerifyWebhookRequest_MultipartForm(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetWebhookSigningKey(testWebhookSigningKey)

	for _, v := range signedTests {
		fields := getSignatureFields(mg.WebhookSigningKey(), v)
		req := buildMultipartFormRequest(fields)

		verified, err := mg.VerifyWebhookRequest(req)
		require.NoError(t, err)

		if v != verified {
			t.Errorf("VerifyWebhookRequest should return '%v' but got '%v'", v, verified)
		}
	}
}

func buildFormRequest(fields map[string]string) *http.Request {
	values := url.Values{}

	for k, v := range fields {
		values.Add(k, v)
	}

	r := strings.NewReader(values.Encode())
	req, _ := http.NewRequest("POST", "/", r)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	return req
}

func buildMultipartFormRequest(fields map[string]string) *http.Request {
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)

	for k, v := range fields {
		_ = writer.WriteField(k, v)
	}

	writer.Close()

	req, _ := http.NewRequest("POST", "/", buf)
	req.Header.Set("Content-type", writer.FormDataContentType())

	return req
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
