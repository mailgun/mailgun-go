//go:build integration

package mailgun_test

import (
	"context"
	"testing"
	"time"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationListDomainsIter(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		require.NoError(t, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var domains []mtypes.Domain
	for page, err := range mg.ListDomainsIter(ctx, nil) {
		require.NoError(t, err)
		domains = append(domains, page...)
	}

	t.Logf("TestListDomains: %d domains retrieved", len(domains))
	assert.NotEmpty(t, domains)
}
