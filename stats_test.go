package mailgun

import (
	"testing"

	"github.com/facebookgo/ensure"
)

func TestGetStats(t *testing.T) {
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)

	totalCount, stats, err := mg.GetStats(-1, -1, nil, "sent", "opened")
	ensure.Nil(t, err)

	t.Logf("Total Count: %d\n", totalCount)
	t.Logf("Id\tEvent\tCreatedAt\tTotalCount\t\n")
	for _, stat := range stats {
		t.Logf("%s\t%s\t%s\t%d\t\n", stat.Id, stat.Event, stat.CreatedAt, stat.TotalCount)
	}
}

func TestDeleteTag(t *testing.T) {
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ensure.Nil(t, mg.DeleteTag("newsletter"))
}
