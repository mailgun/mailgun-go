package mailgun_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

func TestGetComplaints(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL())

	ctx := context.Background()

	it := mg.ListComplaints(testDomain, nil)
	var page []mailgun.Complaint
	for it.Next(ctx, &page) {
	}
	require.NoError(t, it.Err())
}

func TestGetComplaintFromRandomNoComplaint(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	_, err := mg.GetComplaint(ctx, testDomain, randomString(64, "")+"@example.com")
	require.NotNil(t, err)

	var ure *mailgun.UnexpectedResponseError
	require.ErrorAs(t, err, &ure)
	require.Equal(t, http.StatusNotFound, ure.Actual)
}

func TestCreateDeleteComplaint(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	var hasComplaint = func(email string) bool {
		t.Logf("hasComplaint: %s\n", email)
		it := mg.ListComplaints(testDomain, nil)
		require.NoError(t, it.Err())

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
	require.False(t, hasComplaint(randomMail))

	require.NoError(t, mg.CreateComplaint(ctx, testDomain, randomMail))
	require.True(t, hasComplaint(randomMail))
	require.NoError(t, mg.DeleteComplaint(ctx, testDomain, randomMail))
	require.False(t, hasComplaint(randomMail))
}

func TestCreateDeleteComplaintList(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	var hasComplaint = func(email string) bool {
		t.Logf("hasComplaint: %s\n", email)
		it := mg.ListComplaints(testDomain, nil)
		require.NoError(t, it.Err())

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

	addresses := []string{
		strings.ToLower(randomString(64, "")) + "@example1.com",
		strings.ToLower(randomString(64, "")) + "@example2.com",
		strings.ToLower(randomString(64, "")) + "@example3.com",
	}

	require.NoError(t, mg.CreateComplaints(ctx, testDomain, addresses))

	for _, address := range addresses {
		require.True(t, hasComplaint(address))
		require.NoError(t, mg.DeleteComplaint(ctx, testDomain, address))
		require.False(t, hasComplaint(address))
	}
}
