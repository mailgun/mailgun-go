package mailgun_test

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go"
)

func TestExports(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	ctx := context.Background()
	list, err := mg.ListExports(ctx, "")
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(list), 0)

	err = mg.CreateExport(ctx, "/domains")
	ensure.Nil(t, err)

	list, err = mg.ListExports(ctx, "")
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(list), 1)

	ensure.DeepEqual(t, list[0].ID, "0")
	ensure.DeepEqual(t, list[0].URL, "/domains")
	ensure.DeepEqual(t, list[0].Status, "complete")

	export, err := mg.GetExport(ctx, "0")
	ensure.Nil(t, err)
	ensure.DeepEqual(t, export.ID, "0")
	ensure.DeepEqual(t, export.URL, "/domains")
	ensure.DeepEqual(t, export.Status, "complete")
}

func TestExportsLink(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	ctx := context.Background()
	url, err := mg.GetExportLink(ctx, "12")
	ensure.Nil(t, err)
	ensure.StringContains(t, url, "/some/s3/url")
}
