package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

func TestDomainConnection(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	info, err := mg.GetDomainConnection(ctx, testDomain)
	require.NoError(t, err)

	require.True(t, info.RequireTLS)
	require.True(t, info.SkipVerification)

	info.RequireTLS = false
	err = mg.UpdateDomainConnection(ctx, testDomain, info)
	require.NoError(t, err)

	info, err = mg.GetDomainConnection(ctx, testDomain)
	require.NoError(t, err)
	require.False(t, info.RequireTLS)
	require.True(t, info.SkipVerification)
}
