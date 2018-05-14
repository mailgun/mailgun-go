# Mailgun with Go

===============

[![Build Status](https://img.shields.io/travis/forrest321/mailgun-go/master.svg)](https://travis-ci.org/forrest321/mailgun-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/forrest321/mailgun-go)](https://goreportcard.com/report/github.com/forrest321/mailgun-go)
[![GoDoc](https://godoc.org/gopkg.in/mailgun/mailgun-go.v1?status.svg)](https://godoc.org/gopkg.in/mailgun/mailgun-go.v1)

Go library for interacting with the [Mailgun](https://mailgun.com/) [API](https://documentation.mailgun.com/api_reference.html).

## Sending mail via the mailgun CLI

Export your API keys and domain

```bash

$> export MG_API_KEY=your-api-key
$> export MG_DOMAIN=your-domain
$> export MG_PUBLIC_API_KEY=your-public-key
$> export MG_URL="https://api.mailgun.net/v3"

```

Send an email

```bash

$> echo -n 'Hello World' | mailgun send -s "Test subject" address@example.com
```

## Sending mail via the golang library

```go

package main

import (
    "fmt"
    "gopkg.in/mailgun/mailgun-go.v1"
    "log"
)

// Your available domain names can be found here:
// (https://app.mailgun.com/app/domains)
var yourDomain string = "your-domain-name" // e.g. mg.yourcompany.com

// The API Keys are found in your Account Menu, under "Settings":
// (https://app.mailgun.com/app/account/security)

// starts with "key-"
var privateAPIKey string = "your-private-key"

// starts with "pubkey-"
var publicValidationKey string = "your-public-key"

func main() {
    // Create an instance of the Mailgun Client
    mg := mailgun.NewMailgun(yourDomain, privateAPIKey, publicValidationKey)

    sender := "sender@example.com"
    subject := "Fancy subject!"
    body := "Hello from Mailgun Go!"
    recipient := "recipient@example.com"

    sendMessage(mg, sender, subject, body, recipient)
}

func sendMessage(mg mailgun.Mailgun, sender, subject, body, recipient string) {
    message := mg.NewMessage(sender, subject, body, recipient)
    resp, id, err := mg.Send(message)

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("ID: %s Resp: %s\n", id, resp)
}
```

## Installation

Install the go library

```bash
$> go get gopkg.in/mailgun/mailgun-go.v1
```

Install the mailgun CLI

```bash

$> go install github.com/mailgun/mailgun-go/cmd/mailgun/./...

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
