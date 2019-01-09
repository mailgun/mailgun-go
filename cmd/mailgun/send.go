package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/drhodes/golorem"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/pkg/errors"
	"github.com/thrawn01/args"
)

func Send(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)
	var content []byte
	var err error
	var count int

	desc := args.Dedent(`Send emails via the mailgun HTTP API

	Examples:
	   Post a simple email from stdin
	   $ echo -n 'Hello World' | mailgun send -s "Test subject" address@example.com

	   Post a simple email to a specific domain
	   $ echo -n 'Hello World' | mailgun send -s "Test subject" address@example.com -d my-domain.com

	   Post a test lorem ipsum email (random content, and subject)
	   $ mailgun send --lorem address@example.com

	   Post a 10 random test lorem ipsum emails
	   $ mailgun send --lorem address@example.com --count 10`)

	parser.SetDesc(desc)
	parser.AddOption("--subject").Alias("-s").Help("subject of the message")
	parser.AddOption("--tags").IsStringSlice().Alias("-t").Help("comma separated list of tags")
	parser.AddOption("--from").Alias("-f").Env("FROM").Help("from address, defaults to <user>@<hostname>")
	parser.AddOption("--lorem").Alias("-l").IsTrue().Help("generate a random subject and message content")
	parser.AddOption("--count").StoreInt(&count).Default("1").Alias("-c").Help("send the email x number of counts")
	parser.AddArgument("addresses").IsStringSlice().Required().Help("a list of email addresses")

	opts := parser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	// Required for send
	if err := opts.Required([]string{"domain", "api-key"}); err != nil {
		return 1, fmt.Errorf("missing Required option '%s'", err)
	}

	// Default to user@hostname if no from address provided
	if !opts.IsSet("from") {
		host, err := os.Hostname()
		if err != nil {
			return 1, errors.Wrapf(err, "during hostname lookup")
		}
		opts.Set("from", fmt.Sprintf("%s@%s", os.Getenv("USER"), host))
	}

	// If stdin is not open and character device
	if args.IsCharDevice(os.Stdin) {
		// Read the content from stdin
		if content, err = ioutil.ReadAll(os.Stdin); err != nil {
			return 1, errors.Wrap(err, "while reading from stdin")
		}
	}

	subject := opts.String("subject")

	if opts.Bool("lorem") {
		if len(subject) == 0 {
			subject = lorem.Sentence(3, 5)
		}
		if len(content) == 0 {
			content = []byte(lorem.Paragraph(10, 50))
		}
	} else {
		if len(content) == 0 {
			return 1, fmt.Errorf("must provide email body, or use --lorem")
		}
		if len(subject) == 0 {
			return 1, fmt.Errorf("must provide subject, or use --lorem")
		}
	}

	var tags []string
	if opts.IsSet("tags") {
		tags = opts.StringSlice("tags")
	}

	for i := 0; i < count; i++ {
		msg := mg.NewMessage(
			opts.String("from"),
			subject,
			string(content),
			opts.StringSlice("addresses")...)

		// Add any tags if provided
		for _, tag := range tags {
			msg.AddTag(tag)
		}

		ctx := context.Background()
		resp, id, err := mg.Send(ctx, msg)
		if err != nil {
			return 1, errors.Wrap(err, "while sending message")
		}
		fmt.Printf("Id: %s Resp: %s\n", id, resp)
	}
	return 0, nil
}
