//go:build integration

package mailgun_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationMailgunImpl_ListMetrics(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		require.NoError(t, err)
	}

	opts := mailgun.MetricsOptions{
		End:      mailgun.RFC2822Time(time.Now().UTC()),
		Duration: "30d",
		Pagination: mailgun.MetricsPagination{
			Limit: 10,
		},
	}

	iter, err := mg.ListMetrics(opts)
	require.NoError(t, err)

	// create context to list all pages
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	for i := 0; i < 2; i++ {
		var resp mailgun.MetricsResponse
		more := iter.Next(ctx, &resp)
		if iter.Err() != nil {
			require.NoError(t, err)
		}

		t.Logf("Page %d: Start: %s; End: %s; Pagination: %+v\n",
			i+1, resp.Start, resp.End, resp.Pagination)

		assert.GreaterOrEqual(t, len(resp.Items), 1)

		for _, item := range resp.Items {
			b, _ := json.Marshal(item)
			t.Logf("%s\n", b)
		}

		if !more {
			t.Log("no more pages")
			break
		}
	}
}

func TestIntegrationWebhooksCRUD(t *testing.T) {
	// Arrange

	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		require.NoError(t, err)
	}

	const name = "permanent_fail"
	ctx := context.Background()
	urls := []string{"https://example.com/1", "https://example.com/2"}

	err = mg.DeleteWebhook(ctx, name)
	if err != nil {
		// 200 or 404 is expected
		status := mailgun.GetStatusFromErr(err)
		require.Equal(t, http.StatusNotFound, status, err)
	}
	time.Sleep(3 * time.Second)

	defer func() {
		// Cleanup
		_ = mg.DeleteWebhook(ctx, name)
	}()

	// Act

	err = mg.CreateWebhook(ctx, name, urls)
	require.NoError(t, err)
	time.Sleep(3 * time.Second)

	// Assert

	gotUrls, err := mg.GetWebhook(ctx, name)
	require.NoError(t, err)
	t.Logf("Webhooks: %v", urls)
	assert.ElementsMatch(t, urls, gotUrls)
}
