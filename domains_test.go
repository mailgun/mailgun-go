package mailgun_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDomain            = "mailgun.test"
	testKey               = "api-fake-key"
	testWebhookSigningKey = "webhook-signing-key"
)

func TestListDomains(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	it := mg.ListDomains(nil)
	var page []mtypes.Domain
	for it.Next(ctx, &page) {
		for _, d := range page {
			t.Logf("TestListDomains: %#v\n", d)
		}
	}
	t.Logf("TestListDomains: %d domains retrieved\n", it.TotalCount)
	require.NoError(t, it.Err())
	assert.True(t, it.TotalCount != 0)
}

func TestGetSingleDomain(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	it := mg.ListDomains(nil)
	var page []mtypes.Domain
	require.True(t, it.Next(ctx, &page))
	require.NoError(t, it.Err())

	dr, err := mg.GetDomain(ctx, page[0].Name, nil)
	require.NoError(t, err)
	require.True(t, len(dr.ReceivingDNSRecords) != 0)
	require.True(t, len(dr.SendingDNSRecords) != 0)

	t.Logf("TestGetSingleDomain: %#v\n", dr)
	for _, rxd := range dr.ReceivingDNSRecords {
		t.Logf("TestGetSingleDomains:   %#v\n", rxd)
	}
	for _, txd := range dr.SendingDNSRecords {
		t.Logf("TestGetSingleDomains:   %#v\n", txd)
	}
}

func TestGetSingleDomainNotExist(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	_, err = mg.GetDomain(ctx, "unknown.domain", nil)
	if err == nil {
		t.Fatal("Did not expect a domain to exist")
	}
	var ure *mailgun.UnexpectedResponseError
	require.ErrorAs(t, err, &ure)
	require.Equal(t, http.StatusNotFound, ure.Actual)
}

func TestAddUpdateDeleteDomain(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	// First, we need to add the domain.
	_, err = mg.CreateDomain(ctx, "mx.mailgun.test",
		&mailgun.CreateDomainOptions{SpamAction: mtypes.SpamActionTag, Password: "supersecret", WebScheme: "http"})
	require.NoError(t, err)

	// Then, we update it.
	err = mg.UpdateDomain(ctx, "mx.mailgun.test",
		&mailgun.UpdateDomainOptions{WebScheme: "https"})
	require.NoError(t, err)

	// Next, we delete it.
	require.NoError(t, mg.DeleteDomain(ctx, "mx.mailgun.test"))
}

func TestDomainVerify(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	_, err = mg.VerifyAndReturnDomain(ctx, testDomain)
	require.NoError(t, err)
}

func TestListIPDomains(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	it := mg.ListIPDomains("192.172.1.1", nil)
	var page []mtypes.DomainIPs
	for it.Next(ctx, &page) {
		for _, d := range page {
			t.Logf("TestListDomains: %#v\n", d)
		}
	}
	t.Logf("TestListDomains: %d domains retrieved\n", it.TotalCount)
	require.NoError(t, it.Err())
	assert.True(t, it.TotalCount != 0)
}
