package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLimits(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	limits, err := mg.GetTagLimits(ctx, testDomain)
	require.NoError(t, err)

	assert.Equal(t, 50000, limits.Limit)
	assert.Equal(t, 5000, limits.Count)
}
