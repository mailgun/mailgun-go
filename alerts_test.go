package mailgun_test

import (
	"context"
	"encoding/hex"
	"encoding/json"
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

type AlertsWebhookReq struct {
	Signature Signature `json:"signature"`
	EventData any       `json:"event_data"`
}

type Signature struct {
	// Number of seconds passed since January 1, 1970.
	Timestamp int64 `json:"timestamp"`
	// Randomly generated string.
	Token string `json:"token"`
}

func TestCalcAlertsHMAC(t *testing.T) {
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
	testBody, err := json.Marshal(body)
	require.NoError(t, err)

	tests := map[string]struct {
		body              []byte
		webhookSigningKey string
		wantSign          string
		wantErr           error
	}{
		"positive": {
			body:              testBody,
			webhookSigningKey: testWebhookSigningKey,
			wantSign:          "9900bf2f2ae23f99dcb3b660906a20d3cdc89e67ee61cc7522f1f4d661240e04",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			b, err := mailgun.CalcAlertsHMAC(tt.body, tt.webhookSigningKey)
			require.Equal(t, tt.wantErr, err)

			gotSign := hex.EncodeToString(b)
			assert.Equal(t, tt.wantSign, gotSign)
		})
	}
}

// TODO(vtopc): add tests for VerifyAlertsWebhookSign
