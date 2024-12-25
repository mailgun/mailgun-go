package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListMetrics(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL1())

	start, _ := mailgun.NewRFC2822Time("Tue, 24 Sep 2024 00:00:00 +0000")
	end, _ := mailgun.NewRFC2822Time("Tue, 24 Oct 2024 00:00:00 +0000")

	opts := mailgun.MetricsOptions{
		Start: start,
		End:   end,
		Pagination: mailgun.MetricsPagination{
			Limit: 10,
		},
	}
	// filter by domain
	opts.Filter.BoolGroupAnd = []mailgun.MetricsFilterPredicate{{
		Attribute:     "domain",
		Comparator:    "=",
		LabeledValues: []mailgun.MetricsLabeledValue{{Label: testDomain, Value: testDomain}},
	}}

	wantResp := mailgun.MetricsResponse{
		Start:      start,
		End:        end,
		Resolution: "day",
		Duration:   "30d",
		Dimensions: []string{"time"},
		Items: []mailgun.MetricsItem{
			{
				Dimensions: []mailgun.MetricsDimension{{
					Dimension:    "time",
					Value:        "Tue, 24 Sep 2024 00:00:00 +0000",
					DisplayValue: "Tue, 24 Sep 2024 00:00:00 +0000",
				}},
				Metrics: mailgun.Metrics{
					SentCount:      ptr(uint64(4)),
					DeliveredCount: ptr(uint64(3)),
					OpenedCount:    ptr(uint64(2)),
					FailedCount:    ptr(uint64(1)),
				},
			},
		},
		Pagination: mailgun.MetricsPagination{
			Sort:  "",
			Skip:  0,
			Limit: 10,
			Total: 1,
		},
	}

	it, err := mg.ListMetrics(opts)
	require.NoError(t, err)

	var page mailgun.MetricsResponse
	ctx := context.Background()
	more := it.Next(ctx, &page)
	require.Nil(t, it.Err())
	assert.False(t, more)
	assert.Equal(t, wantResp, page)
}
