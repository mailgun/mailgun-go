package mailgun_test

import (
	"context"
	"testing"

	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateInboxPlacementTest(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	opts := mtypes.CreateInboxPlacementTestOptions{
		From:    "test@example.com",
		Subject: "test subject",
	}
	got, err := mg.CreateInboxPlacementTest(context.Background(), opts)
	require.NoError(t, err)

	assert.Equal(t, "result-id", got.ResultID)
	assert.Equal(t, "ibp-123@mailgun.net,seed@domain.tld", got.MailingList)
}
