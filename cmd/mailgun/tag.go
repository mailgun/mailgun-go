package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mailgun/mailgun-go/v3"
	"github.com/thrawn01/args"
)

func Tag(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`Manage tags via the mailgun HTTP API

	Examples:
	   list all available tags
	   $ mailgun tag list

	   list tags with a specific prefix
	   $ mailgun tag list -p foo

	   get a single tag
	   $ mailgun tag get my-tag

	   delete a tag
	   $ mailgun tag delete my-tag`)

	parser.SetDesc(desc)

	// Commands
	parser.AddCommand("list", ListTag)
	parser.AddCommand("get", GetTag)
	parser.AddCommand("delete", DeleteTag)

	// Run the command chosen by our user
	return parser.ParseAndRun(nil, mg)
}

func ListTag(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`list tags via the mailgun HTTP API

	Examples:
	   list all available tags
	   $ mailgun tag list

	   list the first 2,000 tags
	   $ mailgun tag list -l 2000

	   list tags with a specific prefix
	   $ mailgun tag list -p foo`)
	parser.SetDesc(desc)
	parser.AddOption("--prefix").Alias("-p").Help("list only tags with the given prefix")
	parser.AddOption("--limit").Alias("-l").IsInt().Help("Limit the result set")

	opts := parser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	// Calculate our request limit
	limit := opts.Int("limit")

	// Create the tag iterator
	it := mg.ListTags(&mailgun.ListTagOptions{
		Limit:  limit,
		Prefix: opts.String("prefix"),
	})

	var ctx = context.Background()
	var count int
	var page []mailgun.Tag
	for it.Next(ctx, &page) {
		for _, tag := range page {
			fmt.Printf("%s\n", tag.Value)
			count += 1
			if limit != 0 && count > limit {
				return 0, nil
			}
		}
	}
	if it.Err() != nil {
		return 1, it.Err()
	}
	return 0, nil
}

func GetTag(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`get metatdata about a tag via the mailgun HTTP API

	Examples:
	   fetch the tag metatdata and print it in json
	   $ mailgun tag get my-tag`)
	parser.SetDesc(desc)
	parser.AddArgument("tag").Required().Help("the tag to retrieve")

	opts := parser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	tag, err := mg.GetTag(context.Background(), opts.String("tag"))
	if err != nil {
		return 1, err
	}
	output, err := json.Marshal(tag)
	if err != nil {
		return 1, fmt.Errorf("Json Error: %s\n", err)
	}
	fmt.Print(string(output))
	return 0, nil
}

func DeleteTag(parser *args.ArgParser, data interface{}) (int, error) {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`delete a tag via the mailgun HTTP API

	Examples:
	   delete my-tag
	   $ mailgun tag delete my-tag`)
	parser.SetDesc(desc)
	parser.AddArgument("tag").Required().Help("the tag to delete")

	opts := parser.ParseSimple(nil)
	if opts == nil {
		return 1, nil
	}

	err := mg.DeleteTag(context.Background(), opts.String("tag"))
	if err != nil {
		return 1, err
	}
	return 0, nil
}
