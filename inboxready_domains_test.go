package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddDomainToMonitoring(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	req := mtypes.AddDomainToMonitoringOptions{
		Domain: testDomain,
	}
	got, err := mg.AddDomainToMonitoring(context.Background(), req)
	require.NoError(t, err)

	assert.Equal(t, testDomain, got.Domain.Name)
}
