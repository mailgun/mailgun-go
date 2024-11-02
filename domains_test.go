package mailgun_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDomain = "mailgun.test"
	testKey    = "api-fake-key"
)

func TestListDomains(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	it := mg.ListDomains(nil)
	var page []mailgun.Domain
	for it.Next(ctx, &page) {
		for _, d := range page {
			t.Logf("TestListDomains: %#v\n", d)
		}
	}
	t.Logf("TestListDomains: %d domains retrieved\n", it.TotalCount)
	require.NoError(t, it.Err())
	require.True(t, it.TotalCount != 0)
}

func TestGetSingleDomain(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	it := mg.ListDomains(nil)
	var page []mailgun.Domain
	require.True(t, it.Next(ctx, &page))
	require.NoError(t, it.Err())

	dr, err := mg.GetDomain(ctx, page[0].Name)
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
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	_, err := mg.GetDomain(ctx, "unknown.domain")
	if err == nil {
		t.Fatal("Did not expect a domain to exist")
	}
	var ure *mailgun.UnexpectedResponseError
	require.ErrorAs(t, err, &ure)
	require.Equal(t, http.StatusNotFound, ure.Actual)
}

func TestAddUpdateDeleteDomain(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	// First, we need to add the domain.
	_, err := mg.CreateDomain(ctx, "mx.mailgun.test",
		&mailgun.CreateDomainOptions{SpamAction: mailgun.SpamActionTag, Password: "supersecret", WebScheme: "http"})
	require.NoError(t, err)

	// Then, we update it.
	err = mg.UpdateDomain(ctx, "mx.mailgun.test",
		&mailgun.UpdateDomainOptions{WebScheme: "https"})
	require.NoError(t, err)

	// Next, we delete it.
	require.NoError(t, mg.DeleteDomain(ctx, "mx.mailgun.test"))
}

func TestDomainConnection(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	info, err := mg.GetDomainConnection(ctx, testDomain)
	require.NoError(t, err)

	require.True(t, info.RequireTLS)
	require.True(t, info.SkipVerification)

	info.RequireTLS = false
	err = mg.UpdateDomainConnection(ctx, testDomain, info)
	require.NoError(t, err)

	info, err = mg.GetDomainConnection(ctx, testDomain)
	require.NoError(t, err)
	require.False(t, info.RequireTLS)
	require.True(t, info.SkipVerification)
}

func TestDomainTracking(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	info, err := mg.GetDomainTracking(ctx, testDomain)
	require.NoError(t, err)

	require.False(t, info.Unsubscribe.Active)
	require.True(t, len(info.Unsubscribe.HTMLFooter) != 0)
	require.True(t, len(info.Unsubscribe.TextFooter) != 0)
	require.True(t, info.Click.Active)
	require.True(t, info.Open.Active)

	// Click Tracking
	err = mg.UpdateClickTracking(ctx, testDomain, "no")
	require.NoError(t, err)

	info, err = mg.GetDomainTracking(ctx, testDomain)
	require.NoError(t, err)
	require.False(t, info.Click.Active)

	// Open Tracking
	err = mg.UpdateOpenTracking(ctx, testDomain, "no")
	require.NoError(t, err)

	info, err = mg.GetDomainTracking(ctx, testDomain)
	require.NoError(t, err)
	require.False(t, info.Open.Active)

	// Unsubscribe
	err = mg.UpdateUnsubscribeTracking(ctx, testDomain, "yes", "<h2>Hi</h2>", "Hi")
	require.NoError(t, err)

	info, err = mg.GetDomainTracking(ctx, testDomain)
	require.NoError(t, err)
	assert.True(t, info.Unsubscribe.Active)
	assert.Equal(t, "<h2>Hi</h2>", info.Unsubscribe.HTMLFooter)
	assert.Equal(t, "Hi", info.Unsubscribe.TextFooter)
}

func TestDomainVerify(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	_, err := mg.VerifyDomain(ctx, testDomain)
	require.NoError(t, err)
}

func TestDomainVerifyAndReturn(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	_, err := mg.VerifyAndReturnDomain(ctx, testDomain)
	require.NoError(t, err)
}

func TestDomainDkimSelector(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	// Update Domain DKIM selector
	err := mg.UpdateDomainDkimSelector(ctx, testDomain, "gotest")
	require.NoError(t, err)
}

func TestDomainTrackingWebPrefix(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	// Update Domain Tracking Web Prefix
	err := mg.UpdateDomainTrackingWebPrefix(ctx, testDomain, "gotest")
	require.NoError(t, err)
}
