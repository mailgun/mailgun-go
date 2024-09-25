package mailgun_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/facebookgo/ensure"
	"github.com/mailgun/errors"
	"github.com/mailgun/mailgun-go/v4"
)

func TestTemplateCRUD(t *testing.T) {
	mg := mailgun.NewMailgun(testDomain, testKey)
	mg.SetAPIBase(server.URL())
	ctx := context.Background()

	findTemplate := func(name string) bool {
		it := mg.ListTemplates(nil)

		var page []mailgun.Template
		for it.Next(ctx, &page) {
			for _, template := range page {
				if template.Name == name {
					return true
				}
			}
		}
		ensure.Nil(t, it.Err())
		return false
	}

	const (
		Name        = "Mailgun-Go-TestTemplateCRUD"
		Description = "Mailgun-Go Test Template Description"
		UpdatedDesc = "Mailgun-Go Test Updated Description"
	)

	tmpl := mailgun.Template{
		Name:        Name,
		Description: Description,
	}

	// Create a template
	ensure.Nil(t, mg.CreateTemplate(ctx, &tmpl))
	ensure.DeepEqual(t, tmpl.Name, strings.ToLower(Name))
	ensure.DeepEqual(t, tmpl.Description, Description)

	// Wait the template to show up
	ensure.Nil(t, waitForTemplate(mg, tmpl.Name))

	// Ensure the template is in the list
	ensure.True(t, findTemplate(tmpl.Name))

	// Update the description
	tmpl.Description = UpdatedDesc
	ensure.Nil(t, mg.UpdateTemplate(ctx, &tmpl))

	// Ensure update took
	updated, err := mg.GetTemplate(ctx, tmpl.Name)
	ensure.Nil(t, err)

	ensure.DeepEqual(t, updated.Description, UpdatedDesc)

	// Delete the template
	ensure.Nil(t, mg.DeleteTemplate(ctx, tmpl.Name))
}

func waitForTemplate(mg mailgun.Mailgun, id string) error {
	ctx := context.Background()
	var attempts int
	for attempts <= 5 {
		_, err := mg.GetTemplate(ctx, id)
		if err != nil {
			if mailgun.GetStatusFromErr(err) == 404 {
				time.Sleep(time.Second * 2)
				attempts += 1
				continue
			}
			return err
		}
		return nil

	}
	return errors.Errorf("Waited to long for template '%s' to show up", id)
}
