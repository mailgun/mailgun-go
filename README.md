# Mailgun with Go

[![GoDoc](https://godoc.org/github.com/mailgun/mailgun-go?status.svg)](https://godoc.org/github.com/mailgun/mailgun-go)
[![Build Status](https://img.shields.io/travis/mailgun/mailgun-go/master.svg)](https://travis-ci.org/mailgun/mailgun-go)

Go library for interacting with the [Mailgun](https://mailgun.com/) [API](https://documentation.mailgun.com/api_reference.html).

**NOTE: Backward compatibility has been broken with the v3.0 release which includes versioned paths required by 1.11 
go modules (See [Releasing Modules](https://github.com/golang/go/wiki/Modules#releasing-modules-v2-or-higher)).
 Pin your dependencies to the v1.1.1 or v2.0 tag if you are not ready for v3.0**

## Usage
```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/mailgun/mailgun-go/v3"
)

// Your available domain names can be found here:
// (https://app.mailgun.com/app/domains)
var yourDomain string = "your-domain-name" // e.g. mg.yourcompany.com

// You can find the Private API Key in your Account Menu, under "Settings":
// (https://app.mailgun.com/app/account/security)
var privateAPIKey string = "your-private-key"


func main() {
    // Create an instance of the Mailgun Client
    mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

    sender := "sender@example.com"
    subject := "Fancy subject!"
    body := "Hello from Mailgun Go!"
    recipient := "recipient@example.com"

    // The message object allows you to add attachments and Bcc recipients
    message := mg.NewMessage(sender, subject, body, recipient)

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    // Send the message	with a 10 second timeout
    resp, id, err := mg.Send(ctx, message)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
```

## Get Events
```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/mailgun/mailgun-go/v3"
    "github.com/mailgun/mailgun-go/v3/events"
)

func main() {
    // You can find the Private API Key in your Account Menu, under "Settings":
    // (https://app.mailgun.com/app/account/security)
    mg := mailgun.NewMailgun("your-domain.com", "your-private-key")

	it := mg.ListEvents(&mailgun.ListEventOptions{Limit: 100})

	var page []mailgun.Event

	// The entire operation should not take longer than 30 seconds
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// For each page of 100 events
	for it.Next(ctx, &page) {
		for _, e := range page {
			// You can access some fields via the interface
			fmt.Printf("Event: '%s' TimeStamp: '%s'\n", e.GetName(), e.GetTimestamp())

			// and you can act upon each event by type
			switch event := e.(type) {
			case *events.Accepted:
				fmt.Printf("Accepted: auth: %t\n", event.Flags.IsAuthenticated)
			case *events.Delivered:
				fmt.Printf("Delivered transport: %s\n", event.Envelope.Transport)
			case *events.Failed:
				fmt.Printf("Failed reason: %s\n", event.Reason)
			case *events.Clicked:
				fmt.Printf("Clicked GeoLocation: %s\n", event.GeoLocation.Country)
			case *events.Opened:
				fmt.Printf("Opened GeoLocation: %s\n", event.GeoLocation.Country)
			case *events.Rejected:
				fmt.Printf("Rejected reason: %s\n", event.Reject.Reason)
			case *events.Stored:
				fmt.Printf("Stored URL: %s\n", event.Storage.URL)
			case *events.Unsubscribed:
				fmt.Printf("Unsubscribed client OS: %s\n", event.ClientInfo.ClientOS)
			}
		}
	}
}
```

## Event Polling
The mailgun library has built-in support for polling the events api
```go
package main

import (
    "context"
    "time"

    "github.com/mailgun/mailgun-go/v3"
)

func main() {
    // You can find the Private API Key in your Account Menu, under "Settings":
    // (https://app.mailgun.com/app/account/security)
    mg := mailgun.NewMailgun("your-domain.com", "your-private-key")

	begin := time.Now().Add(time.Second * -3)

	// Very short poll interval
	it := mg.PollEvents(&mailgun.ListEventOptions{
		// Only events with a timestamp after this date/time will be returned
		Begin: &begin,
		// How often we poll the api for new events
		PollInterval: time.Second * 30,
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

    // Poll until our email event arrives
    var page []mailgun.Event
    for it.Poll(ctx, &page) {
        for _, e := range page {
            // Do something with event
        }
    }
}
```

# Email Validations
```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/mailgun/mailgun-go/v3"
)

// If your plan does not include email validations but you have an account,
// use your Public Validation api key. If your plan does include email validations,
// use your Private API key. You can find both the Private and
// Public Validation API Keys in your Account Menu, under "Settings":
// (https://app.mailgun.com/app/account/security)
var apiKey string = "your-api-key"

func main() {
    // Create an instance of the Validator
    v := mailgun.NewEmailValidator(apiKey)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

    email, err := v.ValidateEmail(ctx, "recipient@example.com", false)
    if err != nil {
        panic(err)
    }
    fmt.Printf("Valid: %t\n", email.IsValid)
}
```

## Webhook Handling
```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"

    "github.com/mailgun/mailgun-go/v3"
    "github.com/mailgun/mailgun-go/v3/events"
)

func main() {

    // You can find the Private API Key in your Account Menu, under "Settings":
    // (https://app.mailgun.com/app/account/security)
    mg := mailgun.NewMailgun("your-domain.com", "private-api-key")

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

		// Parse the event provided by the webhook payload
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

	fmt.Println("Serve on :9090...")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		fmt.Printf("serve error: %s\n", err)
		os.Exit(1)
	}
}
```

## Using Templates

Templates enable you to create message templates on your Mailgun account and then populate the data variables at send-time. This allows you to have your layout and design managed on the server and handle the data on the client. The template variables are added as a JSON stringified `X-Mailgun-Variables` header. For example, if you have a template to send a password reset link, you could do the following:

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/mailgun/mailgun-go/v3"
)

// Your available domain names can be found here:
// (https://app.mailgun.com/app/domains)
var yourDomain string = "your-domain-name" // e.g. mg.yourcompany.com

// You can find the Private API Key in your Account Menu, under "Settings":
// (https://app.mailgun.com/app/account/security)
var privateAPIKey string = "your-private-key"


func main() {
    // Create an instance of the Mailgun Client
    mg := mailgun.NewMailgun(yourDomain, privateAPIKey)

    sender := "sender@example.com"
    subject := "Fancy subject!"
    body := ""
    recipient := "recipient@example.com"

    // The message object allows you to add attachments and Bcc recipients
		message := mg.NewMessage(sender, subject, body, recipient)
		message.SetTemplate("passwordReset")
		message.AddTemplateVariable("passwordResetLink", "some link to your site unique to your user")

    ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
    defer cancel()

    // Send the message	with a 10 second timeout
    resp, id, err := mg.Send(ctx, message)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
```

The official mailgun documentation includes examples using this library. Go
[here](https://documentation.mailgun.com/en/latest/api_reference.html#api-reference)
and click on the "Go" button at the top of the page.

### EU Region
European customers will need to change the default API Base to access your domains

```go
mg := mailgun.NewMailgun("your-domain.com", "private-api-key")
mg.SetAPIBase(mailgun.APIBaseEU)
```
## Installation

If you are using [golang modules](https://github.com/golang/go/wiki/Modules) make sure you
include the `/v3` at the end of your import paths
```bash
$ go get github.com/mailgun/mailgun-go/v3
```

If you are **not** using golang modules, you can drop the `/v3` at the end of the import path.
As long as you are using the latest 1.10 or 1.11 golang release, import paths that end in `/v3`
in your code should work fine even if you do not have golang modules enabled for your project.
```bash
$ go get github.com/mailgun/mailgun-go
```

**NOTE for go dep users**

Using version 3 of the mailgun-go library with go dep currently results in the following error
```
"github.com/mailgun/mailgun-go/v3/events", which contains malformed code: no package exists at ...
```
This is a known bug in go dep. You can follow the PR to fix this bug [here](https://github.com/golang/dep/pull/1963)
Until this bug is fixed, the best way to use version 3 of the mailgun-go library is to use the golang community 
supported [golang modules](https://github.com/golang/go/wiki/Modules).

## Testing

*WARNING* - running the tests will cost you money!

To run the tests various environment variables must be set. These are:

* `MG_DOMAIN` is the domain name - this is a value registered in the Mailgun admin interface.
* `MG_PUBLIC_API_KEY` is the Public Validation API key - you can get this value from the Mailgun [security page](https://app.mailgun.com/app/account/security)
* `MG_API_KEY` is the Private API key - you can get this value from the Mailgun [security page](https://app.mailgun.com/app/account/security)
* `MG_EMAIL_TO` is the email address used in various sending tests.

and finally

* `MG_SPEND_MONEY` if this value is set the part of the test that use the API to actually send email will be run - be aware *this will count on your quota* and *this _will_ cost you money*.

The code is released under a 3-clause BSD license. See the LICENSE file for more information.
