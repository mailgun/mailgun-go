package mailgun_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestValidateEmailProviderLookupParameter(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		assert.Equal(t, http.MethodGet, req.Method)
		assert.Equal(t, "/v4/address/validate", req.URL.Path)
		assert.Equal(t, "foo@mailgun.com", req.FormValue("address"))
		assert.Equal(t, "true", req.FormValue("provider_lookup"))
		assert.Empty(t, req.FormValue("mailbox_verification"))
		_, err := fmt.Fprint(w, `{"address":"foo@mailgun.com","result":"deliverable","risk":"low","reason":[]}`)
		require.NoError(t, err)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(srv.URL)
	require.NoError(t, err)

	ctx := context.Background()

	_, err = mg.ValidateEmail(ctx, "foo@mailgun.com", true)
	require.NoError(t, err)
}
