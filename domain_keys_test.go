package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/require"
)

func TestListAllDomainsKeys(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	it := mg.ListAllDomainsKeys(nil)
	var page []mtypes.DomainKey
	require.True(t, it.Next(ctx, &page))
	require.NoError(t, it.Err())
	require.Equal(t, 2, len(page))

	it = mg.ListAllDomainsKeys(&mailgun.ListDomainKeysOptions{Limit: 1})
	var pageWithLimit []mtypes.DomainKey
	require.True(t, it.Next(ctx, &pageWithLimit))
	require.NoError(t, it.Err())
	require.Equal(t, 1, len(pageWithLimit))
}

func TestUpdateDomainDkimSelector(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	// Update Domain DKIM selector
	err = mg.UpdateDomainDkimSelector(ctx, testDomain, "gotest")
	require.NoError(t, err)
}
