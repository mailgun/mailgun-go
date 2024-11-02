package mailgun_test

import (
	"context"
	"os"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

var server mailgun.MockServer

// Setup and shutdown the mailgun mock server for the entire test suite
func TestMain(m *testing.M) {
	server = mailgun.NewMockServer()
	defer server.Stop()
	os.Exit(m.Run())
}

func TestListIPS(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	list, err := mg.ListIPS(ctx, false)
	require.NoError(t, err)
	require.Len(t, list, 2)

	ip, err := mg.GetIP(ctx, list[0].IP)
	require.NoError(t, err)

	ensure.DeepEqual(t, ip.IP, list[0].IP)
	ensure.DeepEqual(t, ip.Dedicated, true)
	ensure.DeepEqual(t, ip.RDNS, "luna.mailgun.net")
}

func TestDomainIPS(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	err := mg.AddDomainIP(ctx, "192.172.1.1")
	require.NoError(t, err)

	list, err := mg.ListDomainIPS(ctx)
	require.NoError(t, err)

	require.Len(t, list, 1)
	ensure.DeepEqual(t, list[0].IP, "192.172.1.1")

	err = mg.DeleteDomainIP(ctx, "192.172.1.1")
	require.NoError(t, err)

	list, err = mg.ListDomainIPS(ctx)
	require.NoError(t, err)

	require.Len(t, list, 0)
}
