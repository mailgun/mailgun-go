package mailgun_test

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
)

func TestRouteCRUD(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	var countRoutes = func() int {
		it := mg.ListRoutes(nil)
		var page []mailgun.Route
		it.Next(ctx, &page)
		ensure.Nil(t, it.Err())
		return it.TotalCount
	}

	routeCount := countRoutes()

	newRoute, err := mg.CreateRoute(ctx, mailgun.Route{
		Priority:    1,
		Description: "Sample Route",
		Expression:  "match_recipient(\".*@samples.mailgun.org\")",
		Actions: []string{
			"forward(\"http://example.com/messages/\")",
			"stop()",
		},
	})
	ensure.Nil(t, err)
	ensure.True(t, newRoute.Id != "")

	defer func() {
		ensure.Nil(t, mg.DeleteRoute(ctx, newRoute.Id))
		_, err = mg.GetRoute(ctx, newRoute.Id)
		ensure.NotNil(t, err)
	}()

	newCount := countRoutes()
	ensure.False(t, newCount <= routeCount)

	theRoute, err := mg.GetRoute(ctx, newRoute.Id)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, newRoute, theRoute)

	changedRoute, err := mg.UpdateRoute(ctx, newRoute.Id, mailgun.Route{
		Priority: 2,
	})
	ensure.Nil(t, err)
	ensure.DeepEqual(t, changedRoute.Priority, 2)
	ensure.DeepEqual(t, len(changedRoute.Actions), 2)
}

func TestRoutesIterator(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	it := mg.ListRoutes(&mailgun.ListOptions{Limit: 2})

	var firstPage, secondPage, previousPage, lastPage []mailgun.Route
	var ctx = context.Background()

	// Calling Last() is invalid unless you first use First() or Next()
	ensure.False(t, it.Last(ctx, &lastPage))
	ensure.True(t, len(lastPage) == 0)

	// Get our first page
	ensure.True(t, it.Next(ctx, &firstPage))
	ensure.Nil(t, it.Err())
	ensure.True(t, len(firstPage) != 0)
	firstIterator := *it

	// Get our second page
	ensure.True(t, it.Next(ctx, &secondPage))
	ensure.Nil(t, it.Err())
	ensure.True(t, len(secondPage) != 0)

	// Pages should be different
	ensure.NotDeepEqual(t, firstPage, secondPage)
	ensure.True(t, firstIterator.TotalCount != 0)

	// Previous()
	ensure.True(t, it.First(ctx, &firstPage))
	ensure.True(t, it.Next(ctx, &secondPage))

	ensure.True(t, it.Previous(ctx, &previousPage))
	ensure.True(t, len(previousPage) != 0)
	ensure.DeepEqual(t, previousPage[0].Id, firstPage[0].Id)

	// First()
	ensure.True(t, it.First(ctx, &firstPage))
	ensure.True(t, len(firstPage) != 0)

	// Calling first resets the iterator to the first page
	ensure.True(t, it.Next(ctx, &secondPage))
	ensure.NotDeepEqual(t, firstPage, secondPage)

	// Last()
	ensure.True(t, it.Last(ctx, &firstPage))
	ensure.True(t, len(firstPage) != 0)

}
