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

	iter := mg.ListIPWarmups(nil)

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
