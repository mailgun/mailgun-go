package mailgun_test

import (
	"os"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go"
	"github.com/mailgun/mailgun-go/mock"
)

var server mock.MailgunServer

// Setup and shutdown the mailgun mock server for the entire test suite
func TestMain(m *testing.M) {
	server = mock.NewServer()
	defer server.Stop()
	os.Exit(m.Run())
}

func TestListIPS(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	list, err := mg.ListIPS(false)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(list), 2)

	ip, err := mg.GetIP(list[0].IP)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, ip.IP, list[0].IP)
	ensure.DeepEqual(t, ip.Dedicated, true)
	ensure.DeepEqual(t, ip.RDNS, "luna.mailgun.net")
}

func TestDomainIPS(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	err = mg.AddDomainIP("192.172.1.1")
	ensure.Nil(t, err)

	list, err := mg.ListDomainIPS()
	ensure.Nil(t, err)

	ensure.DeepEqual(t, len(list), 1)
	ensure.DeepEqual(t, list[0].IP, "192.172.1.1")

	err = mg.DeleteDomainIP("192.172.1.1")
	ensure.Nil(t, err)

	list, err = mg.ListDomainIPS()
	ensure.Nil(t, err)

	ensure.DeepEqual(t, len(list), 0)
}
