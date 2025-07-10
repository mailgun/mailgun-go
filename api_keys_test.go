package mailgun_test

import (
	"context"
	"github.com/mailgun/mailgun-go/v5"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAPIKeys(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()
	list, err := mg.ListAPIKeys(ctx, nil)
	require.NoError(t, err)
	require.Len(t, list, 0)
}
