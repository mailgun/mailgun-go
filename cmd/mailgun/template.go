package main

import (
	"context"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/mailgun/mailgun-go/v3"
	"github.com/thrawn01/args"
)

func Templates(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`Manage templates via the mailgun HTTP API

	Examples:
	   list all available templates
	   $ mailgun templates list`)

	parser.SetDesc(desc)

	// Commands
	parser.AddCommand("list", ListTemplates)

	// Run the command chosen by our user
	return parser.ParseAndRun(nil, mg)
}

func ListTemplates(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`list templates via the mailgun HTTP API

	Examples:
	   list all available tags
	   $ mailgun templates list`)

	parser.SetDesc(desc)
	parser.AddOption("--limit").Alias("-l").IsInt().Help("Limit the page size")

	opts := parser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	limit := opts.Int("limit")

	// Create the event iterator
	it := mg.ListTemplates(&mailgun.ListOptions{
		Limit: limit,
	})

	var page []mailgun.Template
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	for it.Next(ctx, &page) {
		for _, event := range page {
			spew.Printf("%+v\n", event)
		}
	}
	cancel()
	if it.Err() != nil {
		return 1, it.Err()
	}
	return 0, nil
}
