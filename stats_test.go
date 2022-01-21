package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v4"

	"github.com/facebookgo/ensure"
)

func TestListStats(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	stats, err := mg.GetStats(ctx, []string{"accepted", "delivered"}, nil)
	ensure.Nil(t, err)

	if len(stats) > 0 {
		firstStatsTotal := stats[0]
		t.Logf("Time: %s\n", firstStatsTotal.Time)
		t.Logf("Accepted Total: %d\n", firstStatsTotal.Accepted.Total)
		t.Logf("Delivered Total: %d\n", firstStatsTotal.Delivered.Total)
	}
}
