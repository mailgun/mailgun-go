package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListMetrics(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	start, _ := mtypes.NewRFC2822Time("Tue, 24 Sep 2024 00:00:00 +0000")
	end, _ := mtypes.NewRFC2822Time("Tue, 24 Oct 2024 00:00:00 +0000")

	opts := mtypes.MetricsRequest{
		Start: start,
		End:   end,
		Pagination: mtypes.MetricsPagination{
			Limit: 10,
		},
	}
	// filter by domain
	opts.Filter.BoolGroupAnd = []mtypes.MetricsFilterPredicate{{
		Attribute:     "domain",
		Comparator:    "=",
		LabeledValues: []mtypes.MetricsLabeledValue{{Label: testDomain, Value: testDomain}},
	}}

	wantResp := mtypes.MetricsResponse{
		Start:      start,
		End:        end,
		Resolution: "day",
		Duration:   "30d",
		Dimensions: []string{"time"},
		Items: []mtypes.MetricsItem{
			{
				Dimensions: []mtypes.MetricsDimension{{
					Dimension:    "time",
					Value:        "Tue, 24 Sep 2024 00:00:00 +0000",
					DisplayValue: "Tue, 24 Sep 2024 00:00:00 +0000",
				}},
				Metrics: mtypes.Metrics{
					SentCount:      ptr(uint64(4)),
					DeliveredCount: ptr(uint64(3)),
					OpenedCount:    ptr(uint64(2)),
					FailedCount:    ptr(uint64(1)),
				},
			},
		},
		Pagination: mtypes.MetricsPagination{
			Sort:  "",
			Skip:  0,
			Limit: 10,
			Total: 1,
		},
	}

	it, err := mg.ListMetrics(opts)
	require.NoError(t, err)

	var page mtypes.MetricsResponse
	ctx := context.Background()
	more := it.Next(ctx, &page)
	require.Nil(t, it.Err())
	assert.False(t, more)
	assert.Equal(t, wantResp, page)
}
