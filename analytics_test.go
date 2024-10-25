package mailgun_test

import (
	"context"
	"testing"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// var mockAcceptedIncomingCount uint64 = 10
//
// var expectedResponse = mailgun.MetricsResponse{
// 	// Start:      "Mon, 15 Apr 2024 00:00:00 +0000",
// 	Dimensions: []string{"time"},
// 	Items: []MetricsItem{
// 		{
// 			Dimensions: []MetricsDimension{
// 				{
// 					Dimension:    "time",
// 					Value:        "Mon, 15 Apr 2024 00:00:00 +0000",
// 					DisplayValue: "Mon, 15 Apr 2024 00:00:00 +0000",
// 				},
// 			},
// 			Metrics: Metrics{
// 				AcceptedIncomingCount: &mockAcceptedIncomingCount,
// 				ClickedRate:           "0.8300",
// 			},
// 		},
// 	},
// }

func TestListMetrics(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL1())

	opts := mailgun.MetricsOptions{
		End:      mailgun.RFC2822Time(time.Now().UTC()),
		Duration: "30d",
		Pagination: mailgun.MetricsPagination{
			Limit: 10,
		},
	}
	it, err := mg.ListMetrics(opts)
	require.NoError(t, err)

	var page mailgun.MetricsResponse
	ctx := context.Background()
	more := it.Next(ctx, &page)
	require.Nil(t, it.Err())
	assert.False(t, more)
	assert.Len(t, page.Items, 1)
}
