package mailgun_test

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

func TestLimits(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	limits, err := mg.GetTagLimits(ctx, testDomain)
	require.NoError(t, err)

	ensure.DeepEqual(t, limits.Limit, 50000)
	ensure.DeepEqual(t, limits.Count, 5000)
}
