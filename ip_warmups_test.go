package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/require"
)

func TestListIPWarmups(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	iter := mg.ListIPWarmups()

	ctx := context.Background()

	var page []mtypes.IPWarmup
	var count int
	for iter.Next(ctx, &page) {
		for _, ip := range page {
			t.Logf("IP Warmup: %#v\n", ip)
			count++
		}
	}
	require.NoError(t, iter.Err())
	require.Equal(t, 2, count)
}

func TestGetIPWarmupStatus(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	detail, err := mg.GetIPWarmupStatus(ctx, "1.0.0.1")
	require.NoError(t, err)
	require.Equal(t, "1.0.0.1", detail.IP)
	require.Len(t, detail.StageHistory, 2)
	t.Logf("IP Warmup: %#v\n", detail)
}

func TestCreateWarmupPlan(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	require.NoError(t, mg.CreateWarmupPlan(ctx, "1.0.0.1"))
}

func TestCancelWarmupPlan(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	require.NoError(t, mg.CancelWarmupPlan(ctx, "1.0.0.1"))
}
