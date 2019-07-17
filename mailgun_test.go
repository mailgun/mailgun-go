package mailgun_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v3"
)

const domain = "valid-mailgun-domain"
const apiKey = "valid-mailgun-api-key"

func TestMailgun(t *testing.T) {
	m := mailgun.NewMailgun(domain, apiKey)

	ensure.DeepEqual(t, m.Domain(), domain)
	ensure.DeepEqual(t, m.APIKey(), apiKey)
	ensure.DeepEqual(t, m.Client(), http.DefaultClient)

	client := new(http.Client)
	m.SetClient(client)
	ensure.DeepEqual(t, client, m.Client())
}

func TestInvalidBaseAPI(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase("https://localhost")

	ctx := context.Background()
	_, err := mg.GetDomain(ctx, "unknown.domain")
	ensure.NotNil(t, err)
	ensure.DeepEqual(t, err.Error(), `BaseAPI must end with a /v2, /v3 or /v4; setBaseAPI("https://host/v3")`)
}
