// +build acceptance,spendMoney

package acceptance

import (
	"fmt"
	mailgun "github.com/mailgun/mailgun-go"
	"io/ioutil"
	"strings"
	"testing"
	"time"
	"text/tabwriter"
	"os"
)

const (
	fromUser       = "\"Joe, Example\" <joe@example.com>"
	exampleSubject = "Joe's Example Subject"
	exampleText    = "Testing some Mailgun awesomeness!"
	exampleHtml    = "<html><head /><body><p>Testing some <a href=\"http://google.com?q=abc&r=def&s=ghi\">Mailgun HTML awesomeness!</a> at www.kc5tja@yahoo.com</p></body></html>"

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

func TestSendPlainAt(t *testing.T) {
	toUser := reqEnv(t, "MG_EMAIL_TO")
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	publicApiKey := reqEnv(t, "MG_PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	m.SetDeliveryTime(time.Now().Add(5 * time.Minute))
	msg, id, err := mg.Send(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestSendPlainAt:MSG(" + msg + "),ID(" + id + ")")
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

func TestSendMIME(t *testing.T) {
	toUser := reqEnv(t, "MG_EMAIL_TO")
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, "")
	m := mailgun.NewMIMEMessage(ioutil.NopCloser(strings.NewReader(exampleMime)), toUser)
	msg, id, err := mg.Send(m)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("TestSendMIME:MSG(" + msg + "),ID(" + id + ")")
}

func TestGetStoredMessage(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, "")
	msgs, err := mg.GetStoredMessages()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Number of messages: ", len(msgs))
	tw := &tabwriter.Writer{}
	tw.Init(os.Stdout, 2, 8, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "From\tTo\tSubject\t# Attachments\t")
	for _, m := range msgs {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%d\t\n", m.From, m.Recipients, m.Subject, len(m.Attachments))
	}
	tw.Flush()
}
