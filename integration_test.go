//go:build integration

package mailgun_test

import (
	"context"
	"encoding/json"
	"os"
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
	if err != nil {
		require.NoError(t, err)
	}

	for i := 0; i < 2; i++ {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp := mailgun.MetricsResponse{}
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
				os.Exit(0)
			}
		}()
	}
}
