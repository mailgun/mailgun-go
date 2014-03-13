// +build acceptance

package acceptance

import (
	"testing"
	mailgun "github.com/mailgun/mailgun-go"
)

func TestGetComplaints(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	publicApiKey := reqEnv(t, "MG_PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	n, complaints, err := mg.GetComplaints(-1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(complaints) != n {
		t.Fatalf("Expected %d complaints; got %d", n, len(complaints))
	}
}
