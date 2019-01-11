package mailgun

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestGetComplaints(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	it := mg.ListComplaints(nil)
	var page []Complaint
	for it.Next(ctx, &page) {
		//spew.Dump(page)
	}
	ensure.Nil(t, err)
}

func TestGetComplaintFromRandomNoComplaint(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	_, err = mg.GetComplaint(ctx, randomString(64, "")+"@example.com")
	ensure.NotNil(t, err)

	ure, ok := err.(*UnexpectedResponseError)
	ensure.True(t, ok)
	ensure.DeepEqual(t, ure.Actual, http.StatusNotFound)
}

func TestCreateDeleteComplaint(t *testing.T) {
	if reason := SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	var hasComplaint = func(email string) bool {
		t.Logf("hasComplaint: %s\n", email)
		it := mg.ListComplaints(nil)
		ensure.Nil(t, err)

		var page []Complaint
		for it.Next(ctx, &page) {
			for _, complaint := range page {
				t.Logf("Complaint Address: %s\n", complaint.Address)
				if complaint.Address == email {
					return true
				}
			}
		}
		return false
	}

	randomMail := strings.ToLower(randomString(64, "")) + "@example.com"
	ensure.False(t, hasComplaint(randomMail))

	ensure.Nil(t, mg.CreateComplaint(ctx, randomMail))
	ensure.True(t, hasComplaint(randomMail))
	ensure.Nil(t, mg.DeleteComplaint(ctx, randomMail))
	ensure.False(t, hasComplaint(randomMail))
}
