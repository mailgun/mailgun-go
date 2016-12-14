package mailgun

import (
	"testing"
)

func TestCreateUnsubscriber(t *testing.T) {
	email := reqEnv(t, "MG_EMAIL_ADDR")
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}

	// Create unsubscription record
	err = mg.Unsubscribe(email, "*")
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetUnsubscribes(t *testing.T) {
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}
	n, us, err := mg.GetUnsubscribes(DefaultLimit, DefaultSkip)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Received %d out of %d unsubscribe records.\n", len(us), n)
	if len(us) > 0 {
		t.Log("ID\tAddress\tCreated At\tTag\t")
		for _, u := range us {
			t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tag)
		}
	}
}

func TestGetUnsubscriptionByAddress(t *testing.T) {
	email := reqEnv(t, "MG_EMAIL_ADDR")
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}
	n, us, err := mg.GetUnsubscribesByAddress(email)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Received %d out of %d unsubscribe records.\n", len(us), n)
	if len(us) > 0 {
		t.Log("ID\tAddress\tCreated At\tTag\t")
		for _, u := range us {
			t.Logf("%s\t%s\t%s\t%s\t\n", u.ID, u.Address, u.CreatedAt, u.Tag)
		}
	}
}

func TestCreateDestroyUnsubscription(t *testing.T) {
	email := reqEnv(t, "MG_EMAIL_ADDR")
	mg, err := NewMailgunFromEnv()
	if err != nil {
		t.Fatalf("NewMailgunFromEnv() error - %s", err.Error())
	}

	// Create unsubscription record
	err = mg.Unsubscribe(email, "*")
	if err != nil {
		t.Fatal(err)
	}

	// Destroy the unsubscription record
	err = mg.RemoveUnsubscribe(email)
	if err != nil {
		t.Fatal(err)
	}
}
