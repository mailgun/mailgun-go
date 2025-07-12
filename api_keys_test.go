package mailgun_test

import (
	"context"
	"github.com/mailgun/mailgun-go/v5"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestListAPIKeys(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	list, err := mg.ListAPIKeys(ctx, nil)
	require.NoError(t, err)
	require.Len(t, list, 2)
}

func TestCreateAPIKeys(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	key, err := mg.CreateAPIKey(ctx, "basic", nil)
	require.NoError(t, err)
	require.Equal(t, "1", key.ID)
	require.Equal(t, "basic", key.Role)
}

func TestDeleteAPIKey(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	msg, err := mg.DeleteAPIKey(ctx, "1")
	require.NoError(t, err)
	require.Equal(t, "success", msg)
}

func TestRegeneratePublicAPIKey(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	resp, err := mg.RegeneratePublicAPIKey(ctx)
	require.NoError(t, err)
	require.Equal(t, "public-1", resp.Key)
	require.Equal(t, "success", resp.Message)
}
