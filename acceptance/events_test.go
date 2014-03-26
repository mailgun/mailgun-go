package acceptance

import (
	"testing"
	"os"
	"text/tabwriter"
	"github.com/mailgun/mailgun-go"
	"fmt"
)

func TestGetEvents(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, "")
	events, links, err := mg.GetEvents(mailgun.GetEventsOptions{})
	if err != nil {
		t.Fatal(err)
	}
	tw := &tabwriter.Writer{}
	tw.Init(os.Stdout, 2, 8, 2, ' ', tabwriter.AlignRight)
	fmt.Fprintln(tw, "Event\tTimestamp\t")
	for _, event := range events {
		fmt.Fprintf(tw, "%s\t%v\t\n", event["event"], int64(event["timestamp"].(float64)))
	}
	tw.Flush()
	fmt.Printf("%d events dumped\n\n", len(events))
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
