package main

import (
	"fmt"
	"os"

	"github.com/mailgun/log"
	"github.com/mailgun/mailgun-go"
	"github.com/thrawn01/args"
)

func checkErr(msg string, err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s - %s\n", msg, err)
		os.Exit(1)
	}
}

func main() {
	log.InitWithConfig(log.Config{Name: "console"})
	desc := args.Dedent(`CLI for mailgun api

	Examples:
	   Set your credentials in the environment
	   export MG_DOMAIN=your-domain-name
	   export MG_API_KEY=your-api-key
	   export MG_PUBLIC_API_KEY=your-public-api-key

	   Post a simple event message from stdin
	   $ echo -n 'Hello World' | mailgun send -s "Test subject" address@example.com

	 Help:
	   For detailed help on send
	   $ mailgun send -h`)

	parser := args.NewParser(args.EnvPrefix("MG_"), args.Desc(desc, args.IsFormated))
	parser.AddOption("--verbose").Alias("-v").IsTrue().Help("be verbose")
	parser.AddOption("--url").Env("URL").Default(mailgun.ApiBase).Help("url to the mailgun api")
	parser.AddOption("--api-key").Env("API_KEY").Help("mailgun api key")
	parser.AddOption("--public-api-key").Env("PUBLIC_API_KEY").Help("mailgun public api key")
	parser.AddOption("--domain").Env("DOMAIN").Help("mailgun api key")

	// Commands
	parser.AddCommand("send", Send)
	parser.AddCommand("tag", Tag)

	// Parser and set global options
	opts := parser.ParseArgsSimple(nil)
	if opts.Bool("verbose") {
		mailgun.Debug = true
	}

	// Initialize our mailgun object
	mg := mailgun.NewMailgun(
		opts.String("domain"),
		opts.String("api-key"),
		opts.String("public-api-key"))

	// Set our api url
	mg.SetAPIBase(opts.String("url"))

	// Run the command chosen by our user
	retCode, err := parser.RunCommand(mg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(retCode)
}
