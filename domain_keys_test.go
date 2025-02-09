package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

func TestDomainDkimSelector(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	// Update Domain DKIM selector
	err = mg.UpdateDomainDkimSelector(ctx, testDomain, true)
	require.NoError(t, err)
}
