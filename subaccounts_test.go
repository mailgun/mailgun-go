package mailgun_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
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
	ensure.NotNil(t, iterator)

	ctx := context.Background()

	var page []mailgun.Subaccount
	for iterator.Next(ctx, &page) {
		for _, d := range page {
			t.Logf("TestListSubaccounts: %#v\n", d)
		}
	}
	t.Logf("TestListSubaccounts: %d subaccounts retrieved\n", iterator.Total)
	ensure.Nil(t, iterator.Err())
	ensure.True(t, iterator.Total != 0)
}

func TestSubaccountDetails(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	iterator := mg.ListSubaccounts(nil)
	ensure.NotNil(t, iterator)

	page := []mailgun.Subaccount{}
	ensure.True(t, iterator.Next(context.Background(), &page))
	ensure.Nil(t, iterator.Err())

	resp, err := mg.SubaccountDetails(ctx, page[0].Id)
	ensure.Nil(t, err)
	ensure.NotNil(t, resp)
}

func TestSubaccountDetailsStatusNotFound(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	_, err := mg.SubaccountDetails(ctx, "unexisting.id")
	if err == nil {
		t.Fatal("Did not expect a subaccount to exist")
	}
	ure, ok := err.(*mailgun.UnexpectedResponseError)
	ensure.True(t, ok)
	ensure.DeepEqual(t, ure.Actual, http.StatusNotFound)
}

func TestCreateSubaccount(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	resp, err := mg.CreateSubaccount(ctx, testSubaccountName)
	ensure.Nil(t, err)
	ensure.NotNil(t, resp)
}

func TestEnableSubaccountAlreadyEnabled(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	_, err := mg.EnableSubaccount(ctx, testEnabledSubaccountId)
	ensure.Nil(t, err)
}

func TestEnableSubaccount(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	resp, err := mg.EnableSubaccount(ctx, testDisabledSubaccountId)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp.Item.Status, "enabled")
}

func TestDisableSubaccount(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	resp, err := mg.DisableSubaccount(ctx, testEnabledSubaccountId)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, resp.Item.Status, "disabled")
}

func TestDisableSubaccountAlreadyDisabled(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	_, err := mg.DisableSubaccount(ctx, testDisabledSubaccountId)
	ensure.Nil(t, err)
}
