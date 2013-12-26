package mailgun

import (
	"testing"
)

const domain = "valid-mailgun-domain"
const apiKey = "valid-mailgun-api-key"
const publicApiKey = "valid-mailgun-public-api-key"

func TestMailgun(t *testing.T) {
	m := NewMailgun(domain, apiKey, publicApiKey)

	if domain != m.Domain() {
		t.Fatal("Domain not equal!")
	}

	if apiKey != m.ApiKey() {
		t.Fatal("ApiKey not equal!")
	}

	if publicApiKey != m.PublicApiKey() {
		t.Fatal("PublicApiKey not equal!")
	}
}
