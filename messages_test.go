package mailgun_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/stretchr/testify/require"
)

const (
	exampleHtml    = "<html><head /><body><p>Testing some <a href=\"http://google.com?q=abc&r=def&s=ghi\">Mailgun HTML awesomeness!</a> at www.kc5tja@yahoo.com</p></body></html>"
	exampleAMPHtml = `<!doctype html><html ⚡4email><head><meta charset="utf-8"><script async src="https://cdn.ampproject.org/v0.js"></script><style amp4email-boilerplate>body{visibility:hidden}</style><style amp-custom>h1{margin: 1rem;}</style></head><body><h1>Hello, I am an AMP EMAIL!</h1></body></html>`
	exampleMime    = `Content-Type: text/plain; charset="ascii"
Subject: Joe's Example Subject
From: Joe Example <joe@example.com>
To: BARGLEGARF <sam.falvo@rackspace.com>
Content-Transfer-Encoding: 7bit
Date: Thu, 6 Mar 2014 00:37:52 +0000

Testing some Mailgun MIME awesomeness!
`
	templateText  = "Greetings %recipient.name%!  Your reserved seat is at table %recipient.table%."
	exampleDomain = "testDomain"
	exampleAPIKey = "testAPIKey"
)

func init() {
	mailgun.Debug = true
	mailgun.CaptureCurlOutput = true
	mailgun.RedactCurlAuth = true
}

func spendMoney(t *testing.T, tFunc func()) {
	ok := os.Getenv("MG_SPEND_MONEY")
	if ok != "" {
		tFunc()
	} else {
		t.Log("Money spending not allowed, not running function.")
	}
}

func TestSendMGPlain(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendPlain:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGPlainWithTracking(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
		m.SetTracking(true)
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendPlainWithTracking:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGPlainAt(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
		m.SetDeliveryTime(time.Now().Add(5 * time.Minute))
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendPlainAt:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGSTO(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
		m.SetSTOPeriod("24h")
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendMGSTO:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGHtml(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
		m.SetHTML(exampleHtml)
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendHtml:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGAMPHtml(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
		m.SetHTML(exampleHtml)
		m.SetAMPHtml(exampleAMPHtml)
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendHtml:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGTracking(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText+"Tracking!\n", toUser)
		m.SetTracking(false)
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendTracking:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGTrackingClicksHtmlOnly(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
		m.SetHTML(exampleHtml)
		options := mailgun.TrackingOptions{
			Tracking:       true,
			TrackingClicks: "htmlonly",
			TrackingOpens:  true,
		}
		m.SetTrackingOptions(&options)
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendHtml:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGTag(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText+"Tags Galore!\n", toUser)
		m.AddTag("FooTag")
		m.AddTag("BarTag")
		m.AddTag("BlortTag")
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendTag:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGMIME(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMIMEMessage(io.NopCloser(strings.NewReader(exampleMime)), toUser)
		msg, id, err := mg.Send(ctx, m)
		require.NoError(t, err)
		t.Log("TestSendMIME:MSG(" + msg + "),ID(" + id + ")")
	})
}

func TestSendMGBatchFailRecipients(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")

		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText+"Batch\n")
		for i := 0; i < mailgun.MaxNumberOfRecipients; i++ {
			m.AddRecipient("") // We expect this to indicate a failure at the API
		}
		err := m.AddRecipientAndVariables(toUser, nil)
		// In case of error the SDK didn't send the message,
		// OR the API didn't check for empty To: headers.
		ensure.NotNil(t, err)
	})
}

func TestSendMGBatchRecipientVariables(t *testing.T) {
	if reason := mailgun.SkipNetworkTest(); reason != "" {
		t.Skip(reason)
	}

	spendMoney(t, func() {
		toUser := os.Getenv("MG_EMAIL_TO")
		mg, err := mailgun.NewMailgunFromEnv()
		require.NoError(t, err)

		ctx := context.Background()
		m := mailgun.NewMessage(fromUser, exampleSubject, templateText)
		err = m.AddRecipientAndVariables(toUser, map[string]interface{}{
			"name":  "Joe Cool Example",
			"table": 42,
		})
		require.NoError(t, err)
		_, _, err = mg.Send(ctx, m)
		require.NoError(t, err)
	})
}

func TestSendMGOffline(t *testing.T) {
	const (
		exampleDomain  = "testDomain"
		exampleAPIKey  = "testAPIKey"
		toUser         = "test@test.com"
		exampleMessage = "Queue. Thank you"
		exampleID      = "<20111114174239.25659.5817@samples.mailgun.org>"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, fmt.Sprintf("/v3/%s/messages", exampleDomain), req.URL.Path)
		require.Equal(t, fromUser, req.FormValue("from"))
		require.Equal(t, exampleSubject, req.FormValue("subject"))
		require.Equal(t, exampleText, req.FormValue("text"))
		require.Equal(t, toUser, req.FormValue("to"))
		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")
	ctx := context.Background()

	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)
}

func TestSendMGSeparateDomain(t *testing.T) {
	const (
		exampleDomain = "testDomain"
		signingDomain = "signingDomain"

		exampleAPIKey  = "testAPIKey"
		toUser         = "test@test.com"
		exampleMessage = "Queue. Thank you"
		exampleID      = "<20111114174239.25659.5817@samples.mailgun.org>"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, fmt.Sprintf("/v3/%s/messages", signingDomain), req.URL.Path)
		require.Equal(t, fromUser, req.FormValue("from"))
		require.Equal(t, exampleSubject, req.FormValue("subject"))
		require.Equal(t, exampleText, req.FormValue("text"))
		require.Equal(t, toUser, req.FormValue("to"))
		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")

	ctx := context.Background()
	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	m.AddDomain(signingDomain)

	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)
}

func TestSendMGMessageVariables(t *testing.T) {
	const (
		exampleDomain         = "testDomain"
		exampleAPIKey         = "testAPIKey"
		toUser                = "test@test.com"
		exampleMessage        = "Queue. Thank you"
		exampleID             = "<20111114174239.25659.5820@samples.mailgun.org>"
		exampleStrVarKey      = "test-str-key"
		exampleStrVarVal      = "test-str-val"
		exampleBoolVarKey     = "test-bool-key"
		exampleBoolVarVal     = "false"
		exampleMapVarKey      = "test-map-key"
		exampleMapVarStrVal   = `{"test":"123"}`
		exampleTemplateStrVal = `{"templateVariable":{"key":{"nested":"yes","status":"test"}}}`
	)
	var (
		exampleMapVarVal        = map[string]string{"test": "123"}
		exampleTemplateVariable = map[string]interface{}{
			"key": map[string]string{
				"nested": "yes",
				"status": "test",
			},
		}
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, fmt.Sprintf("/v3/%s/messages", exampleDomain), req.URL.Path)

		require.Equal(t, fromUser, req.FormValue("from"))
		require.Equal(t, exampleSubject, req.FormValue("subject"))
		require.Equal(t, exampleText, req.FormValue("text"))
		require.Equal(t, toUser, req.FormValue("to"))
		require.Equal(t, exampleMapVarStrVal, req.FormValue("v:"+exampleMapVarKey))
		require.Equal(t, exampleBoolVarVal, req.FormValue("v:"+exampleBoolVarKey))
		require.Equal(t, exampleStrVarVal, req.FormValue("v:"+exampleStrVarKey))
		require.Equal(t, exampleTemplateStrVal, req.FormValue("h:X-Mailgun-Variables"))
		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")

	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	m.AddVariable(exampleStrVarKey, exampleStrVarVal)
	m.AddVariable(exampleBoolVarKey, false)
	m.AddVariable(exampleMapVarKey, exampleMapVarVal)
	m.AddTemplateVariable("templateVariable", exampleTemplateVariable)

	msg, id, err := mg.Send(context.Background(), m)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)
}

func TestAddRecipientsError(t *testing.T) {
	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText)

	for i := 0; i < 1000; i++ {
		recipient := fmt.Sprintf("recipient_%d@example.com", i)
		require.NoError(t, m.AddRecipient(recipient))
	}

	err := m.AddRecipient("recipient_1001@example.com")
	ensure.NotNil(t, err)
	require.EqualError(t, err, "recipient limit exceeded (max 1000)")
}

func TestAddRecipientAndVariablesError(t *testing.T) {
	var err error

	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText)

	for i := 0; i < 1000; i++ {
		recipient := fmt.Sprintf("recipient_%d@example.com", i)
		err = m.AddRecipientAndVariables(recipient, map[string]interface{}{"id": i})
		require.NoError(t, err)
	}

	err = m.AddRecipientAndVariables("recipient_1001@example.com", map[string]interface{}{"id": 1001})
	require.Equal(t, err.Error(), "recipient limit exceeded (max 1000)")
}

func TestSendDomainError(t *testing.T) {
	cases := []struct {
		domain  string
		isValid bool
	}{
		{"http://example.com", false},
		{"example.com", true},
		{"mail.example.com", true},
		{"mail.service.example.com", true},
		{"http://example.com?email=yes", false},
		{"https://example.com", false},
		{"smtp://example.com", false},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rsp := `{
				"message":"Queued. Thank you",
				"id":"<20111114174239.25659.5817@samples.mailgun.org>"
				}`
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	for _, c := range cases {
		ctx := context.Background()
		mg := mailgun.NewMailgun(c.domain, exampleAPIKey)
		mg.SetAPIBase(srv.URL + "/v3")
		m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, "test@test.com")

		_, _, err := mg.Send(ctx, m)
		if c.isValid {
			require.NoError(t, err)
		} else {
			require.EqualError(t, err, "you called Send() with a domain that contains invalid characters")
		}
	}
}

func TestSendEOFError(t *testing.T) {
	const (
		exampleDomain = "testDomain"
		exampleAPIKey = "testAPIKey"
		toUser        = "test@test.com"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		panic("")
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")

	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	_, _, err := mg.Send(context.Background(), m)
	ensure.NotNil(t, err)
	ensure.StringContains(t, err.Error(), "remote server prematurely closed connection: Post ")
	ensure.StringContains(t, err.Error(), "EOF")
}

func TestHasRecipient(t *testing.T) {
	const (
		exampleDomain = "testDomain"
		exampleAPIKey = "testAPIKey"
		recipient     = "test@test.com"
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, fmt.Sprintf("/v3/%s/messages", exampleDomain), req.URL.Path)
		fmt.Fprint(w, `{"message":"Queued, Thank you", "id":"<20111114174239.25659.5820@samples.mailgun.org>"}`)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")

	// No recipient
	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText)
	_, _, err := mg.Send(context.Background(), m)
	require.EqualError(t, err, "message not valid")

	// Provided Bcc
	m = mailgun.NewMessage(fromUser, exampleSubject, exampleText)
	m.AddBCC(recipient)
	_, _, err = mg.Send(context.Background(), m)
	require.NoError(t, err)

	// Provided cc
	m = mailgun.NewMessage(fromUser, exampleSubject, exampleText)
	m.AddCC(recipient)
	_, _, err = mg.Send(context.Background(), m)
	require.Nil(t, err)
}

func TestResendStored(t *testing.T) {
	const (
		exampleDomain  = "testDomain"
		exampleAPIKey  = "testAPIKey"
		toUser         = "test@test.com"
		exampleMessage = "Queue. Thank you"
		exampleID      = "<20111114174239.25659.5820@samples.mailgun.org>"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, "/v3/some-url", req.URL.Path)
		require.Equal(t, toUser, req.FormValue("to"))

		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")

	msg, id, err := mg.ReSend(context.Background(), srv.URL+"/v3/some-url")
	ensure.NotNil(t, err)
	require.Equal(t, err.Error(), "must provide at least one recipient")

	msg, id, err = mg.ReSend(context.Background(), srv.URL+"/v3/some-url", toUser)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)
}

func TestAddOverrideHeader(t *testing.T) {
	const (
		exampleDomain  = "testDomain"
		exampleAPIKey  = "testAPIKey"
		toUser         = "test@test.com"
		exampleMessage = "Queue. Thank you"
		exampleID      = "<20111114174239.25659.5817@samples.mailgun.org>"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, fmt.Sprintf("/v3/%s/messages", exampleDomain), req.URL.Path)
		require.Equal(t, "custom-value", req.Header.Get("CustomHeader"))
		require.Equal(t, "example.com", req.Host)

		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")
	mg.AddOverrideHeader("Host", "example.com")
	mg.AddOverrideHeader("CustomHeader", "custom-value")
	ctx := context.Background()

	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	m.SetRequireTLS(true)
	m.SetSkipVerification(true)

	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)

	ensure.StringContains(t, mg.GetCurlOutput(), "Host:")
}

func TestOnBehalfOfSubaccount(t *testing.T) {
	const (
		exampleDomain  = "testDomain"
		exampleAPIKey  = "testAPIKey"
		toUser         = "test@test.com"
		exampleMessage = "Queue. Thank you"
		exampleID      = "<20111114174239.25659.5817@samples.mailgun.org>"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, fmt.Sprintf("/v3/%s/messages", exampleDomain), req.URL.Path)
		require.Equal(t, "custom-value", req.Header.Get("CustomHeader"))
		require.Equal(t, "example.com", req.Host)
		require.Equal(t, "mailgun.subaccount", req.Header.Get(mailgun.OnBehalfOfHeader))

		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")
	mg.AddOverrideHeader("Host", "example.com")
	mg.AddOverrideHeader("CustomHeader", "custom-value")
	mg.SetOnBehalfOfSubaccount("mailgun.subaccount")
	ctx := context.Background()

	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	m.SetRequireTLS(true)
	m.SetSkipVerification(true)

	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)

	ensure.StringContains(t, mg.GetCurlOutput(), "Host:")
}

func TestCaptureCurlOutput(t *testing.T) {
	const (
		exampleDomain  = "testDomain"
		exampleAPIKey  = "testAPIKey"
		toUser         = "test@test.com"
		exampleMessage = "Queue. Thank you"
		exampleID      = "<20111114174239.25659.5817@samples.mailgun.org>"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, fmt.Sprintf("/v3/%s/messages", exampleDomain), req.URL.Path)
		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")
	ctx := context.Background()

	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)

	ensure.StringContains(t, mg.GetCurlOutput(), "curl")
	t.Logf("%s", mg.GetCurlOutput())
}

func TestSendTLSOptions(t *testing.T) {
	const (
		exampleDomain  = "testDomain"
		exampleAPIKey  = "testAPIKey"
		toUser         = "test@test.com"
		exampleMessage = "Queue. Thank you"
		exampleID      = "<20111114174239.25659.5817@samples.mailgun.org>"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, http.MethodPost, req.Method)
		require.Equal(t, fmt.Sprintf("/v3/%s/messages", exampleDomain), req.URL.Path)
		require.Equal(t, fromUser, req.FormValue("from"))
		require.Equal(t, exampleSubject, req.FormValue("subject"))
		require.Equal(t, exampleText, req.FormValue("text"))
		require.Equal(t, toUser, req.FormValue("to"))
		require.Equal(t, "true", req.FormValue("o:require-tls"))
		require.Equal(t, "true", req.FormValue("o:skip-verification"))
		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")
	ctx := context.Background()

	m := mailgun.NewMessage(fromUser, exampleSubject, exampleText, toUser)
	m.SetRequireTLS(true)
	m.SetSkipVerification(true)

	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)
}

func TestSendTemplate(t *testing.T) {
	const (
		exampleDomain  = "testDomain"
		exampleAPIKey  = "testAPIKey"
		toUser         = "test@test.com"
		exampleMessage = "Queue. Thank you"
		exampleID      = "<20111114174239.25659.5817@samples.mailgun.org>"
		templateName   = "my-template"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, templateName, req.FormValue("template"))
		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")
	ctx := context.Background()

	m := mailgun.NewMessage(fromUser, exampleSubject, "", toUser)
	m.SetTemplate(templateName)

	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)
}

func TestSendTemplateOptions(t *testing.T) {
	const (
		exampleDomain      = "testDomain"
		exampleAPIKey      = "testAPIKey"
		toUser             = "test@test.com"
		exampleMessage     = "Queue. Thank you"
		exampleID          = "<20111114174239.25659.5817@samples.mailgun.org>"
		templateName       = "my-template"
		templateVersionTag = "initial"
		templateRenderText = "yes"
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		require.Equal(t, templateName, req.FormValue("template"))
		require.Equal(t, templateVersionTag, req.FormValue("t:version"))
		require.Equal(t, templateRenderText, req.FormValue("t:text"))
		rsp := fmt.Sprintf(`{"message":"%s", "id":"%s"}`, exampleMessage, exampleID)
		fmt.Fprint(w, rsp)
	}))
	defer srv.Close()

	mg := mailgun.NewMailgun(exampleDomain, exampleAPIKey)
	mg.SetAPIBase(srv.URL + "/v3")
	ctx := context.Background()

	m := mailgun.NewMessage(fromUser, exampleSubject, "", toUser)
	m.SetTemplate(templateName)
	m.SetTemplateRenderText(true)
	m.SetTemplateVersion(templateVersionTag)

	msg, id, err := mg.Send(ctx, m)
	require.NoError(t, err)
	require.Equal(t, exampleMessage, msg)
	require.Equal(t, exampleID, id)
}
