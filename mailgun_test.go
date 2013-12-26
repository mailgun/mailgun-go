package mailgun

import (
	"testing"
)

const DOMAIN = "valid-mailgun-domain"
const API_KEY = "valid-mailgun-api-key"

func TestMailgun(t *testing.T) {
	m := NewMailgun(DOMAIN, API_KEY)
	if API_KEY != m.ApiKey() {
		t.Fatal("ApiKey not equal!")
	}
	if DOMAIN != m.Domain() {
		t.Fatal("Domain not equal!")
	}
}
