package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("providerLookup=false", func(t *testing.T) {
		ev, err := mg.ValidateEmail(ctx, "foo@mailgun.com", false)
		require.NoError(t, err)

		assert.False(t, ev.IsDisposableAddress)
		assert.False(t, ev.IsRoleAddress)
		require.Len(t, ev.Reason, 1)
		assert.Equal(t, "no-reason", ev.Reason[0])
		assert.Equal(t, "low", ev.Risk)
		assert.Equal(t, "deliverable", ev.Result)
		assert.Equal(t, "disengaged", ev.Engagement.Behavior)
		assert.False(t, ev.Engagement.Engaging)
		assert.False(t, ev.Engagement.IsBot)
	})

	// When providerLookup (provider_lookup) is true, the mock simulates a mailbox
	// provider lookup failure by returning "smtp_timeout" as the reason.
	// https://documentation.mailgun.com/docs/validate/single-valid-ir#reason-explanation
	t.Run("providerLookup=true", func(t *testing.T) {
		ev, err := mg.ValidateEmail(ctx, "foo@mailgun.com", true)
		require.NoError(t, err)

		assert.False(t, ev.IsDisposableAddress)
		assert.False(t, ev.IsRoleAddress)
		require.Len(t, ev.Reason, 1)
		assert.Equal(t, "smtp_timeout", ev.Reason[0])
		assert.Equal(t, "low", ev.Risk)
		assert.Equal(t, "deliverable", ev.Result)
		assert.Equal(t, "disengaged", ev.Engagement.Behavior)
		assert.False(t, ev.Engagement.Engaging)
		assert.False(t, ev.Engagement.IsBot)
	})
}
