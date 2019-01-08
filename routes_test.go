package mailgun

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestRouteCRUD(t *testing.T) {
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	var countRoutes = func() int {
		count, _, err := mg.ListRoutes(ctx, DefaultLimit, DefaultSkip)
		ensure.Nil(t, err)
		return count
	}

	routeCount := countRoutes()

	newRoute, err := mg.CreateRoute(ctx, Route{
		Priority:    1,
		Description: "Sample Route",
		Expression:  "match_recipient(\".*@samples.mailgun.org\")",
		Actions: []string{
			"forward(\"http://example.com/messages/\")",
			"stop()",
		},
	})
	ensure.Nil(t, err)
	ensure.True(t, newRoute.ID != "")

	defer func() {
		ensure.Nil(t, mg.DeleteRoute(ctx, newRoute.ID))
		_, err = mg.GetRoute(ctx, newRoute.ID)
		ensure.NotNil(t, err)
	}()

	newCount := countRoutes()
	ensure.False(t, newCount <= routeCount)

	theRoute, err := mg.GetRoute(ctx, newRoute.ID)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, newRoute, theRoute)

	changedRoute, err := mg.UpdateRoute(ctx, newRoute.ID, Route{
		Priority: 2,
	})
	ensure.Nil(t, err)
	ensure.DeepEqual(t, changedRoute.Priority, 2)
	ensure.DeepEqual(t, len(changedRoute.Actions), 2)
}
