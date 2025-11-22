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

	_, err = mg.VerifyDomain(ctx, testDomain)
	require.NoError(t, err)
}

func TestCreateDomainWithExtendedOptions(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	// Test creating domain with all extended options
	messageTTL := 86400
	_, err = mg.CreateDomain(ctx, "extended.mailgun.test",
		&mailgun.CreateDomainOptions{
			SpamAction:             mtypes.SpamActionTag,
			Password:               "supersecret",
			WebScheme:              "https",
			Wildcard:               true,
			ForceDKIMAuthority:     true,
			DKIMKeySize:            2048,
			ArchiveTo:              "https://archive.example.com/messages",
			DKIMHostName:           "dkim.extended.mailgun.test",
			DKIMSelector:           "mailgun",
			ForceRootDKIMHost:      false,
			EncryptIncomingMessage: true,
			PoolID:                 "pool123",
			RequireTLS:             true,
			SkipVerification:       false,
			WebPrefix:              "tracking",
			MessageTTL:             messageTTL,
		})
	require.NoError(t, err)

	// Verify the domain was created correctly in the mock by checking the stored values
	domains := server.DomainList()
	var found bool
	for _, dc := range domains {
		if dc.Domain.Name != "extended.mailgun.test" {
			continue
		}
		found = true
		assert.Equal(t, mtypes.SpamActionTag, dc.Domain.SpamAction)
		assert.Equal(t, "https", dc.Domain.WebScheme)
		assert.Equal(t, true, dc.Domain.Wildcard)
		assert.Equal(t, "https://archive.example.com/messages", dc.Domain.ArchiveTo)
		assert.Equal(t, "dkim.extended.mailgun.test", dc.Domain.DKIMHost)
		assert.Equal(t, true, dc.Domain.EncryptIncomingMessage)
		assert.Equal(t, true, dc.Domain.RequireTLS)
		assert.Equal(t, false, dc.Domain.SkipVerification)
		assert.Equal(t, "tracking", dc.Domain.WebPrefix)
		assert.Equal(t, 86400, dc.Domain.MessageTTL)
		break
	}
	assert.True(t, found, "Domain should exist in mock server")

	// Clean up
	require.NoError(t, mg.DeleteDomain(ctx, "extended.mailgun.test"))
}

func TestUpdateDomainWithExtendedOptions(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	// First create a domain
	_, err = mg.CreateDomain(ctx, "update-extended.mailgun.test",
		&mailgun.CreateDomainOptions{SpamAction: mtypes.SpamActionTag, Password: "supersecret"})
	require.NoError(t, err)

	// Update with extended options
	requireTLS := true
	skipVerification := false
	useAutoSecurity := true
	messageTTL := 172800

	err = mg.UpdateDomain(ctx, "update-extended.mailgun.test",
		&mailgun.UpdateDomainOptions{
			WebScheme:                  "https",
			WebPrefix:                  "email",
			RequireTLS:                 &requireTLS,
			SkipVerification:           &skipVerification,
			UseAutomaticSenderSecurity: &useAutoSecurity,
			ArchiveTo:                  "https://archive.example.com/messages",
			MailFromHost:               "mail.update-extended.mailgun.test",
			MessageTTL:                 &messageTTL,
		})
	require.NoError(t, err)

	// Verify the domain was updated correctly in the mock by checking the stored values
	domains := server.DomainList()
	var found bool
	for _, dc := range domains {
		if dc.Domain.Name != "update-extended.mailgun.test" {
			continue
		}
		found = true
		assert.Equal(t, "https", dc.Domain.WebScheme)
		assert.Equal(t, "email", dc.Domain.WebPrefix)
		assert.Equal(t, true, dc.Domain.RequireTLS)
		assert.Equal(t, false, dc.Domain.SkipVerification)
		assert.Equal(t, true, dc.Domain.UseAutomaticSenderSecurity)
		assert.Equal(t, "https://archive.example.com/messages", dc.Domain.ArchiveTo)
		assert.Equal(t, "mail.update-extended.mailgun.test", dc.Domain.MailFromHost)
		assert.Equal(t, 172800, dc.Domain.MessageTTL)
		break
	}
	assert.True(t, found, "Domain should exist in mock server")

	// Clean up
	require.NoError(t, mg.DeleteDomain(ctx, "update-extended.mailgun.test"))
}
