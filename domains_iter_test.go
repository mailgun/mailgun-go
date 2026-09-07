package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListDomainsIter(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	var domains []mtypes.Domain
	for page, err := range mg.ListDomainsIter(ctx, nil) {
		require.NoError(t, err)
		domains = append(domains, page...)
	}

	t.Logf("TestListDomains: %d domains retrieved", len(domains))
	assert.NotEmpty(t, domains)
}

func TestListDomainsIter_Pagination(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	const limit = 1
	wantDomains := []string{"mailgun.test", "example.com"}

	ctx := context.Background()

	_, err = mg.CreateDomain(ctx, "example.com", &mailgun.CreateDomainOptions{
		SpamAction: mtypes.SpamActionTag,
		Password:   "supersecret",
		WebScheme:  "http",
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		require.NoError(t, mg.DeleteDomain(context.Background(), "example.com"))
	})

	var (
		pages      int
		gotDomains []string
	)
	for page, err := range mg.ListDomainsIter(ctx, &mailgun.ListDomainsOptions{Limit: limit}) {
		require.NoError(t, err)
		require.Len(t, page, limit)
		pages++
		gotDomains = append(gotDomains, page[0].Name)
	}

	assert.Equal(t, len(wantDomains), pages)
	assert.Equal(t, wantDomains, gotDomains)
}

func TestListDomainsIter_Error(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase("http://localhost:9/invalid")
	require.NoError(t, err)

	ctx := context.Background()

	var iterations int
	for _, err := range mg.ListDomainsIter(ctx, nil) {
		iterations++
		assert.Error(t, err)
	}

	assert.Equal(t, 1, iterations)
}
