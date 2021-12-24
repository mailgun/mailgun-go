package mailgun_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
)

func TestGetComplaints(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	it := mg.ListComplaints(nil)
	var page []mailgun.Complaint
	for it.Next(ctx, &page) {
		//spew.Dump(page)
	}
	ensure.Nil(t, it.Err())
}

func TestGetComplaintFromRandomNoComplaint(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	_, err := mg.GetComplaint(ctx, randomString(64, "")+"@example.com")
	ensure.NotNil(t, err)

	ure, ok := err.(*mailgun.UnexpectedResponseError)
	ensure.True(t, ok)
	ensure.DeepEqual(t, ure.Actual, http.StatusNotFound)
}

func TestCreateDeleteComplaint(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	var hasComplaint = func(email string) bool {
		t.Logf("hasComplaint: %s\n", email)
		it := mg.ListComplaints(nil)
		ensure.Nil(t, it.Err())

		var page []mailgun.Complaint
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
