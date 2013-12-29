package mailgun

import (
)

type Message struct {
	From    string
	To      []string
	Cc      []string
	Bcc     []string
	Subject string
	Text    string
	Html    string
}

func (m *Message) AddRecipient(recipient string) {
	m.To = append(m.To, recipient)
}

func (m *Message) AddCC(recipient string) {
	m.Cc = append(m.Cc, recipient)
}

func (m *Message) AddBCC(recipient string) {
	m.Bcc = append(m.Bcc, recipient)
}

func (m *Message) validateMessage() bool {
	if m == nil {
		return false
	}

	if m.From == "" {
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

func validateAddressList(list []string, requireOne bool) bool {
	hasOne := false

	if list == nil {
		return !requireOne
	} else {
		for _, a := range list {
			if a == "" {
				return false
			} else {
				hasOne = hasOne || true
			}
		}
	}

	return hasOne
}
