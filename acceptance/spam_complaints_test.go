// +build acceptance

package acceptance

import (
	"github.com/mailgun/mailgun-go"
	"testing"
)

func TestGetComplaints(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	publicApiKey := reqEnv(t, "MG_PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	n, complaints, err := mg.GetComplaints(-1, -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(complaints) != n {
		t.Fatalf("Expected %d complaints; got %d", n, len(complaints))
	}
}

func TestGetComplaintFromBazNoComplaint(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	publicApiKey := reqEnv(t, "MG_PUBLIC_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, publicApiKey)
	complaint, err := mg.GetSingleComplaint("baz@example.com")
	if err != nil {
		t.Fatal(err)
	}
	if complaint.Count != 0 {
		t.Fatal("Expected baz@example.com to have zero complaints.")
	}
}

func TestCreateDeleteComplaint(t *testing.T) {
	domain := reqEnv(t, "MG_DOMAIN")
	apiKey := reqEnv(t, "MG_API_KEY")
	mg := mailgun.NewMailgun(domain, apiKey, "")
	var check = func(count int) {
		c, _, err := mg.GetComplaints(mailgun.DefaultLimit, mailgun.DefaultSkip)
		if err != nil {
			t.Fatal(err)
		}
		if c != count {
			t.Fatalf("Expected baz@example.com to have %d complaints; got %d", count, c)
		}
	}

	check(0)

	err := mg.CreateComplaint("baz@example.com")
	if err != nil {
		t.Fatal(err)
	}

	check(1)

	err = mg.DeleteComplaint("baz@example.com")
	if err != nil {
		t.Fatal(err)
	}

	check(0)
}
