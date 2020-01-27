package mailgun

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v3/events"
)

func TestParseErrors(t *testing.T) {
	_, err := ParseEvent([]byte(""))
	ensure.DeepEqual(t, err.Error(), "failed to recognize event: EOF")

	_, err = ParseEvent([]byte(`{"event": "unknown_event"}`))
	ensure.DeepEqual(t, err.Error(), "unsupported event: 'unknown_event'")

	_, err = ParseEvent([]byte(`{
		"event": "accepted",
		"timestamp": "1420255392.850187"
	}`))
	errMsg := "failed to parse event 'accepted': parse error: expected number near offset 59 of '1420255392...'"
	ensure.DeepEqual(t, err.Error(), errMsg)
}

func TestParseSuccess(t *testing.T) {
	event, err := ParseEvent([]byte(`{
		"event": "accepted",
		"timestamp": 1420255392.850187,
		"user-variables": {
			"custom": "value",
			"parent": {"child": "user defined variable"},
			"a-list": [1,2,3,4,5]
		},
		"envelope": {
		    "sender": "noreply@example.com",
		    "transport": "smtp",
		    "mail-from": null,
		    "targets": "me@example.com"
		},
		"message": {
		    "headers": {
		        "to": "me@example.com",
		        "subject": "Test message going through the bus.",
		        "message-id": "20150103032312.125890.23497@example.com",
		        "from": "Example <noreply@example.com>"
		  },
		  "recipients": [
		    "me@example.com"
		  ],
		  "attachments": [],
		  "size": 6830
		},
		"tags": [
		    "sell_email_new"
		],
		"campaigns": [
		    {
		        "id": "d2yb8",
		        "name": "wantlist"
		    }
		],
		"recipient": "dude@example.com",
		"recipient-domain": "example.com",
		"method": "http",
		"flags": {
		  "is-system-test": false,
		  "is-big": false,
		  "is-test-mode": false,
		  "is-authenticated": false,
		  "is-routed": null
		}
	}`))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, reflect.TypeOf(event).String(), "*events.Accepted")
	subject := event.(*events.Accepted).Message.Headers.Subject
	ensure.DeepEqual(t, subject, "Test message going through the bus.")

	// Make sure the next event parsing attempt will zero the fields.
	event2, err := ParseEvent([]byte(`{
		"event": "accepted",
		"timestamp": 1533922516.538978,
		"recipient": "someone@example.com"
	}`))
	ensure.Nil(t, err)

	ensure.DeepEqual(t, event2.GetTimestamp(),
		time.Date(2018, 8, 10, 17, 35, 16, 538978048, time.UTC))
	ensure.DeepEqual(t, event2.(*events.Accepted).Message.Headers.Subject, "")
	// Make sure the second attempt of Parse doesn't overwrite the first event struct.
	ensure.DeepEqual(t, event.(*events.Accepted).Recipient, "dude@example.com")

	ensure.DeepEqual(t, event.(*events.Accepted).UserVariables.(map[string]interface{})["custom"], "value")
	child := event.(*events.Accepted).UserVariables.(map[string]interface{})["parent"].(map[string]interface{})["child"]
	ensure.DeepEqual(t, child, "user defined variable")
	aList := event.(*events.Accepted).UserVariables.(map[string]interface{})["a-list"].([]interface{})
	ensure.DeepEqual(t, aList, []interface{}{1.0, 2.0, 3.0, 4.0, 5.0})
}

func TestParseSuccessInvalidUserVariables(t *testing.T) {
	event, err := ParseEvent([]byte(`{
		"event": "accepted",
		"timestamp": 1420255392.850187,
		"user-variables": "Could not load user-variables. They were either truncated or invalid JSON"
	}`))
	ensure.Nil(t, err)
	ensure.DeepEqual(t, reflect.TypeOf(event).String(), "*events.Accepted")
	ensure.DeepEqual(t, event.(*events.Accepted).UserVariables, "Could not load user-variables. They were either truncated or invalid JSON")
}

func TestParseResponse(t *testing.T) {
	// Make sure the next event parsing attempt will zero the fields.
	evnts, err := parseResponse([]byte(`{
		"items": [
			{
				"event": "accepted",
				"timestamp": 1533922516.538978,
				"recipient": "someone@example.com"
			},
			{
				"event": "delivered",
				"timestamp": 1533922516.538978,
				"recipient": "test@mailgun.test"
			}
		],
		"paging": {
			"next": "https://next",
			"first": "https://first",
			"last": "https://last",
			"previous": "https://prev"
		}
	}`))
	ensure.Nil(t, err)

	ensure.DeepEqual(t, evnts[0].GetName(), "accepted")
	ensure.DeepEqual(t, evnts[0].(*events.Accepted).Recipient, "someone@example.com")

	ensure.DeepEqual(t, evnts[1].GetName(), "delivered")
	ensure.DeepEqual(t, evnts[1].(*events.Delivered).Recipient, "test@mailgun.test")
}

func TestTimeStamp(t *testing.T) {
	var event events.Generic
	ts := time.Date(2018, 8, 10, 17, 35, 16, 538978048, time.UTC)
	event.SetTimestamp(ts)
	ensure.DeepEqual(t, event.GetTimestamp(), ts)

	event.Timestamp = 1546899001.019501
	ensure.DeepEqual(t, event.GetTimestamp(),
		time.Date(2019, 1, 7, 22, 10, 01, 19501056, time.UTC))
}

func TestEventNames(t *testing.T) {
	for name := range EventNames {
		event, err := ParseEvent([]byte(fmt.Sprintf(`{"event": "%s"}`, name)))
		ensure.Nil(t, err)
		ensure.DeepEqual(t, event.GetName(), name)
	}
}

func TestEventMessageWithAttachment(t *testing.T) {
	body := []byte(`{
        "event": "delivered",
        "message": {
            "headers": {
                "to": "alice@example.com",
                "message-id": "alice@example.com",
                "from": "Bob <bob@example.com>",
                "subject": "Hi Alice"},
                "attachments": [{"filename": "doc.pdf",
                                 "content-type": "application/pdf",
                                 "size": 139214}],
                "size": 142698}}`)
	event, err := ParseEvent(body)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, event.(*events.Delivered).Message.Attachments[0].FileName, "doc.pdf")
}

func TestStored(t *testing.T) {
	var key = "WyJlODlmYTdhNTE3IiwghODAyMTQ4MDEzMTUiXSLCAib2xkY29rZSJd"
	var url = "https://api.mailgun.net/v2/domains/mg/messages/Wy9JZlJd"
	body := []byte(fmt.Sprintf(`{
        "event": "stored",
        "storage": {
            "key": "%s",
            "url": "%s"
        }}`, key, url))
	event, err := ParseEvent(body)
	ensure.Nil(t, err)
	ensure.DeepEqual(t, event.(*events.Stored).Storage.Key, key)
	ensure.DeepEqual(t, event.(*events.Stored).Storage.URL, url)
}
