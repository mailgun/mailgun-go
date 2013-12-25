package mailgun

import (
	"testing"
)

const DOMAIN = "banzon.dk"
const API_KEY = "testingapikey"

func TestMailgun(t *testing.T) {
	m := NewMailgun(DOMAIN, API_KEY)
	if API_KEY != m.ApiKey() {
		t.Fatal("ApiKey not equal!")
	}
}
