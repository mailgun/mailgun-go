package mailgun_test

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

func TestExports(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	list, err := mg.ListExports(ctx, "")
	require.NoError(t, err)
	require.Len(t, list, 0)

	err = mg.CreateExport(ctx, "/domains")
	require.NoError(t, err)

	list, err = mg.ListExports(ctx, "")
	require.NoError(t, err)
	require.Len(t, list, 1)

	ensure.DeepEqual(t, list[0].ID, "0")
	ensure.DeepEqual(t, list[0].URL, "/domains")
	ensure.DeepEqual(t, list[0].Status, "complete")

	export, err := mg.GetExport(ctx, "0")
	require.NoError(t, err)
	ensure.DeepEqual(t, export.ID, "0")
	ensure.DeepEqual(t, export.URL, "/domains")
	ensure.DeepEqual(t, export.Status, "complete")
}

func TestExportsLink(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	url, err := mg.GetExportLink(ctx, "12")
	require.NoError(t, err)
	ensure.StringContains(t, url, "/some/s3/url")
}
