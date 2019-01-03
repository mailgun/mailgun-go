package mailgun_test

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go"
)

func TestExports(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	list, err := mg.ListExports("")
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(list), 0)

	err = mg.CreateExport("/domains")
	ensure.Nil(t, err)

	list, err = mg.ListExports("")
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(list), 1)

	ensure.DeepEqual(t, list[0].ID, "0")
	ensure.DeepEqual(t, list[0].URL, "/domains")
	ensure.DeepEqual(t, list[0].Status, "complete")

	export, err := mg.GetExport("0")
	ensure.Nil(t, err)
	ensure.DeepEqual(t, export.ID, "0")
	ensure.DeepEqual(t, export.URL, "/domains")
	ensure.DeepEqual(t, export.Status, "complete")
}

func TestExportsLink(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	url, err := mg.GetExportLink("12")
	ensure.Nil(t, err)
	ensure.StringContains(t, url, "/some/s3/url")
}
