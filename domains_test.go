package mailgun_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v3"
)

const (
	testDomain = "mailgun.test"
	testKey    = "api-fake-key"
)

func TestGetDomains(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	n, domains, err := mg.ListDomains(ctx, nil)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(domains) != 0, true)

	t.Logf("TestGetDomains: %d domains retrieved\n", n)
	for _, d := range domains {
		t.Logf("TestGetDomains: %#v\n", d)
	}
}

func TestGetSingleDomain(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	_, domains, err := mg.ListDomains(ctx, nil)
	ensure.Nil(t, err)

	dr, rxDnsRecords, txDnsRecords, err := mg.GetDomain(ctx, domains[0].Name)
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
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	_, _, _, err := mg.GetDomain(ctx, "unknown.domain")
	if err == nil {
		t.Fatal("Did not expect a domain to exist")
	}
	ure, ok := err.(*mailgun.UnexpectedResponseError)
	ensure.True(t, ok)
	ensure.DeepEqual(t, ure.Actual, http.StatusNotFound)
}

func TestAddDeleteDomain(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	// First, we need to add the domain.
	ensure.Nil(t, mg.CreateDomain(ctx, "mx.mailgun.test", "supersecret",
		&mailgun.CreateDomainOptions{SpamAction: mailgun.SpamActionTag}))
	// Next, we delete it.
	ensure.Nil(t, mg.DeleteDomain(ctx, "mx.mailgun.test"))
}

func TestDomainConnection(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	info, err := mg.GetDomainConnection(ctx, testDomain)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, info.RequireTLS, true)
	ensure.DeepEqual(t, info.SkipVerification, true)

	info.RequireTLS = false
	err = mg.UpdateDomainConnection(ctx, testDomain, info)
	ensure.Nil(t, err)

	info, err = mg.GetDomainConnection(ctx, testDomain)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, info.RequireTLS, false)
	ensure.DeepEqual(t, info.SkipVerification, true)
}

func TestDomainTracking(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	info, err := mg.GetDomainTracking(ctx, testDomain)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, info.Unsubscribe.Active, false)
	ensure.DeepEqual(t, len(info.Unsubscribe.HTMLFooter) != 0, true)
	ensure.DeepEqual(t, len(info.Unsubscribe.TextFooter) != 0, true)
	ensure.DeepEqual(t, info.Click.Active, true)
	ensure.DeepEqual(t, info.Open.Active, true)
}
