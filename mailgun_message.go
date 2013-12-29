package mailgun

import (
	"net/mail"
)

type MailgunMessage struct {
	From    *mail.Address
	To      []*mail.Address
	Cc      []*mail.Address
	Bcc     []*mail.Address
	Subject string
	Text    string
	Html    string
}

func (m *MailgunMessage) AddRecipient(recipient *mail.Address) {
	m.To = append(m.To, recipient)
}

func (m *MailgunMessage) AddCC(recipient *mail.Address) {
	m.Cc = append(m.Cc, recipient)
}

func (m *MailgunMessage) AddBCC(recipient *mail.Address) {
	m.Bcc = append(m.Bcc, recipient)
}

func (m *MailgunMessage) validateMessage() bool {
	if m == nil {
		return false
	}

	if m.From == nil {
		return false
	}

	if m.From.Address == "" {
		return false
	}

	if !validateAddressList(m.To, true) {
		return false
	}

	if !validateAddressList(m.Cc, false) {
		return false
	}

	if !validateAddressList(m.Bcc, false) {
		return false
	}

	if m.Text == "" {
		return false
	}

	return true
}

func validateAddressList(list []*mail.Address, requireOne bool) bool {
	hasOne := false

	if list == nil {
		return !requireOne
	} else {
		for _, a := range list {
			if a.Address == "" {
				return false
			} else {
				hasOne = hasOne || true
			}
		}
	}

	return hasOne
}
