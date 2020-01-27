package mailgun_test

import (
	"context"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
)

func TestLimits(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()
	limits, err := mg.GetTagLimits(ctx, testDomain)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, limits.Limit, 50000)
	ensure.DeepEqual(t, limits.Count, 5000)
}
