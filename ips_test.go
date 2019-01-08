package mailgun_test

import (
	"context"
	"os"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go"
)

var server mailgun.MockServer

// Setup and shutdown the mailgun mock server for the entire test suite
func TestMain(m *testing.M) {
	server = mailgun.NewMockServer()
	defer server.Stop()
	os.Exit(m.Run())
}

func TestListIPS(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	ctx := context.Background()
	list, err := mg.ListIPS(ctx, false)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(list), 2)

	ip, err := mg.GetIP(ctx, list[0].IP)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, ip.IP, list[0].IP)
	ensure.DeepEqual(t, ip.Dedicated, true)
	ensure.DeepEqual(t, ip.RDNS, "luna.mailgun.net")
}

func TestDomainIPS(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	ctx := context.Background()
	err = mg.AddDomainIP(ctx, "192.172.1.1")
	ensure.Nil(t, err)

	list, err := mg.ListDomainIPS(ctx)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, len(list), 1)
	ensure.DeepEqual(t, list[0].IP, "192.172.1.1")

	err = mg.DeleteDomainIP(ctx, "192.172.1.1")
	ensure.Nil(t, err)

	list, err = mg.ListDomainIPS(ctx)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, len(list), 0)
}
