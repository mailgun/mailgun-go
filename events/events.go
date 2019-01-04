package events

import "errors"

type MailingListError struct {
	Message string
}

type MailingListEvent struct {
	Address string `json:"address"`
	ListID  string `json:"list-id"`
	SID     string `json:"sid"`
}

type MailingListMember struct {
	Subscribed bool
	Address    string
	Name       string
	Vars       []string
}

type DeliveryStatus struct {
	Message     *string     `json:"message,omitempty"`
	Code        interface{} `json:"code,omitempty"`
	Description *string     `json:"description,omitempty"`
	Retry       *int        `json:"retry-seconds,omitempty"`
}

type EventFlags struct {
	Authenticated bool `json:"is-authenticated"`
	Batch         bool `json:"is-batch"`
	Big           bool `json:"is-big"`
	Callback      bool `json:"is-callback"`
	DelayedBounce bool `json:"is-delayed-bounce"`
	SystemTest    bool `json:"is-system-test"`
	TestMode      bool `json:"is-test-mode"`
}

type ClientInfo struct {
	ClientType *ClientType `json:"client-type,omitempty"`
	ClientOS   *string     `json:"client-os,omitempty"`
	ClientName *string     `json:"client-name,omitempty"`
	DeviceType *DeviceType `json:"device-type,omitempty"`
	UserAgent  *string     `json:"user-agent,omitempty"`
}

type Geolocation struct {
	Country *string `json:"country,omitempty"`
	Region  *string `json:"region,omitempty"`
	City    *string `json:"city,omitempty"`
}

type Storage struct {
	URL string `json:"url"`
	Key string `json:"key"`
}

type Batch struct {
	ID string `json:"id"`
}

type Envelope struct {
	Sender      *string          `json:"sender,omitempty"`
	SendingHost *string          `json:"sending-host,omitempty"`
	SendingIP   *IP              `json:"sending-ip,omitempty"`
	Targets     *string          `json:"targets,omitempty"`
	Transport   *TransportMethod `json:"transport,omitempty"`
}

type StoredAttachment struct {
	Size        int    `json:"size"`
	Url         string `json:"url"`
	Name        string `json:"name"`
	ContentType string `json:"content-type"`
}

type EventMessage struct {
	Headers     map[string]string  `json:"headers,omitempty"`
	Recipients  []string           `json:"recipients,omitempty"`
	Attachments []StoredAttachment `json:"attachments,omitempty"`
	Size        *int               `json:"size,omitempty"`
}

func (em *EventMessage) ID() (string, error) {
	if em != nil && em.Headers != nil {
		if id, ok := em.Headers["message-id"]; ok {
			return id, nil
		}
	}
	return "", errors.New("message id not set")
}
