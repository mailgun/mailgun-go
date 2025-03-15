package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouteCRUD(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	var countRoutes = func() int {
		it := mg.ListRoutes(nil)
		var page []mtypes.Route
		it.Next(ctx, &page)
		require.NoError(t, it.Err())
		return it.TotalCount
	}

	routeCount := countRoutes()

	newRoute, err := mg.CreateRoute(ctx, mtypes.Route{
		Priority:    1,
		Description: "Sample Route",
		Expression:  "match_recipient(\".*@samples.mailgun.org\")",
		Actions: []string{
			"forward(\"http://example.com/messages/\")",
			"stop()",
		},
	})
	require.NoError(t, err)
	require.NotEqual(t, "", newRoute.Id)

	defer func() {
		require.NoError(t, mg.DeleteRoute(ctx, newRoute.Id))
		_, err = mg.GetRoute(ctx, newRoute.Id)
		require.NotNil(t, err)
	}()

	newCount := countRoutes()
	require.False(t, newCount <= routeCount)

	theRoute, err := mg.GetRoute(ctx, newRoute.Id)
	require.NoError(t, err)
	assert.Equal(t, newRoute, theRoute)

	changedRoute, err := mg.UpdateRoute(ctx, newRoute.Id, mtypes.Route{
		Priority: 2,
	})
	require.NoError(t, err)
	assert.Equal(t, 2, changedRoute.Priority)
	assert.Len(t, changedRoute.Actions, 2)
}

func TestRoutesIterator(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	it := mg.ListRoutes(&mailgun.ListOptions{Limit: 2})

	var firstPage, secondPage, previousPage, lastPage []mtypes.Route
	var ctx = context.Background()

	// Calling Last() is invalid unless you first use First() or Next()
	require.False(t, it.Last(ctx, &lastPage))
	require.Len(t, lastPage, 0)

	// Get our first page
	require.True(t, it.Next(ctx, &firstPage))
	require.NoError(t, it.Err())
	require.True(t, len(firstPage) != 0)
	firstIterator := *it

	// Get our second page
	require.True(t, it.Next(ctx, &secondPage))
	require.NoError(t, it.Err())
	require.True(t, len(secondPage) != 0)

	// Pages should be different
	require.NotEqual(t, firstPage, secondPage)
	require.True(t, firstIterator.TotalCount != 0)

	// Previous()
	require.True(t, it.First(ctx, &firstPage))
	require.True(t, it.Next(ctx, &secondPage))

	require.True(t, it.Previous(ctx, &previousPage))
	require.True(t, len(previousPage) != 0)
	require.Equal(t, previousPage[0].Id, firstPage[0].Id)

	// First()
	require.True(t, it.First(ctx, &firstPage))
	require.True(t, len(firstPage) != 0)

	// Calling first resets the iterator to the first page
	require.True(t, it.Next(ctx, &secondPage))
	require.NotEqual(t, firstPage, secondPage)

	// Last()
	require.True(t, it.Last(ctx, &firstPage))
	require.True(t, len(firstPage) != 0)
}
