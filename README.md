# Mailgun with Go

[![GoDoc](https://godoc.org/gopkg.in/mailgun/mailgun-go.v1?status.svg)](https://godoc.org/gopkg.in/mailgun/mailgun-go.v1)

Go library for interacting with the [Mailgun](https://mailgun.com/) [API](https://documentation.mailgun.com/api_reference.html).

**NOTE: Backward compatibility has been broken with the v3.0 release which includes versioned paths required by 1.11 
go modules (See [Releasing Modules](https://github.com/golang/go/wiki/Modules#releasing-modules-v2-or-higher)).
 Pin your dependencies to the v1.1.1 or v2.0 tag if you are not ready for v3.0**

## Sending mail via the mailgun CLI

Export your API keys and domain

```bash
$ export MG_API_KEY=your-api-key
$ export MG_DOMAIN=your-domain
$ export MG_PUBLIC_API_KEY=your-public-key
$ export MG_URL="https://api.mailgun.net/v3"
```

Send the email

```bash
$ echo -n 'Hello World' | mailgun send -s "Test subject" address@example.com
```

## Sending mail via the golang library
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

// The API Keys are found in your Account Menu, under "Settings":
// (https://app.mailgun.com/app/account/security)

// starts with "key-"
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
    mg := mailgun.NewMailgun("your-domain.com", "your-api-key")

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
    mg := mailgun.NewMailgun("your-domain.com", "your-api-key")

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
// use your public api key (starts with "pubkey-"). If your plan does include
// email validations, use your private api key (starts with "key-")
var apiKey string = "your-key"

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

The official mailgun documentation includes examples using this library. Go
[here](https://documentation.mailgun.com/en/latest/api_reference.html#api-reference)
and click on the "Go" button at the top of the page.

## Installation

Install the go library

```bash
$ go get github.com/mailgun/mailgun-go/v3
```

Install the mailgun CLI

```bash
$ go install github.com/mailgun/mailgun-go/v3/cmd/mailgun/./...
```

## Testing

*WARNING* - running the tests will cost you money!

To run the tests various environment variables must be set. These are:

* `MG_DOMAIN` is the domain name - this is a value registered in the Mailgun admin interface.
* `MG_PUBLIC_API_KEY` is the public API key - you can get this value from the Mailgun admin interface.
* `MG_API_KEY` is the (private) API key - you can get this value from the Mailgun admin interface.
* `MG_EMAIL_TO` is the email address used in various sending tests.

and finally

* `MG_SPEND_MONEY` if this value is set the part of the test that use the API to actually send email will be run - be aware *this will count on your quota* and *this _will_ cost you money*.

The code is released under a 3-clause BSD license. See the LICENSE file for more information.
