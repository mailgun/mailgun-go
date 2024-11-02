package mailgun_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

func TestEmailValidationV4(t *testing.T) {
	v := mailgun.NewEmailValidator(testKey)
	// API Base is set to `http://server/v4`
	v.SetAPIBase(server.URL4())
	ctx := context.Background()

	ev, err := v.ValidateEmail(ctx, "foo@mailgun.com", false)
	require.NoError(t, err)

	require.True(t, ev.IsValid)
	require.Equal(t, "", ev.MailboxVerification)
	require.False(t, ev.IsDisposableAddress)
	require.False(t, ev.IsRoleAddress)
	require.Equal(t, "", ev.Parts.DisplayName)
	require.Equal(t, "foo", ev.Parts.LocalPart)
	require.Equal(t, "mailgun.com", ev.Parts.Domain)
	require.Equal(t, "", ev.Reason)
	require.True(t, len(ev.Reasons) != 0)
	require.Equal(t, "no-reason", ev.Reasons[0])
	require.Equal(t, "unknown", ev.Risk)
	require.Equal(t, "deliverable", ev.Result)
	require.Equal(t, "disengaged", ev.Engagement.Behavior)
	require.False(t, ev.Engagement.Engaging)
	require.False(t, ev.Engagement.IsBot)
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

	require.True(t, ev.IsValid)
	require.Equal(t, "unknown", ev.MailboxVerification)
	require.False(t, ev.IsDisposableAddress)
	require.False(t, ev.IsRoleAddress)
	require.Equal(t, "", ev.Parts.DisplayName)
	require.Equal(t, "some_email", ev.Parts.LocalPart)
	require.Equal(t, "aol.com", ev.Parts.Domain)
	require.Equal(t, "no-reason", ev.Reason)
}
