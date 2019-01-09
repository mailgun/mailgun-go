package mailgun_test

import (
	"context"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v3"
)

func ExampleMailgunImpl_ValidateEmail() {
	v := mailgun.NewEmailValidator("my_public_api_key")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ev, err := v.ValidateEmail(ctx, "joe@example.com", false)
	if err != nil {
		log.Fatal(err)
	}
	if !ev.IsValid {
		log.Fatal("Expected valid e-mail address")
	}
	log.Printf("Parts local_part=%s domain=%s display_name=%s", ev.Parts.LocalPart, ev.Parts.Domain, ev.Parts.DisplayName)
	if ev.DidYouMean != "" {
		log.Printf("The address is syntactically valid, but perhaps has a typo.")
		log.Printf("Did you mean %s instead?", ev.DidYouMean)
	}
}

func ExampleMailgunImpl_ParseAddresses() {
	v := mailgun.NewEmailValidator("my_public_api_key")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	addressesThatParsed, unparsableAddresses, err := v.ParseAddresses(ctx, "Alice <alice@example.com>", "bob@example.com", "example.com")
	if err != nil {
		log.Fatal(err)
	}
	hittest := map[string]bool{
		"Alice <alice@example.com>": true,
		"bob@example.com":           true,
	}
	for _, a := range addressesThatParsed {
		if !hittest[a] {
			log.Fatalf("Expected %s to be parsable", a)
		}
	}
	if len(unparsableAddresses) != 1 {
		log.Fatalf("Expected 1 address to be unparsable; got %d", len(unparsableAddresses))
	}
}

func ExampleMailgunImpl_UpdateList() {
	mg := mailgun.NewMailgun("example.com", "my_api_key")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := mg.UpdateMailingList(ctx, "joe-stat@example.com", mailgun.MailingList{
		Name:        "Joe Stat",
		Description: "Joe's status report list",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func ExampleMailgunImpl_Send_constructed() {
	mg := mailgun.NewMailgun("example.com", "my_api_key")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	m := mg.NewMessage(
		"Excited User <me@example.com>",
		"Hello World",
		"Testing some Mailgun Awesomeness!",
		"baz@example.com",
		"bar@example.com",
	)
	m.SetTracking(true)
	m.SetDeliveryTime(time.Now().Add(24 * time.Hour))
	m.SetHtml("<html><body><h1>Testing some Mailgun Awesomeness!!</h1></body></html>")
	_, id, err := mg.Send(ctx, m)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Message id=%s", id)
}

func ExampleMailgunImpl_Send_mime() {
	exampleMime := `Content-Type: text/plain; charset="ascii"
Subject: Joe's Example Subject
From: Joe Example <joe@example.com>
To: BARGLEGARF <bargle.garf@example.com>
Content-Transfer-Encoding: 7bit
Date: Thu, 6 Mar 2014 00:37:52 +0000

Testing some Mailgun MIME awesomeness!
`

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	mg := mailgun.NewMailgun("example.com", "my_api_key")
	m := mg.NewMIMEMessage(ioutil.NopCloser(strings.NewReader(exampleMime)), "bargle.garf@example.com")
	_, id, err := mg.Send(ctx, m)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Message id=%s", id)
}

func ExampleMailgunImpl_GetRoutes() {
	mg := mailgun.NewMailgun("example.com", "my_api_key")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	n, routes, err := mg.ListRoutes(ctx, mailgun.DefaultLimit, mailgun.DefaultSkip)
	if err != nil {
		log.Fatal(err)
	}
	if n > len(routes) {
		log.Printf("More routes exist than has been returned.")
	}
	for _, r := range routes {
		log.Printf("Route pri=%d expr=%s desc=%s", r.Priority, r.Expression, r.Description)
	}
}

func ExampleMailgunImpl_UpdateRoute() {
	mg := mailgun.NewMailgun("example.com", "my_api_key")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err := mg.UpdateRoute(ctx, "route-id-here", mailgun.Route{
		Priority: 2,
	})
	if err != nil {
		log.Fatal(err)
	}
}
