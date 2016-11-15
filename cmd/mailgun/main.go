package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/drhodes/golorem"
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
	parser.AddOption("--url").Alias("-u").Env("URL").Default(mailgun.BaseURL).Help("url to the mailgun api")
	parser.AddOption("--api-key").Alias("-a").Env("API_KEY").Help("mailgun api key")
	parser.AddOption("--public-api-key").Alias("-p").Env("PUBLIC_API_KEY").Help("mailgun public api key")
	parser.AddOption("--domain").Alias("-d").Env("DOMAIN").Help("mailgun api key")

	// Commands
	parser.AddCommand("send", Send)

	// Parser and set global options
	opts := parser.ParseArgsSimple(nil)
	mailgun.BaseURL = opts.String("url")
	if opts.Bool("verbose") {
		mailgun.Debug = true
	}

	// Initialize our mailgun object
	mg := mailgun.NewMailgun(
		opts.String("domain"),
		opts.String("api-key"),
		opts.String("public-api-key"))

	// Run the command chosen by our user
	retCode, err := parser.RunCommand(mg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	os.Exit(retCode)
}

func Send(parser *args.ArgParser, data interface{}) int {
	mg := data.(mailgun.Mailgun)
	var content []byte
	var err error
	var count int

	log.InitWithConfig(log.Config{Name: "console"})
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
	parser.AddOption("--lorem").Alias("-l").IsTrue().Help("generate a randome subject and message content")
	parser.AddOption("--count").StoreInt(&count).Default("1").Alias("-c").Help("send the email x number of counts")
	parser.AddArgument("addresses").IsStringSlice().Required().Help("a list of email addresses")

	opts := parser.ParseArgsSimple(nil)

	// Required for send
	if err := opts.Required([]string{"domain", "api-key"}); err != nil {
		fmt.Fprintf(os.Stderr, "Missing Required option '%s'", err)
		return 1
	}

	// Default to user@hostname if no from address provided
	if !opts.IsSet("from") {
		host, err := os.Hostname()
		checkErr("Hostname Error", err)
		opts.Set("from", fmt.Sprintf("%s@%s", os.Getenv("USER"), host))
	}

	// Read the content from stdin
	if !opts.Bool("lorem") {
		// If stdin is not open and character device
		if !args.IsCharDevice(os.Stdin) {
			parser.PrintHelp()
			return 1
		}
		content, err = ioutil.ReadAll(os.Stdin)
		checkErr("Error reading stdin", err)
	}

	subject := opts.String("subject")
	for i := 0; i < count; i++ {
		if opts.Bool("lorem") {
			subject = lorem.Sentence(3, 5)
			content = []byte(lorem.Paragraph(10, 50))
		}
		resp, id, err := mg.Send(mg.NewMessage(
			opts.String("from"),
			subject,
			string(content),
			opts.StringSlice("addresses")...))
		checkErr("Message Error", err)
		fmt.Printf("Id: %s Resp: %s\n", id, resp)
	}
	return 0
}
