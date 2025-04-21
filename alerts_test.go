package mailgun_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
