package mailgun_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/assert"
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
	assert.EqualError(t, err, `APIBase must end with a /v1, /v2, /v3 or /v4; SetAPIBase("https://host/v3")`)
}

func TestValidBaseAPI(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp mailgun.DomainResponse
		b, err := json.Marshal(resp)
		ensure.Nil(t, err)

		w.Write(b)
	}))

	apiBases := []string{
		fmt.Sprintf("%s/v3", testServer.URL),
		fmt.Sprintf("%s/proxy/v3", testServer.URL),
	}

	for _, apiBase := range apiBases {
		mg := mailgun.NewMailgun(testDomain, testKey)
		mg.SetAPIBase(apiBase)

		ctx := context.Background()
		_, err := mg.GetDomain(ctx, "unknown.domain")
		ensure.Nil(t, err)
	}
}

func ptr[T any](v T) *T {
	return &v
}
