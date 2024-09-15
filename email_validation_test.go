package mailgun_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
)

func TestEmailValidationV4(t *testing.T) {
	v := mailgun.NewEmailValidator(testKey)
	// API Base is set to `http://server/v4`
	v.SetAPIBase(server.URL4())
	ctx := context.Background()

	ev, err := v.ValidateEmail(ctx, "foo@mailgun.com", false)
	ensure.Nil(t, err)

	ensure.True(t, ev.IsValid)
	ensure.DeepEqual(t, ev.MailboxVerification, "")
	ensure.False(t, ev.IsDisposableAddress)
	ensure.False(t, ev.IsRoleAddress)
	ensure.True(t, ev.Parts.DisplayName == "")
	ensure.DeepEqual(t, ev.Parts.LocalPart, "foo")
	ensure.DeepEqual(t, ev.Parts.Domain, "mailgun.com")
	ensure.DeepEqual(t, ev.Reason, "")
	ensure.True(t, len(ev.Reasons) != 0)
	ensure.DeepEqual(t, ev.Reasons[0], "no-reason")
	ensure.DeepEqual(t, ev.Risk, "unknown")
	ensure.DeepEqual(t, ev.Result, "deliverable")
	ensure.DeepEqual(t, ev.Engagement.Behavior, "disengaged")
	ensure.DeepEqual(t, ev.Engagement.Engaging, false)
	ensure.False(t, ev.Engagement.IsBot)
}

func TestParseAddresses(t *testing.T) {
	v := mailgun.NewEmailValidator(testKey)
	v.SetAPIBase(server.URL())
	ctx := context.Background()

	addressesThatParsed, unparsableAddresses, err := v.ParseAddresses(ctx,
		"Alice <alice@example.com>",
		"bob@example.com",
		"example.com")
	ensure.Nil(t, err)
	hittest := map[string]bool{
		"Alice <alice@example.com>": true,
		"bob@example.com":           true,
	}
	for _, a := range addressesThatParsed {
		ensure.True(t, hittest[a])
	}
	ensure.True(t, len(unparsableAddresses) == 1)
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
	ensure.Nil(t, err)

	ensure.True(t, ev.IsValid)
	ensure.DeepEqual(t, ev.MailboxVerification, "unknown")
	ensure.False(t, ev.IsDisposableAddress)
	ensure.False(t, ev.IsRoleAddress)
	ensure.True(t, ev.Parts.DisplayName == "")
	ensure.DeepEqual(t, ev.Parts.LocalPart, "some_email")
	ensure.DeepEqual(t, ev.Parts.Domain, "aol.com")
	ensure.DeepEqual(t, ev.Reason, "no-reason")
}
