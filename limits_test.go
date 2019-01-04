package mailgun_test

import (
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go"
)

func TestLimits(t *testing.T) {
	mg, err := mailgun.NewMailgunFromEnv()
	mg.SetAPIBase(server.URL())
	ensure.Nil(t, err)

	limits, err := mg.GetTagLimits(sampleDomain)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, limits.Limit, 50000)
	ensure.DeepEqual(t, limits.Count, 5000)
}
