package mailgun

import (
	"testing"
)

const API_KEY = "testingapikey"

func TestMailgun(t *testing.T) {
	m := NewMailgun(API_KEY)
	if API_KEY != m.ApiKey() {
		t.Fatal("ApiKey not equal!")
	}
}
