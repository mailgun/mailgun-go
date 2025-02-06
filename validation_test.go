package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateEmail(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	ev, err := mg.ValidateEmail(ctx, "foo@mailgun.com", false)
	require.NoError(t, err)

	assert.False(t, ev.IsDisposableAddress)
	assert.False(t, ev.IsRoleAddress)
	assert.True(t, len(ev.Reason) != 0)
	assert.Equal(t, "no-reason", ev.Reason[0])
	assert.Equal(t, "low", ev.Risk)
	assert.Equal(t, "deliverable", ev.Result)
	assert.Equal(t, "disengaged", ev.Engagement.Behavior)
	assert.False(t, ev.Engagement.Engaging)
	assert.False(t, ev.Engagement.IsBot)
}
