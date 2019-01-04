package main

import (
	"github.com/davecgh/go-spew/spew"

	"github.com/mailgun/mailgun-go"
	"github.com/thrawn01/args"
)

func ListEvents(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`list events via the mailgun HTTP API

	Examples:
	   list all available tags
	   $ mailgun events

	   limit each event page to 100 events
	   $ mailgun events -l 100`)

	parser.SetDesc(desc)
	parser.AddOption("--limit").Alias("-l").IsInt().Help("Limit the page size")

	opts := parser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	limit := opts.Int("limit")

	// Create the event iterator
	it := mg.ListEvents(&mailgun.EventsOptions{
		Limit: limit,
	})

	var page []mailgun.Event
	for it.Next(&page) {
		for _, event := range page {
			spew.Printf("%+v\n", event)
		}
	}
	if it.Err() != nil {
		return 1, it.Err()
	}
	return 0, nil
}