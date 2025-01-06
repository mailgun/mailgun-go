package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailValidationV4(t *testing.T) {
	v := mailgun.NewEmailValidator(testKey)
	v.SetAPIBase(server.URL())

	ctx := context.Background()

	ev, err := v.ValidateEmail(ctx, "foo@mailgun.com", false)
	require.NoError(t, err)

	assert.False(t, ev.IsDisposableAddress)
	assert.False(t, ev.IsRoleAddress)
	assert.True(t, len(ev.Reason) != 0)
	assert.Equal(t, "no-reason", ev.Reason[0])
	assert.Equal(t, "unknown", ev.Risk)
	assert.Equal(t, "deliverable", ev.Result)
	assert.Equal(t, "disengaged", ev.Engagement.Behavior)
	assert.False(t, ev.Engagement.Engaging)
	assert.False(t, ev.Engagement.IsBot)
}
