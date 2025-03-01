package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDomainTracking(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	info, err := mg.GetDomainTracking(ctx, testDomain)
	require.NoError(t, err)

	require.False(t, info.Unsubscribe.Active)
	require.True(t, info.Unsubscribe.HTMLFooter != "")
	require.True(t, info.Unsubscribe.TextFooter != "")
	require.True(t, info.Click.Active)
	require.True(t, info.Open.Active)

	// Click Tracking
	err = mg.UpdateClickTracking(ctx, testDomain, "no")
	require.NoError(t, err)

	info, err = mg.GetDomainTracking(ctx, testDomain)
	require.NoError(t, err)
	require.False(t, info.Click.Active)

	// Open Tracking
	err = mg.UpdateOpenTracking(ctx, testDomain, "no")
	require.NoError(t, err)

	info, err = mg.GetDomainTracking(ctx, testDomain)
	require.NoError(t, err)
	require.False(t, info.Open.Active)

	// Unsubscribe
	err = mg.UpdateUnsubscribeTracking(ctx, testDomain, "yes", "<h2>Hi</h2>", "Hi")
	require.NoError(t, err)

	info, err = mg.GetDomainTracking(ctx, testDomain)
	require.NoError(t, err)
	assert.True(t, info.Unsubscribe.Active)
	assert.Equal(t, "<h2>Hi</h2>", info.Unsubscribe.HTMLFooter)
	assert.Equal(t, "Hi", info.Unsubscribe.TextFooter)
}

func TestDomainTrackingWebPrefix(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	// Update Domain Tracking Web Prefix
	err = mg.UpdateDomain(ctx, testDomain, &mailgun.UpdateDomainOptions{
		WebScheme: "",
		WebPrefix: "gotest",
	})
	require.NoError(t, err)
}
