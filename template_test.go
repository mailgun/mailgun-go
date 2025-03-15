package mailgun_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/mailgun/errors"
	"github.com/mailgun/mailgun-go/v5"
	"github.com/mailgun/mailgun-go/v5/mtypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateCRUD(t *testing.T) {
	mg := mailgun.NewMailgun(testKey)
	err := mg.SetAPIBase(server.URL())
	require.NoError(t, err)

	ctx := context.Background()

	findTemplate := func(name string) bool {
		it := mg.ListTemplates(testDomain, nil)

		var page []mtypes.Template
		for it.Next(ctx, &page) {
			for _, template := range page {
				if template.Name == name {
					return true
				}
			}
		}
		require.NoError(t, it.Err())
		return false
	}

	const (
		Name        = "Mailgun-Go-TestTemplateCRUD"
		Description = "Mailgun-Go Test Template Description"
		UpdatedDesc = "Mailgun-Go Test Updated Description"
	)

	tmpl := mtypes.Template{
		Name:        Name,
		Description: Description,
	}

	// Create a template
	require.NoError(t, mg.CreateTemplate(ctx, testDomain, &tmpl))
	assert.Equal(t, strings.ToLower(Name), tmpl.Name)
	assert.Equal(t, Description, tmpl.Description)

	// Wait the template to show up
	require.NoError(t, waitForTemplate(mg, tmpl.Name))

	// Ensure the template is in the list
	require.True(t, findTemplate(tmpl.Name))

	// Update the description
	tmpl.Description = UpdatedDesc
	require.NoError(t, mg.UpdateTemplate(ctx, testDomain, &tmpl))

	// Ensure update took
	updated, err := mg.GetTemplate(ctx, testDomain, tmpl.Name)
	require.NoError(t, err)

	assert.Equal(t, UpdatedDesc, updated.Description)

	// Delete the template
	require.NoError(t, mg.DeleteTemplate(ctx, testDomain, tmpl.Name))
}

func waitForTemplate(mg mailgun.Mailgun, id string) error {
	ctx := context.Background()
	var attempts int
	for attempts <= 5 {
		_, err := mg.GetTemplate(ctx, testDomain, id)
		if err != nil {
			if mailgun.GetStatusFromErr(err) == 404 {
				time.Sleep(time.Second * 2)
				attempts++
				continue
			}

			return err
		}

		return nil
	}

	return errors.Errorf("Waited to long for template '%s' to show up", id)
}
