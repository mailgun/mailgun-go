package main

import (
	"fmt"
	"os"

	"encoding/json"

	"github.com/mailgun/log"
	"github.com/mailgun/mailgun-go"
	"github.com/thrawn01/args"
)

func Tag(parser *args.ArgParser, data interface{}) int {
	mg := data.(mailgun.Mailgun)
	var err error

	log.InitWithConfig(log.Config{Name: "console"})
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

	// Parse the subcommands
	parser.ParseArgsSimple(nil)

	// Run the command chosen by our user
	retCode, err := parser.RunCommand(mg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return 1
	}
	return retCode
}

func ListTag(parser *args.ArgParser, data interface{}) int {
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
	parser.AddOption("--tag").Alias("-t").Help("The tag that marks piviot point for the --page parameter")
	parser.AddOption("--page").Alias("-pg").
		Help("The page direction based off the tag parameter; valid choices are (first, last, next, prev)")

	opts := parser.ParseArgsSimple(nil)

	// Calculate our request limit
	limit := opts.Int("limit")

	// Create the tag iterator
	it := mg.ListTags(&mailgun.TagOptions{
		Limit:  limit,
		Prefix: opts.String("prefix"),
		Page:   opts.String("page"),
		Tag:    opts.String("tag"),
	})

	var count int
	var page mailgun.TagsPage
	for it.Next(&page) {
		for _, tag := range page.Items {
			fmt.Printf("%s\n", tag.Value)
			count += 1
			if limit != 0 && count > limit {
				return 0
			}
		}
	}
	return 0
}

func GetTag(parser *args.ArgParser, data interface{}) int {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`get metatdata about a tag via the mailgun HTTP API

	Examples:
	   fetch the tag metatdata and print it in json
	   $ mailgun tag get my-tag`)
	parser.SetDesc(desc)
	parser.AddArgument("tag").Required().Help("the tag to retrieve")

	opts := parser.ParseArgsSimple(nil)

	tag, err := mg.GetTag(opts.String("tag"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	output, err := json.Marshal(tag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Json Error: %s\n", err)
		return 1
	}
	fmt.Print(string(output))
	return 0
}

func DeleteTag(parser *args.ArgParser, data interface{}) int {
	mg := data.(mailgun.Mailgun)

	desc := args.Dedent(`delete a tag via the mailgun HTTP API

	Examples:
	   delete my-tag
	   $ mailgun tag delete my-tag`)
	parser.SetDesc(desc)
	parser.AddArgument("tag").Required().Help("the tag to delete")

	opts := parser.ParseArgsSimple(nil)

	err := mg.DeleteTag(opts.String("tag"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		return 1
	}
	return 0
}
