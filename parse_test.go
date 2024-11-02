package mailgun

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/mailgun/mailgun-go/v4/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseErrors(t *testing.T) {
	_, err := ParseEvent([]byte(""))
	require.NotNil(t, err)
	// TODO(vtopc): do not compare strings, use errors.Is or errors.As:
	require.Contains(t, err.Error(), "failed to recognize event")

	_, err = ParseEvent([]byte(`{"event": "unknown_event"}`))
	require.EqualError(t, err, "unsupported event: 'unknown_event'")

	_, err = ParseEvent([]byte(`{
		"event": "accepted",
		"timestamp": "1420255392.850187"
	}`))
	// TODO(vtopc): do not compare strings, use errors.Is or errors.As:
	require.Contains(t, err.Error(), "failed to parse event 'accepted'")
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
		"storage": {
			"key": "AgEASDSGFB8y4--TSDGxvccvmQB==",
			"url": "https://storage.eu.mailgun.net/v3/domains/example.com/messages/AgEASDSGFB8y4--TSDGxvccvmQB=="
		},
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
	require.NoError(t, err)
	require.Equal(t, reflect.TypeOf(event).String(), "*events.Accepted")
	subject := event.(*events.Accepted).Message.Headers.Subject
	assert.Equal(t, "Test message going through the bus.", subject)
	assert.Equal(t, "AgEASDSGFB8y4--TSDGxvccvmQB==", event.(*events.Accepted).Storage.Key)

	// Make sure the next event parsing attempt will zero the fields.
	event2, err := ParseEvent([]byte(`{
		"event": "accepted",
		"timestamp": 1533922516.538978,
		"recipient": "someone@example.com"
	}`))
	require.NoError(t, err)

	assert.Equal(t, time.Date(2018, 8, 10, 17, 35, 16, 538978048, time.UTC), event2.GetTimestamp())
	assert.Equal(t, "", event2.(*events.Accepted).Message.Headers.Subject)
	// Make sure the second attempt of Parse doesn't overwrite the first event struct.
	assert.Equal(t, "dude@example.com", event.(*events.Accepted).Recipient)

	assert.Equal(t, "value", event.(*events.Accepted).UserVariables.(map[string]interface{})["custom"])
	child := event.(*events.Accepted).UserVariables.(map[string]interface{})["parent"].(map[string]interface{})["child"]
	assert.Equal(t, "user defined variable", child)
	aList := event.(*events.Accepted).UserVariables.(map[string]interface{})["a-list"].([]interface{})
	assert.Equal(t, []interface{}{1.0, 2.0, 3.0, 4.0, 5.0}, aList)
}

func TestParseSuccessInvalidUserVariables(t *testing.T) {
	event, err := ParseEvent([]byte(`{
		"event": "accepted",
		"timestamp": 1420255392.850187,
		"user-variables": "Could not load user-variables. They were either truncated or invalid JSON"
	}`))
	require.NoError(t, err)
	require.Equal(t, "*events.Accepted", reflect.TypeOf(event).String())
	assert.Equal(t, "Could not load user-variables. They were either truncated or invalid JSON",
		event.(*events.Accepted).UserVariables)
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
	require.NoError(t, err)

	assert.Equal(t, "accepted", evnts[0].GetName())
	assert.Equal(t, "someone@example.com", evnts[0].(*events.Accepted).Recipient)

	assert.Equal(t, "delivered", evnts[1].GetName())
	assert.Equal(t, "test@mailgun.test", evnts[1].(*events.Delivered).Recipient)
}

func TestTimeStamp(t *testing.T) {
	var event events.Generic
	ts := time.Date(2018, 8, 10, 17, 35, 16, 538978048, time.UTC)
	event.SetTimestamp(ts)
	assert.Equal(t, ts, event.GetTimestamp())

	event.Timestamp = 1546899001.019501
	assert.Equal(t, time.Date(2019, 1, 7, 22, 10, 01, 19501056, time.UTC), event.GetTimestamp())
}

func TestEventNames(t *testing.T) {
	for name := range EventNames {
		event, err := ParseEvent([]byte(fmt.Sprintf(`{"event": "%s"}`, name)))
		require.NoError(t, err)
		assert.Equal(t, name, event.GetName())
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
	require.NoError(t, err)
	assert.Equal(t, "doc.pdf", event.(*events.Delivered).Message.Attachments[0].FileName)
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
	require.NoError(t, err)
	assert.Equal(t, key, event.(*events.Stored).Storage.Key)
	assert.Equal(t, url, event.(*events.Stored).Storage.URL)
}
