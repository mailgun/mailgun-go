package mailgun_test

import (
	"context"
	"testing"

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

	require.Equal(t, "0", list[0].ID)
	require.Equal(t, "/domains", list[0].URL)
	require.Equal(t, "complete", list[0].Status)

	export, err := mg.GetExport(ctx, "0")
	require.NoError(t, err)
	require.Equal(t, "0", export.ID)
	require.Equal(t, "/domains", export.URL)
	require.Equal(t, "complete", export.Status)
}

func TestExportsLink(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	url, err := mg.GetExportLink(ctx, "12")
	require.NoError(t, err)
	require.Contains(t, url, "/some/s3/url")
}
