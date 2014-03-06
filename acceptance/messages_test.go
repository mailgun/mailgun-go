// +build acceptance,spendMoney

package acceptance

import (
	"fmt"
	"github.com/mailgun/mailgun-go"
	"testing"
)

const (
	fromUser       = "Joe Example <joe@example.com>"
	exampleSubject = "Joe's Example Subject"
	exampleText    = "Testing some Mailgun awesomeness!"
	exampleHtml    = "<html><head /><body><p>Testing some Mailgun HTML awesomeness!</p></body></html>"

	exampleMime = `Content-Type: text/plain; charset="ascii"
Subject: Joe's Example Subject
From: Joe Example <joe@example.com>
To: BARGLEGARF <sam.falvo@rackspace.com>
Content-Transfer-Encoding: 7bit
Date: Thu, 6 Mar 2014 00:37:52 +0000

Testing some Mailgun MIME awesomeness!
`
)

func TestSendPlain(t *testing.T) {
	toUser := reqEnv(t, "MG_EMAIL_TO")
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	publicApiKey := reqEnv(t, "MG_PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	msg, id, err := mg.Send(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestSendPlain:MSG(" + msg + "),ID(" + id + ")")
}

func TestSendHtml(t *testing.T) {
	toUser := reqEnv(t, "MG_EMAIL_TO")
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	publicApiKey := reqEnv(t, "MG_PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	m.SetHtml(exampleHtml)
	msg, id, err := mg.Send(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestSendHtml:MSG(" + msg + "),ID(" + id + ")")
}

func TestSendTracking(t *testing.T) {
	toUser := reqEnv(t, "MG_EMAIL_TO")
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	publicApiKey := reqEnv(t, "MG_PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText+"Tracking!\n", toUser)
	m.SetTracking(false)
	msg, id, err := mg.Send(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestSendTracking:MSG(" + msg + "),ID(" + id + ")")
}

func TestSendTag(t *testing.T) {
	toUser := reqEnv(t, "MG_EMAIL_TO")
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	publicApiKey := reqEnv(t, "MG_PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText+"Tags Galore!\n", toUser)
	m.AddTag("FooTag")
	m.AddTag("BarTag")
	m.AddTag("BlortTag")
	msg, id, err := mg.Send(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestSendTag:MSG(" + msg + "),ID(" + id + ")")
}

// TODO(sam-falvo): these require changes to the core library or support for sending multipart/mime-encoded content.

// func TestSendMime(t *testing.T) {
func testSendMime(t *testing.T) {
	t.Fatalf("Not Implemented Yet")
}

// func TestSendDeliveryTime(t *testing.T) {
func testSendDeliveryTime(t *testing.T) {
	t.Fatalf("Not Implemented Yet")
}

