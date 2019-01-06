package mailgun_test

import (
	"net/http"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go"
)

const (
	testDomain = "mailgun.test"
)

func TestGetDomains(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	n, domains, err := mg.ListDomains(mailgun.DefaultLimit, mailgun.DefaultSkip)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(domains) != 0, true)

	t.Logf("TestGetDomains: %d domains retrieved\n", n)
	for _, d := range domains {
		t.Logf("TestGetDomains: %#v\n", d)
	}
}

func TestGetSingleDomain(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	_, domains, err := mg.ListDomains(mailgun.DefaultLimit, mailgun.DefaultSkip)
	ensure.Nil(t, err)

	dr, rxDnsRecords, txDnsRecords, err := mg.GetDomain(domains[0].Name)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(rxDnsRecords) != 0, true)
	ensure.DeepEqual(t, len(txDnsRecords) != 0, true)

	t.Logf("TestGetSingleDomain: %#v\n", dr)
	for _, rxd := range rxDnsRecords {
		t.Logf("TestGetSingleDomains:   %#v\n", rxd)
	}
	for _, txd := range txDnsRecords {
		t.Logf("TestGetSingleDomains:   %#v\n", txd)
	}
}

func TestGetSingleDomainNotExist(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)
	_, _, _, err = mg.GetDomain("unknown.domain")
	if err == nil {
		t.Fatal("Did not expect a domain to exist")
	}
	ure, ok := err.(*mailgun.UnexpectedResponseError)
	ensure.True(t, ok)
	ensure.DeepEqual(t, ure.Actual, http.StatusNotFound)
}

func TestAddDeleteDomain(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	// First, we need to add the domain.
	ensure.Nil(t, mg.CreateDomain("mx.mailgun.test", "supersecret", mailgun.Tag, false))
	// Next, we delete it.
	ensure.Nil(t, mg.DeleteDomain("mx.mailgun.test"))
}

func TestDomainConnection(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	info, err := mg.GetDomainConnection(testDomain)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, info.RequireTLS, true)
	ensure.DeepEqual(t, info.SkipVerification, true)

	info.RequireTLS = false
	err = mg.UpdateDomainConnection(testDomain, info)
	ensure.Nil(t, err)

	info, err = mg.GetDomainConnection(testDomain)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, info.RequireTLS, false)
	ensure.DeepEqual(t, info.SkipVerification, true)
}

func TestDomainTracking(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	info, err := mg.GetDomainTracking(testDomain)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, info.Unsubscribe.Active, false)
	ensure.DeepEqual(t, len(info.Unsubscribe.HTMLFooter) != 0, true)
	ensure.DeepEqual(t, len(info.Unsubscribe.TextFooter) != 0, true)
	ensure.DeepEqual(t, info.Click.Active, true)
	ensure.DeepEqual(t, info.Open.Active, true)
}
