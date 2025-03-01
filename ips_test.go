package mailgun_test

import (
	"context"
	"os"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mailgun/mailgun-go/v4/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var server *mocks.Server

// Setup and shutdown the mailgun mock server for the entire test suite
func TestMain(m *testing.M) {
	server = mocks.NewServer()
	defer server.Stop()
	os.Exit(m.Run())
}

func TestListIPS(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	list, err := mg.ListIPS(ctx, false)
	require.NoError(t, err)
	require.Len(t, list, 2)

	ip, err := mg.GetIP(ctx, list[0].IP)
	require.NoError(t, err)

	assert.Equal(t, list[0].IP, ip.IP)
	assert.True(t, ip.Dedicated)
	assert.Equal(t, "luna.mailgun.net", ip.RDNS)
}

func TestDomainIPS(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	err = mg.AddDomainIP(ctx, testDomain, "192.172.1.1")
	require.NoError(t, err)

	list, err := mg.ListDomainIPS(ctx, testDomain)
	require.NoError(t, err)

	require.Len(t, list, 1)
	require.Equal(t, "192.172.1.1", list[0].IP)

	err = mg.DeleteDomainIP(ctx, testDomain, "192.172.1.1")
	require.NoError(t, err)

	list, err = mg.ListDomainIPS(ctx, testDomain)
	require.NoError(t, err)

	require.Len(t, list, 0)
}
