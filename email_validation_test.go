package mailgun_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEmailValidationV4(t *testing.T) {
	v := mailgun.NewEmailValidator(testKey)
	// API Base is set to `http://server/v4`
	v.SetAPIBase(server.URL4())
	ctx := context.Background()

	ev, err := v.ValidateEmail(ctx, "foo@mailgun.com", false)
	require.NoError(t, err)

	assert.True(t, ev.IsValid)
	assert.Equal(t, "", ev.MailboxVerification)
	assert.False(t, ev.IsDisposableAddress)
	assert.False(t, ev.IsRoleAddress)
	assert.Equal(t, "", ev.Parts.DisplayName)
	assert.Equal(t, "foo", ev.Parts.LocalPart)
	assert.Equal(t, "mailgun.com", ev.Parts.Domain)
	assert.Equal(t, "", ev.Reason)
	assert.True(t, len(ev.Reasons) != 0)
	assert.Equal(t, "no-reason", ev.Reasons[0])
	assert.Equal(t, "low", ev.Risk)
	assert.Equal(t, "deliverable", ev.Result)
	assert.Equal(t, "disengaged", ev.Engagement.Behavior)
	assert.False(t, ev.Engagement.Engaging)
	assert.False(t, ev.Engagement.IsBot)
}

func TestParseAddresses(t *testing.T) {
	v := mailgun.NewEmailValidator(testKey)
	v.SetAPIBase(server.URL())
	ctx := context.Background()

	addressesThatParsed, unparsableAddresses, err := v.ParseAddresses(ctx,
		"Alice <alice@example.com>",
		"bob@example.com",
		"example.com")
	require.NoError(t, err)
	hittest := map[string]bool{
		"Alice <alice@example.com>": true,
		"bob@example.com":           true,
	}
	for _, a := range addressesThatParsed {
		require.True(t, hittest[a])
	}
	require.Len(t, unparsableAddresses, 1)
}

func TestUnmarshallResponse(t *testing.T) {
	payload := []byte(`{
		"address": "some_email@aol.com",
		"did_you_mean": null,
		"is_disposable_address": false,
		"is_role_address": false,
		"is_valid": true,
		"mailbox_verification": "unknown",
		"parts":
		{
			"display_name": null,
			"domain": "aol.com",
			"local_part": "some_email"
		},
		"reason": "no-reason"
	}`)
	var ev mailgun.EmailVerification
	err := json.Unmarshal(payload, &ev)
	require.NoError(t, err)

	assert.True(t, ev.IsValid)
	assert.Equal(t, "unknown", ev.MailboxVerification)
	assert.False(t, ev.IsDisposableAddress)
	assert.False(t, ev.IsRoleAddress)
	assert.Equal(t, "", ev.Parts.DisplayName)
	assert.Equal(t, "some_email", ev.Parts.LocalPart)
	assert.Equal(t, "aol.com", ev.Parts.Domain)
	assert.Equal(t, "no-reason", ev.Reason)
}
