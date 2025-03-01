package mailgun_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mocks"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const domain = "valid-mailgun-domain"
const apiKey = "valid-mailgun-api-key" //nolint:gosec // This is a test

var server *mocks.Server

// Setup and shutdown the mailgun mock server for the entire test suite
func TestMain(m *testing.M) {
	server = mocks.NewServer()
	defer server.Stop()
	// TODO: os.Exit will exit, and `defer server.Stop()` will not run
	// switch to testify suite
	//nolint:gocritic // ignored till switched to testify suite
	os.Exit(m.Run())
}

func TestMailgun(t *testing.T) {
	m := mailgun.NewMailgun(apiKey)

	assert.Equal(t, apiKey, m.APIKey())
	assert.Equal(t, http.DefaultClient, m.HTTPClient())

	client := new(http.Client)
	m.SetHTTPClient(client)
	assert.Equal(t, m.HTTPClient(), client)
}

func TestInvalidBaseAPI(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase("https://localhost/v3")
	assert.EqualError(t, err, `APIBase must not contain a version; SetAPIBase("https://host")`)
}

func TestValidBaseAPI(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		var resp mtypes.GetDomainResponse
		b, err := json.Marshal(resp)
		require.NoError(t, err)

		_, err = w.Write(b)
		require.NoError(t, err)
	}))

	apiBases := []string{
		mailgun.APIBase,
		mailgun.APIBaseEU,
		testServer.URL,
	}

	for _, apiBase := range apiBases {
		t.Run(apiBase, func(t *testing.T) {
			mg := mailgun.NewMailgun(testKey)
			err := mg.SetAPIBase(apiBase)
			require.NoError(t, err)
		})
	}
}

func ptr[T any](v T) *T {
	return &v
}
