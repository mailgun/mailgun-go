package mailgun_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/mailgun/mailgun-go/v4/events"
)

func ExampleMailgunImpl_ValidateEmail() {
	v := mailgun.NewEmailValidator("my_public_validation_api_key")

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
	v := mailgun.NewEmailValidator("my_public_validation_api_key")

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

	it := mg.ListRoutes(nil)
	var page []mailgun.Route
	for it.Next(ctx, &page) {
		for _, r := range page {
			log.Printf("Route pri=%d expr=%s desc=%s", r.Priority, r.Expression, r.Description)
		}
	}
	if it.Err() != nil {
		log.Fatal(it.Err())
	}
}

func ExampleMailgunImpl_VerifyWebhookSignature() {
	// Create an instance of the Mailgun Client
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		fmt.Printf("mailgun error: %s\n", err)
		os.Exit(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var payload mailgun.WebhookPayload
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			fmt.Printf("decode JSON error: %s", err)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		verified, err := mg.VerifyWebhookSignature(payload.Signature)
		if err != nil {
			fmt.Printf("verify error: %s\n", err)
			w.WriteHeader(http.StatusNotAcceptable)
			return
		}

		if !verified {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Printf("failed verification %+v\n", payload.Signature)
			return
		}

		fmt.Printf("Verified Signature\n")

		// Parse the raw event to extract the

		e, err := mailgun.ParseEvent(payload.EventData)
		if err != nil {
			fmt.Printf("parse event error: %s\n", err)
			return
		}

		switch event := e.(type) {
		case *events.Accepted:
			fmt.Printf("Accepted: auth: %t\n", event.Flags.IsAuthenticated)
		case *events.Delivered:
			fmt.Printf("Delivered transport: %s\n", event.Envelope.Transport)
		}
	})

	fmt.Println("Running...")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		fmt.Printf("serve error: %s\n", err)
		os.Exit(1)
	}
}

func ExampleMailgunImpl_ListMetrics() {
	mg, err := mailgun.NewMailgunFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	opts := mailgun.MetricsOptions{
		Start: mailgun.RFC2822Time(time.Now().Add(-time.Hour * 24 * 30)),
		Pagination: mailgun.MetricsPagination{
			Limit: 10,
		},
	}

	iter, err := mg.ListMetrics(opts)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 2; i++ {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			resp := mailgun.MetricsResponse{}
			more := iter.Next(ctx, &resp)
			if iter.Err() != nil {
				log.Fatal(iter.Err())
			}

			b, _ := json.Marshal(resp)
			log.Printf("%s", b)

			if !more {
				log.Print("no more pages")
				os.Exit(0)
			}
		}()
	}
	// Output:
}
