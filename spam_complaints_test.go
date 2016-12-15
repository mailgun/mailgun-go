package mailgun

import (
	"strings"
	"testing"
)

func TestGetComplaints(t *testing.T) {
	reqEnv(t, "MG_PUBLIC_API_KEY")
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}
	n, complaints, err := mg.GetComplaints(-1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(complaints) != n {
		t.Fatalf("Expected %d complaints; got %d", n, len(complaints))
	}
}

func TestGetComplaintFromRandomNoComplaint(t *testing.T) {
	reqEnv(t, "MG_PUBLIC_API_KEY")
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}
	_, err = mg.GetSingleComplaint(randomString(64, "") + "@example.com")
	if err == nil {
		t.Fatal("Expected not-found error for missing complaint")
	}
	ure, ok := err.(*UnexpectedResponseError)
	if !ok {
		t.Fatal("Expected UnexpectedResponseError")
	}
	if ure.Actual != 404 {
		t.Fatalf("Expected 404 response code; got %d", ure.Actual)
	}
}

func TestCreateDeleteComplaint(t *testing.T) {
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}
	var hasComplaint = func(email string) bool {
		t.Logf("hasComplaint: %s\n", email)
		_, complaints, err := mg.GetComplaints(DefaultLimit, DefaultSkip)
		if err != nil {
			t.Fatal(err)
		}

		for _, complaint := range complaints {
			t.Logf("Complaint Address: %s\n", complaint.Address)
			if complaint.Address == email {
				return true
			}
		}
		return false
	}

	randomMail := strings.ToLower(randomString(64, "")) + "@example.com"

	if hasComplaint(randomMail) {
		t.Fatalf("Expected no complaints from '%s'", randomMail)
	}

	err = mg.CreateComplaint(randomMail)
	if err != nil {
		t.Fatal(err)
	}

	if !hasComplaint(randomMail) {
		t.Fatalf("Expected complaint from '%s'; but got none", randomMail)
	}

	err = mg.DeleteComplaint(randomMail)
	if err != nil {
		t.Fatal(err)
	}

	if hasComplaint(randomMail) {
		t.Fatalf("Expected no complaints after delete from '%s'", randomMail)
	}
}
