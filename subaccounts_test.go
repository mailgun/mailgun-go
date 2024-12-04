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
	testSubaccountName       = "mailgun.test"
	testEnabledSubaccountId  = "enabled.subaccount"
	testDisabledSubaccountId = "disabled.subaccount"
)

func TestListSubaccounts(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	iterator := mg.ListSubaccounts(nil)
	require.NotNil(t, iterator)

	ctx := context.Background()

	var page []mailgun.Subaccount
	for iterator.Next(ctx, &page) {
		for _, d := range page {
			t.Logf("TestListSubaccounts: %#v\n", d)
		}
	}
	t.Logf("TestListSubaccounts: %d subaccounts retrieved\n", iterator.Total)
	require.NoError(t, iterator.Err())
	require.True(t, iterator.Total != 0)
}

func TestSubaccountDetails(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	iterator := mg.ListSubaccounts(nil)
	require.NotNil(t, iterator)

	page := make([]mailgun.Subaccount, 0, 1)
	require.True(t, iterator.Next(context.Background(), &page))
	require.NoError(t, iterator.Err())

	resp, err := mg.SubaccountDetails(ctx, page[0].Id)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestSubaccountDetailsStatusNotFound(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	_, err := mg.SubaccountDetails(ctx, "unexisting.id")
	if err == nil {
		t.Fatal("Did not expect a subaccount to exist")
	}
	var ure *mailgun.UnexpectedResponseError
	require.ErrorAs(t, err, &ure)
	require.Equal(t, http.StatusNotFound, ure.Actual)
}

func TestCreateSubaccount(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	resp, err := mg.CreateSubaccount(ctx, testSubaccountName)
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestEnableSubaccountAlreadyEnabled(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	_, err := mg.EnableSubaccount(ctx, testEnabledSubaccountId)
	require.NoError(t, err)
}

func TestEnableSubaccount(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	resp, err := mg.EnableSubaccount(ctx, testDisabledSubaccountId)
	require.NoError(t, err)
	assert.Equal(t, "enabled", resp.Item.Status)
}

func TestDisableSubaccount(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	resp, err := mg.DisableSubaccount(ctx, testEnabledSubaccountId)
	require.NoError(t, err)
	assert.Equal(t, "disabled", resp.Item.Status)
}

func TestDisableSubaccountAlreadyDisabled(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	_, err := mg.DisableSubaccount(ctx, testDisabledSubaccountId)
	require.NoError(t, err)
}
