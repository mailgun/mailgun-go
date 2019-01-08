package mailgun

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/facebookgo/ensure"
)

func TestGetComplaints(t *testing.T) {
	reqEnv(t, "MG_PUBLIC_API_KEY")
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	n, complaints, err := mg.ListComplaints(ctx, -1, -1)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, len(complaints), n)
}

func TestGetComplaintFromRandomNoComplaint(t *testing.T) {
	reqEnv(t, "MG_PUBLIC_API_KEY")
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
	mg, err := NewMailgunFromEnv()
	ensure.Nil(t, err)
	ctx := context.Background()

	var hasComplaint = func(email string) bool {
		t.Logf("hasComplaint: %s\n", email)
		_, complaints, err := mg.ListComplaints(ctx, DefaultLimit, DefaultSkip)
		ensure.Nil(t, err)

		for _, complaint := range complaints {
			t.Logf("Complaint Address: %s\n", complaint.Address)
			if complaint.Address == email {
				return true
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
