package acceptance

import (
	"fmt"
	"github.com/mailgun/mailgun-go"
	"os"
	"testing"
	"text/tabwriter"
)

func TestGetEvents(t *testing.T) {
	// Grab the list of events (as many as we can get)
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, "")
	events, links, err := mg.GetEvents(mailgun.GetEventsOptions{})
	if err != nil {
		t.Fatal(err)
	}

	// Print out the kind of event and timestamp.
	// Specifics about each event will depend on the "event" type.
	tw := &tabwriter.Writer{}
	tw.Init(os.Stdout, 2, 8, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "Event\tTimestamp\t")
	for _, event := range events {
		fmt.Fprintf(tw, "%s\t%v\t\n", event["event"], event["timestamp"])
	}
	tw.Flush()
	fmt.Printf("%d events dumped\n\n", len(events))

	// Print out the types of links provided in case more pages of data exist.
	// For brevity, links are truncated.
	fmt.Fprintln(tw, "Link\tDestination\t")
	for k, v := range links {
		if len(v) > 48 {
			v = fmt.Sprintf("%s...", v[0:48])
		}
		fmt.Fprintf(tw, "%s\t%s\t\n", k, v)
	}
	tw.Flush()
	fmt.Printf("%d links given\n\n", len(links))
}
