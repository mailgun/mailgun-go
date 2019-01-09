package main

import (
	"context"

	"github.com/davecgh/go-spew/spew"

	"github.com/mailgun/mailgun-go/v3"
	"github.com/thrawn01/args"
)

func MailingLists(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`Manage mailing-lists via the mailgun HTTP API

	Examples:
	   list all available mailing lists
	   $ mailgun mailing-lists list`)

	parser.SetDesc(desc)

	// Commands
	parser.AddCommand("list", ListMailingList)
	parser.AddCommand("members-list", ListMailingListMembers)

	// Run the command chosen by our user
	return parser.ParseAndRun(nil, mg)
}

func ListMailingList(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`list mailing lists via the mailgun HTTP API

	Examples:
	   list all available mailing lists
	   $ mailgun mailing-lists list`)

	parser.SetDesc(desc)
	parser.AddOption("--limit").Alias("-l").IsInt().Help("Limit the result set")

	opts := parser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	// Create the tag iterator
	it := mg.ListMailingLists(&mailgun.ListOptions{
		Limit: opts.Int("limit"),
	})

	ctx := context.Background()
	var page []mailgun.MailingList
	for it.Next(ctx, &page) {
		for _, list := range page {
			spew.Printf("%+v\n", list)
		}
	}
	if it.Err() != nil {
		return 1, it.Err()
	}
	return 0, nil
}

func ListMailingListMembers(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`Manage mailing-list members via the mailgun HTTP API

	Examples:
	   list all available addresses in a mailing list
	   $ mailgun mailing-lists members-list my-list@my-domain.com`)

	parser.SetDesc(desc)
	parser.AddOption("--limit").Alias("-l").IsInt().Help("Limit the result set")
	parser.AddArgument("address").Required().Help("The mailing list address")

	opts := parser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	// Create the tag iterator
	it := mg.ListMembers(opts.String("address"), &mailgun.ListOptions{
		Limit: opts.Int("limit"),
	})

	ctx := context.Background()
	var page []mailgun.Member
	for it.Next(ctx, &page) {
		for _, list := range page {
			spew.Printf("%+v\n", list)
		}
	}
	if it.Err() != nil {
		return 1, it.Err()
	}
	return 0, nil
}
