package mailgun

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestListStats(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
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

func TestDeleteTag(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := NewMailgunFromEnv()
	ctx := context.Background()

	ensure.Nil(t, err)
	ensure.Nil(t, mg.DeleteTag(ctx, "newsletter"))
}
