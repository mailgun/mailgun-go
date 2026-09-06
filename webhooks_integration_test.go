//go:build integration

package mailgun_test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationWebhooksCRUD(t *testing.T) {
	// Arrange

	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		require.NoError(t, err)
	}

	domain := os.Getenv("MG_DOMAIN")
	require.NotEmpty(t, domain)

	const name = "permanent_fail"
	ctx := context.Background()
	urls := []string{"https://example.com/1", "https://example.com/2"}

	err = mg.DeleteWebhook(ctx, domain, name)
	if err != nil {
		// 200 or 404 is expected
		status := mailgun.GetStatusFromErr(err)
		require.Equal(t, http.StatusNotFound, status, err)
	}
	time.Sleep(3 * time.Second)

	defer func() {
		// Cleanup
		_ = mg.DeleteWebhook(ctx, domain, name)
	}()

	// Act

	err = mg.CreateWebhook(ctx, domain, name, urls)
	require.NoError(t, err)
	time.Sleep(3 * time.Second)

	// Assert

	gotUrls, err := mg.GetWebhook(ctx, domain, name)
	require.NoError(t, err)
	t.Logf("Webhooks: %v", urls)
	assert.ElementsMatch(t, urls, gotUrls)
}
