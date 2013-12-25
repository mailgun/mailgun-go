package mailgun

import (
	"net/mail"
	"testing"
)

const TEST_ADDRESS = "Michael Banzon <michael@banzon.dk>"

func TestMailgunMessage(t *testing.T) {
	m := new(MailgunMessage)
	address, _ := mail.ParseAddress(TEST_ADDRESS)

	checkMessageNotValid(m, t)
	m.From = address
	checkMessageNotValid(m, t)
	m.AddRecipient(address)
	checkMessageNotValid(m, t)
	m.Text = "Hello mail!"
}

func checkMessageNotValid(m *MailgunMessage, t *testing.T) {
	if m.validateMessage() {
		t.Error("Message should not be valid yet!")
	}
}
