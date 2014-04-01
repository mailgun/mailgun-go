// +build acceptance

package acceptance

import (
	"testing"
	"github.com/mailgun/mailgun-go"
)

func TestRouteCRUD(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, "")

	var countRoutes = func() int {
		count, _, err := mg.GetRoutes(mailgun.DefaultLimit, mailgun.DefaultSkip)
		if err != nil {
			t.Fatal(err)
		}
		return count
	}

	routeCount := countRoutes()

	_, err := mg.CreateRoute(mailgun.Route{
		Priority: 1,
		Description: "Sample Route",
		Expression: "match_recipient(\".*@samples.mailgun.org\")",
		Actions: []string{
			"forward(\"http://example.com/messages/\")",
			"stop()",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	newCount := countRoutes()
	if newCount <= routeCount {
		t.Fatalf("Expected %d routes defined; got %d", routeCount+1, newCount)
	}

	// err := mg.DeleteRoute(newRoute.ID)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// newCount = countRoutes()
	// if newCount != routeCount {
	// 	t.Fatalf("Expected %d routes defined; got %d", routeCount, newCount)
	// }
}
