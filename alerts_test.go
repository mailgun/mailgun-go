package mailgun_test

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testAlertWebhookSigningKey = "0102030405060708090a0b0c0d0e0f10"

func TestListAlerts(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	resp, err := mg.ListAlerts(context.Background(), nil)
	require.NoError(t, err)
	assert.Len(t, resp.Events, 2)
}

func TestAddAlert(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	req := mtypes.AlertsEventSettingRequest{
		Channel:   mtypes.AlertsEmailChannel,
		EventType: "ip_listed",
		Settings: mtypes.AlertsChannelSettings{
			Emails: []string{"mail1@example.com", "mail2@example.com"},
		},
	}

	wantResp := mtypes.AlertsEventSettingResponse{
		Channel:    req.Channel,
		DisabledAt: nil,
		EventType:  req.EventType,
		ID:         ptr(uuid.MustParse("12345678-1234-5678-1234-123456789012")),
		Settings:   req.Settings,
	}

	resp, err := mg.AddAlert(context.Background(), req)
	require.NoError(t, err)
	assert.Equal(t, wantResp, *resp)
}

func TestDeleteAlert(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	err = mg.DeleteAlert(context.Background(), uuid.New())
	require.NoError(t, err)
}

func TestCalcAlertsHMAC(t *testing.T) {
	tests := map[string]struct {
		body              []byte
		webhookSigningKey string
		wantHEXSign       string
		wantErr           bool
	}{
		"positive": {
			body:              alertsWebhookBody(t),
			webhookSigningKey: testAlertWebhookSigningKey,
			wantHEXSign:       "8c82d17f19d19baf6cae658e3cf5db3c389309bcccfa490d27a5d39fa036dadf",
		},
		"empty_signing_key": {
			body:              alertsWebhookBody(t),
			webhookSigningKey: "",
			wantErr:           true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := mailgun.CalcAlertsHMAC(tt.body, tt.webhookSigningKey)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				gotSign := hex.EncodeToString(b)
				assert.Equal(t, tt.wantHEXSign, gotSign)
			}
		})
	}
}

func TestVerifyAlertsWebhookSign(t *testing.T) {
	tests := map[string]struct {
		body              []byte
		signHeader        string
		webhookSigningKey string
		want              bool
		wantErr           bool
	}{
		"verified": {
			body:              alertsWebhookBody(t),
			signHeader:        "8c82d17f19d19baf6cae658e3cf5db3c389309bcccfa490d27a5d39fa036dadf",
			webhookSigningKey: testAlertWebhookSigningKey,
			want:              true,
		},
		"not_verified": {
			body:              alertsWebhookBody(t),
			signHeader:        "beef",
			webhookSigningKey: testAlertWebhookSigningKey,
			want:              false,
		},
		"malformed_sign": {
			body:              alertsWebhookBody(t),
			signHeader:        "malformed",
			webhookSigningKey: testAlertWebhookSigningKey,
			wantErr:           true,
			want:              false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			isVerified, err := mailgun.VerifyAlertsWebhookSign(tt.body, tt.signHeader, tt.webhookSigningKey)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, isVerified)
			}
		})
	}
}

func TestVerifyAlertsWebhookSignFromRequest(t *testing.T) {
	tests := map[string]struct {
		body              []byte
		signHeader        string
		webhookSigningKey string
		want              bool
		wantErr           bool
	}{
		"verified": {
			body:              alertsWebhookBody(t),
			signHeader:        "8c82d17f19d19baf6cae658e3cf5db3c389309bcccfa490d27a5d39fa036dadf",
			webhookSigningKey: testAlertWebhookSigningKey,
			want:              true,
		},
		"not_verified": {
			body:              alertsWebhookBody(t),
			signHeader:        "beef",
			webhookSigningKey: testAlertWebhookSigningKey,
			want:              false,
		},
		"missing_sign_header": {
			body:              alertsWebhookBody(t),
			signHeader:        "",
			webhookSigningKey: testAlertWebhookSigningKey,
			want:              false,
		},
		"malformed_sign": {
			body:              alertsWebhookBody(t),
			signHeader:        "malformed",
			webhookSigningKey: testAlertWebhookSigningKey,
			wantErr:           true,
			want:              false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(tt.body))
			req.Header.Set(mtypes.AlertsWebhookSignHeader, tt.signHeader)

			isVerified, err := mailgun.VerifyAlertsWebhookSignFromRequest(req, tt.webhookSigningKey)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, isVerified)
			}

			t.Run("body_is_still_readable", func(t *testing.T) {
				gotBody, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				assert.Equal(t, tt.body, gotBody)
			})
		})
	}
}

func alertsWebhookBody(t *testing.T) []byte {
	t.Helper()

	type Signature struct {
		// Number of seconds passed since January 1, 1970.
		Timestamp int64 `json:"timestamp"`
		// Randomly generated string.
		Token string `json:"token"`
	}

	type AlertsWebhookReq struct {
		Signature Signature `json:"signature"`
		EventData any       `json:"event_data"`
	}

	body := AlertsWebhookReq{
		Signature: Signature{
			Timestamp: 1136239445,
			Token:     "abc",
		},
		EventData: map[string]string{
			"event": "ip_listed",
			"ip":    "1.1.1.1",
		},
	}
	b, err := json.Marshal(body)
	require.NoError(t, err)

	return b
}
